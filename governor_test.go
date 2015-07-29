// governor_test.go
package governor

import (
	"github.com/stretchr/testify/assert"	
	"testing"
	"net/http"
)

type StubConfig struct {
	key, value string
}

func TestConsulAccess(t *testing.T) {
	
	// Test server that always responds with 200 code, and specific payload
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(200)fmt.Fprintln(w, `[{"email":"bob@example.com","status":"sent","reject_reason":"hard-bounce","_id":"1"}]`)}))
	
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// Make a http.Client with the transport
	httpClient := &http.Client{Transport: transport}

	// Make an API client and inject
	client := &Client{server.URL, httpClient}

	expected := StubConfig{key: "ssl_key", value: "/path/to/key"}
	
	attr := GetAttribute(expected.key)
	
	assert.Equal(t, attr, expected.value, "The two words should be equal")
	
}