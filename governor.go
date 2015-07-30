// governor
package governor

import (
	"net/http"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
)

func GetAttribute (key string, defaultClient *http.Client) string {

	// Get client
	config := api.DefaultConfig()
	
	// Override the client if we want more control
	if defaultClient != nil {
		config.HttpClient = defaultClient	
	}
	
	client, _ := api.NewClient(config)
	
	// Key-value end point
	kv := client.KV()
	
	keyValue, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	}
	
	// Get the value and convert to string
	byteValue := keyValue.Value
	
	stringValue := string(byteValue[:])

	return stringValue
}

func GetConfigFromFile (fileName string) string {
	
	// Open the file
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	
	return string(contents)
}