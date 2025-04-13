package models

type AnalysisRequest struct {
	URL string `json:"url" example:"https://example.com"`
}

type HeadingCount struct {
	H1 int `json:"h1" example:"1"`
	H2 int `json:"h2" example:"3"`
	H3 int `json:"h3" example:"3"`
	H4 int `json:"h4" example:"2"`
	H5 int `json:"h5" example:"3"`
	H6 int `json:"h6" example:"2"`
}

type LinkAnalysis struct {
	Internal     int `json:"internal" example:"5"`
	External     int `json:"external" example:"3"`
	Inaccessible int `json:"inaccessible" example:"1"`
}

type AnalysisResponse struct {
	HTMLVersion       string       `json:"htmlVersion" example:"HTML5"`
	Title             string       `json:"title" example:"Example Domain"`
	Headings          HeadingCount `json:"headings"`
	Links             LinkAnalysis `json:"links"`
	ContainsLoginForm bool         `json:"containsLoginForm" example:"false"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode" example:"502"`
	Message    string `json:"message" example:"Failed to analyze URL: HTTP error 404 Not Found"`
}
