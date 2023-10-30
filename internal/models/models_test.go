package internal

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestHeader(t *testing.T) {
	request := httptest.NewRequest("GET", "https://school.com/courses", nil)
	reqWithHeaders := setHeaders(request)
	assert.NotNil(t, reqWithHeaders.Header.Get("apiKey"))
}
