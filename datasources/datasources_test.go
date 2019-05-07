package datasources

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParse(t *testing.T) {
	expected := &url.URL{
		Scheme:   "http",
		Host:     "example.com",
		Path:     "/foo.json",
		RawQuery: "bar",
	}
	u, err := ParseSourceURL("http://example.com/foo.json?bar")
	assert.NoError(t, err)
	assert.EqualValues(t, expected, u)
}
