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
	
	expected := StubConfig{key: "ssl_key", value: "/path/to/key"}

	// Make a stub file
	stubFileName := "config_file.conf"
	stubContent := fmt.Sprintf(`{"%s": "%s"}`, expected.key, expected.value)
	err := ioutil.WriteFile(stubFileName, []byte(stubContent), 0644)
	if err != nil {
		panic(err)
	}
	
	// Cleanup the file that was made
	defer os.Remove(stubFileName)
	
	// Load the config file
	config := GetConfigFromFile(stubFileName)
	
	// Check that the key exists
	value, ok := config[expected.key]	
	assert.True(t, ok)
	assert.Equal(t, value, expected.value)
	
}

func TestConsulWriteToDisk(t *testing.T) {
	
	expected_file := "config.conf"
	
	stubMap := map[string]string{
		expected_file: `{"ssl_key": "/path/to/key"}`,
	}

	// Generate a file from the stub data
	MakeConfigFiles(stubMap)
	
	// Check that the config file was created
	_, err := os.Stat(expected_file)
	assert.Nil(t, err)

	// This is only called if the first assertion passes
	defer os.Remove(expected_file)

	// Check the contents is sensible
	contents, err := ioutil.ReadFile(expected_file)
	if err != nil {
		panic(err)
	}
	
	// Check the file contains relevant information
	assert.Contains(
		t,
		string(contents),
		"ssl_key",
		"The output should contain specific text, but it does not",
	)
	
}