// governor_test.go
package governor

import (
	"fmt"
	"github.com/stretchr/testify/assert"	
	"testing"
	"net/http"
	"net/url"
	"net/http/httptest"
	"encoding/base64"
	"io/ioutil"
	"os"
)

type StubConfig struct {
	key, value string
}

func TestConsulAccess(t *testing.T) {
	
	expected := StubConfig{key: "ssl_key", value: "/path/to/key"}
	
	// Test server that always responds with 200 code, and specific payload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		valueBase64 := base64.StdEncoding.EncodeToString([]byte(expected.value))

		response := fmt.Sprintf(`[{
			"CreateIndex": 100,
		    "ModifyIndex": 200,
		    "LockIndex": 200,
		    "Key": "ssl_key",
		    "Flags": 0,
		    "Value": "%s",
		    "Session": "adf4238a-882b-9ddc-4a9d-5b6758e4159e"
  		}]`, valueBase64)

		fmt.Fprintln(w, response)
	}))
	
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// Make a http.Client with the transport
	httpClient := &http.Client{Transport: transport}

	attr := GetAttribute(expected.key, httpClient)
	
	assert.Equal(t, attr, expected.value, "The two words should be equal")
	
}

func TestParseFileCorrectly(t *testing.T) {
	
	// Make a stub file
	stubContent := []byte("key config_file.conf")
	err := ioutil.WriteFile("config_file.conf", stubContent, 0644)
	if err != nil{
		panic(err)
	}
	
	// Load the config file
	config := GetConfigFromFile("config_file.conf")
	assert.Equal(t, len(config), 1, "The length of the config should be 1")
	
	err = os.Remove("config_file.conf")
	if err != nil {
		panic(err)
	}
	
}