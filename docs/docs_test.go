package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerInfoInitialization(t *testing.T) {
	// initialized good?
	assert.NotNil(t, SwaggerInfo, "SwaggerInfo should not be nil")

	assert.Equal(t, "Web Page Analyzer API", SwaggerInfo.Title)
	assert.Equal(t, "1.0", SwaggerInfo.Version)
	assert.Equal(t, "localhost:8080", SwaggerInfo.Host)
	assert.Equal(t, "/", SwaggerInfo.BasePath)
	assert.Contains(t, SwaggerInfo.Schemes, "http")

	// template exists ?
	assert.NotEmpty(t, SwaggerInfo.SwaggerTemplate)
}
