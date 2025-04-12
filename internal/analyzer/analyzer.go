package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/maheshjq/web-analyzer_v1/internal/models"
	"golang.org/x/net/html"
)

// Analyzer handles webpage analysis
type Analyzer struct {
	client *http.Client
}

// NewAnalyzer creates a new Analyzer instance
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

	// Detect HTML version
	result.HTMLVersion = detectHTMLVersion(doc)

	// Extract title
	result.Title = extractTitle(doc)

	// Count headings
	countHeadings(doc, &result.Headings)

	// Analyze links
	result.Links = analyzeLinks(doc, baseURL.Host, a.client)

	// Detect login form
	result.ContainsLoginForm = detectLoginForm(doc)

	return result, nil
}

// detectHTMLVersion determines the HTML version based on the doctype
func detectHTMLVersion(doc *html.Node) string {
	// Check for doctype
	if doc.Type == html.DocumentNode {
		for child := doc.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.DoctypeNode {
				doctype := child.Data
				
				// Check for HTML5
				if strings.Contains(doctype, "html") && !strings.Contains(doctype, "DTD") {
					return "HTML5"
				}
				
				// Check for HTML 4
				if strings.Contains(doctype, "HTML 4") || strings.Contains(doctype, "HTML 4.01") {
					return "HTML 4.01"
				}
				
				// Check for XHTML
				if strings.Contains(doctype, "XHTML") {
					if strings.Contains(doctype, "1.0") {
						return "XHTML 1.0"
					}
					if strings.Contains(doctype, "1.1") {
						return "XHTML 1.1"
					}
				}
				
				// Return something generic if we found a doctype but couldn't identify it precisely
				return "Unknown DOCTYPE: " + doctype
			}
		}
	}

	// If no doctype is found, check for HTML5 elements
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

// extractTitle extracts the title from the document
func extractTitle(doc *html.Node) string {
	var title string
	var crawler func(*html.Node)
	
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
				return
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	
	crawler(doc)
	return title
}

// countHeadings counts the number of each heading type (h1-h6)
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

// analyzeLinks analyzes and categorizes links in the document
func analyzeLinks(doc *html.Node, host string, client *http.Client) models.LinkAnalysis {
	var links []string
	
	// Extract all links
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
					break
				}
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	
	crawler(doc)
	
	// Categorize links (internal vs external)
	var internal, external, inaccessible int
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	// Create a client with shorter timeout for checking link accessibility
	linkClient := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}
	
	// Process each link
	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			
			// Skip empty links and javascript
			if l == "" || strings.HasPrefix(l, "javascript:") {
				return
			}
			
			// Check if internal or external
			isInternal := isInternalLink(l, host)
			
			mu.Lock()
			if isInternal {
				internal++
			} else {
				external++
			}
			
			// Check accessibility (only for HTTP(S) links)
			if strings.HasPrefix(l, "http") && !isAccessibleLink(l, linkClient) {
				inaccessible++
			}
			mu.Unlock()
		}(link)
	}
	
	wg.Wait()
	
	return models.LinkAnalysis{
		Internal:     internal,
		External:     external,
		Inaccessible: inaccessible,
	}
}

// isInternalLink determines if a link is internal (same domain) or external
func isInternalLink(href, host string) bool {
	// Fragment or empty links are internal
	if href == "" || strings.HasPrefix(href, "#") {
		return true
	}
	
	// Relative links are internal
	if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../") {
		return true
	}
	
	// Absolute links need checking
	u, err := url.Parse(href)
	if err != nil {
		return false
	}
	
	return u.Host == host || u.Host == ""
}

// isAccessibleLink checks if a link is accessible
func isAccessibleLink(link string, client *http.Client) bool {
	// Fragment links are considered accessible
	if strings.HasPrefix(link, "#") {
		return true
	}
	
	// Try to fetch the link
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
	
	// Look for indicators of a login form
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
			
			// Check attribute indicators
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