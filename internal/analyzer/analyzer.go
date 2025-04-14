package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
)

type Analyzer struct {
	client *http.Client
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Analyze performs a full analysis of the webpage at the given URL
func (a *Analyzer) Analyze(targetURL string) (*models.AnalysisResponse, error) {
	// Fetch the page
	resp, err := a.client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Parse the base URL for link analysis
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Analyze the document
	result := &models.AnalysisResponse{
		Headings: models.HeadingCount{},
		Links:    models.LinkAnalysis{},
	}

	result.HTMLVersion = detectHTMLVersion(doc)

	result.Title = extractTitle(doc)

	countHeadings(doc, &result.Headings)

	result.Links = analyzeLinks(doc, baseURL.Host, a.client)

	result.ContainsLoginForm = detectLoginForm(doc)

	return result, nil
}

func detectHTMLVersion(doc *html.Node) string {
	if doc.Type == html.DocumentNode {
		for child := doc.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.DoctypeNode {
				// HTML5: <!DOCTYPE html> with no attributes
				if strings.ToLower(child.Data) == "html" && len(child.Attr) == 0 {
					return "HTML5"
				}
				// Check for older versions
				for _, attr := range child.Attr {
					if attr.Key == "public" {
						pubID := strings.ToLower(attr.Val)
						if strings.Contains(pubID, "html 4") {
							return "HTML 4.01"
						} else if strings.Contains(pubID, "xhtml 1.0") {
							return "XHTML 1.0"
						} else if strings.Contains(pubID, "xhtml 1.1") {
							return "XHTML 1.1"
						}
					}
				}
				return "Unknown DOCTYPE"
			}
		}
	}
	// Check for HTML5 elements if no doctype
	html5Elements := []string{"article", "aside", "audio", "canvas", "footer", "header", "nav", "section", "video"}
	for _, element := range html5Elements {
		if findElement(doc, element) {
			return "HTML5 (No DOCTYPE)"
		}
	}
	return "Unknown (No DOCTYPE)"
}

// findElement searches for a specific element in the document
func findElement(n *html.Node, tagName string) bool {
	if n.Type == html.ElementNode && strings.ToLower(n.Data) == tagName {
		return true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if findElement(c, tagName) {
			return true
		}
	}

	return false
}

func extractTitle(doc *html.Node) string {
	var title string
	var findTitle func(*html.Node) bool
	findTitle = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "title" {
			// Get text content with HTML tags stripped
			var b strings.Builder
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					b.WriteString(c.Data)
				} else if c.Type == html.ElementNode {
					// Recursively get text from child elements
					var childText string
					for child := c.FirstChild; child != nil; child = child.NextSibling {
						if child.Type == html.TextNode {
							childText += child.Data
						}
					}
					b.WriteString(childText)
				}
			}
			title = b.String()
			return true // Stop searching after finding first title
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findTitle(c) {
				return true // Propagate the "found" signal up
			}
		}
		return false
	}
	findTitle(doc)
	return title
}

func countHeadings(doc *html.Node, headings *models.HeadingCount) {
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				headings.H1++
			case "h2":
				headings.H2++
			case "h3":
				headings.H3++
			case "h4":
				headings.H4++
			case "h5":
				headings.H5++
			case "h6":
				headings.H6++
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(doc)
}

func analyzeLinks(doc *html.Node, host string, client *http.Client) models.LinkAnalysis {
	var links []string
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}
	extractLinks(doc)

	// Define a result struct for channel communication
	type linkResult struct {
		isInternal     bool
		isExternal     bool
		isInaccessible bool
	}

	resultCh := make(chan linkResult, len(links))
	var wg sync.WaitGroup

	for _, link := range links {
		if link == "" || strings.HasPrefix(link, "javascript:") {
			continue // Skip empty or js links
		}

		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			result := linkResult{
				isInternal:     isInternalLink(l, host),
				isExternal:     !isInternalLink(l, host),
				isInaccessible: false,
			}

			if strings.HasPrefix(l, "http") {
				fmt.Println("Checking accessibility for:", l) // Debug print
				if !isAccessibleLink(l, client) {
					fmt.Println("Inaccessible link found:", l) // Debug print
					result.isInaccessible = true
				}
			}

			resultCh <- result
		}(link)
	}

	// Close the channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Count the results
	var internal, external, inaccessible int
	for result := range resultCh {
		if result.isInternal {
			internal++
		}
		if result.isExternal {
			external++
		}
		if result.isInaccessible {
			inaccessible++
		}
	}

	return models.LinkAnalysis{
		Internal:     internal,
		External:     external,
		Inaccessible: inaccessible,
	}
}

func isInternalLink(href, host string) bool {
	if href == "" || strings.HasPrefix(href, "#") {
		return true
	}
	if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../") {
		return true
	}
	u, err := url.Parse(href)
	if err != nil || (u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	return u.Host == host || u.Host == ""
}

func isAccessibleLink(link string, client *http.Client) bool {
	// Fragment links are considered accessible
	if strings.HasPrefix(link, "#") {
		return true
	}

	resp, err := client.Head(link)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 2xx and 3xx status codes are considered accessible
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// detectLoginForm detects if the document contains a login form
func detectLoginForm(doc *html.Node) bool {
	var hasLoginForm bool

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if hasLoginForm {
			return
		}

		if n.Type == html.ElementNode && n.Data == "form" {
			// Check for form attributes that suggest a login form
			var formAction, formId, formClass string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "action":
					formAction = strings.ToLower(attr.Val)
				case "id":
					formId = strings.ToLower(attr.Val)
				case "class":
					formClass = strings.ToLower(attr.Val)
				}
			}

			if strings.Contains(formAction, "login") ||
				strings.Contains(formAction, "signin") ||
				strings.Contains(formId, "login") ||
				strings.Contains(formId, "signin") ||
				strings.Contains(formClass, "login") ||
				strings.Contains(formClass, "signin") {
				hasLoginForm = true
				return
			}

			// Check form elements for password inputs
			var formCrawler func(*html.Node)
			formCrawler = func(node *html.Node) {
				if node.Type == html.ElementNode && node.Data == "input" {
					for _, attr := range node.Attr {
						if attr.Key == "type" && attr.Val == "password" {
							hasLoginForm = true
							return
						}
					}
				}

				for c := node.FirstChild; c != nil; c = c.NextSibling {
					formCrawler(c)
				}
			}

			formCrawler(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(doc)
	return hasLoginForm
}
