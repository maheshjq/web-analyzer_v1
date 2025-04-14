package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalysisResponseSerialization(t *testing.T) {
	response := AnalysisResponse{
		HTMLVersion: "HTML5",
		Title:       "Test Page",
		Headings: HeadingCount{
			H1: 1,
			H2: 2,
			H3: 3,
		},
		Links: LinkAnalysis{
			Internal:     5,
			External:     3,
			Inaccessible: 1,
		},
		ContainsLoginForm: true,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Deserialize back to struct
	var decoded AnalysisResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, response.HTMLVersion, decoded.HTMLVersion)
	assert.Equal(t, response.Title, decoded.Title)
	assert.Equal(t, response.Headings.H1, decoded.Headings.H1)
	assert.Equal(t, response.Headings.H2, decoded.Headings.H2)
	assert.Equal(t, response.Headings.H3, decoded.Headings.H3)
	assert.Equal(t, response.Links.Internal, decoded.Links.Internal)
	assert.Equal(t, response.Links.External, decoded.Links.External)
	assert.Equal(t, response.Links.Inaccessible, decoded.Links.Inaccessible)
	assert.Equal(t, response.ContainsLoginForm, decoded.ContainsLoginForm)
}

func TestErrorResponseSerialization(t *testing.T) {
	errorResp := ErrorResponse{
		StatusCode: 404,
		Message:    "Page not found",
	}

	jsonData, err := json.Marshal(errorResp)
	assert.NoError(t, err)

	var decoded ErrorResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, errorResp.StatusCode, decoded.StatusCode)
	assert.Equal(t, errorResp.Message, decoded.Message)
}
