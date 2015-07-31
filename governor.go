// governor
package main

import (
	"os"
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"encoding/json"
)

const (
	CONSUL_ADDRESS string = "CONSUL_HOST"
	CONSUL_PORT string = "CONSUL_PORT"
)

func GetAttribute (key string, defaultClient *http.Client) string {

	// Get client
	config := api.DefaultConfig()
	
	// Override the client if we want more control
	if defaultClient != nil {
		config.HttpClient = defaultClient	
	}
	
	// Check environment variables
	consulPort := os.Getenv(CONSUL_PORT)
	if consulPort == "" {
		consulPort = "8500"
	}
	
	if consulAddress := os.Getenv(CONSUL_ADDRESS); consulAddress != "" {
		config.Address = consulAddress + ":" + consulPort
	}
	log.Println("Set the address to: ", config.Address)
	
	// Load the client
	client, _ := api.NewClient(config)
	
	// Key-value end point
	kv := client.KV()
	
	keyValue, _, err := kv.Get(key, nil)
	if err != nil {
		log.Println("error!!", err)
		panic(err)
	}
	
	// Get the value and convert to string
	byteValue := keyValue.Value
	
	stringValue := string(byteValue[:])

	return stringValue
}

func GetConfigFromFile (fileName string) map[string]string {
	
	// Open the file
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	
	// Unmarshal the JSON content
	configMap := make(map[string]string)
	json.Unmarshal(contents, &configMap)
	
	return configMap
}

func MakeConfigFiles (configMap map[string]string) {
	
	for filePath, fileContents := range configMap {
		// Write the file with its relevant contents
		ioutil.WriteFile(filePath, []byte(fileContents), 0644)		
	}

}

func checkFileExists(fileName string) {
	_, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("You need to specify a config file that exists.")
		panic(err)
	}	
}

func Govern (configFile string, defaultClient *http.Client) {
	
	// Parse the config file
	configMap := GetConfigFromFile(configFile)
	
	// Obtain the content from Consul and place in map
	var configContent string
	outputConfigMap := make(map[string]string)
	
	for consulKey, configPath := range configMap {
		configContent = GetAttribute(consulKey, defaultClient)
		
		outputConfigMap[configPath] = configContent
	}
	
	// Make the config files
	MakeConfigFiles(outputConfigMap)
}

func main () {
	
	// Definitions of allowed input flags
	configFilePtr := flag.String("config", "govern.conf", "Config file.")
	
	// Parse all the flags based on definitions
	flag.Parse()
	
	// Check config file exists
	checkFileExists(*configFilePtr)
	
	// Runtime routine
	Govern(*configFilePtr, nil)
}
