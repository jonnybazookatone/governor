// governor
package governor

import (
	"github.com/hashicorp/consul/api"
)

func GetAttribute (key string) string {

	// Get client
	client, _ := api.NewClient(api.DefaultConfig())
	
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