// governor
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	CONSUL_ADDRESS string = "CONSUL_HOST"
	CONSUL_PORT    string = "CONSUL_PORT"
)

func GetAttribute(key string, defaultClient *http.Client) string {

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

	log.Println("Attempting to retrieve key: ", key)
	keyValue, _, err := kv.Get(key, nil)
	if err != nil {
		log.Fatal("Error raised when attempting to get keys from consul: ", err)
	}
	if keyValue == nil {
		log.Fatal("Key supplied returned a nil value - does it exist: ", keyValue)
	}

	// Get the value and convert to string
	byteValue := keyValue.Value

	stringValue := string(byteValue[:])
	log.Println("Consul returned", stringValue)

	return stringValue
}

func GetConfigFromFile(fileName string) map[string]string {

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

func MakeConfigFiles(configMap map[string]string) {

	for filePath, fileContents := range configMap {

		// Does the output folder exist? If not, make it
		dirPath, _ := filepath.Abs(filepath.Dir(filePath))
		src, err := os.Stat(dirPath)
		if src == nil {
			// Create folder
			log.Println("Folder does not exist, making:", dirPath)
			err := os.MkdirAll(dirPath, 0777)
			if err != nil {
				log.Fatal("Unexpected error: ", err)
			}

		} else if err != nil {
			log.Fatal("Unexpected error: ", src, err)
		}

		// Write the file with its relevant contents
		log.Printf("Writing config file %s\n", filePath)
		ioutil.WriteFile(filePath, []byte(fileContents), 0777)
	}

}

func checkFileExists(fileName string) {
	_, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("You need to specify a config file that exists.")
		panic(err)
	}
}

func Govern(configFile string, defaultClient *http.Client) {

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

func main() {

	// Definitions of allowed input flags
	configFilePtr := flag.String("c", "govern.conf", "Config file.")

	// Parse all the flags based on definitions
	flag.Parse()

	// Check config file exists
	checkFileExists(*configFilePtr)

	// Runtime routine
	log.Println("Using config file: ", *configFilePtr)
	Govern(*configFilePtr, nil)
}
