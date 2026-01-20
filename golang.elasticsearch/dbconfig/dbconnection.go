package config

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	once   sync.Once
	client *elasticsearch.Client
)

func Connection() *elasticsearch.Client {
	once.Do(func() {

		cfg := elasticsearch.Config{
			Addresses: []string{os.Getenv("ES_HOST")}, // Ensure this is https://...
			Username:  os.Getenv("ES_USER"),
			Password:  os.Getenv("ES_PASSWORD"),
			// Add this to skip SSL verification for local self-signed certs
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		// FIX: Use = instead of := to assign to the global 'client' variable
		var err error
		client, err = elasticsearch.NewClient(cfg)
		if err != nil {
			log.Fatalf("Error creating Elasticsearch client: %s", err)
		}

		// Verify general cluster connection
		res, err := client.Info()
		if err != nil {
			log.Fatalf("Error connecting to Elasticsearch: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Fatalf("Elasticsearch returned an error: %s", res.String())
		}

		log.Println("Successfully connected to Elasticsearch cluster")

		// Check if the specific index exists
		// indexName := "golang125"
		// existsRes, err := client.Indices.Exists([]string{indexName})
		// if err != nil {
		// 	log.Printf("Error checking index existence: %s", err)
		// } else {
		// 	defer existsRes.Body.Close()
		// 	if existsRes.StatusCode == 200 {
		// 		log.Printf("Successfully connected to index: %s", indexName)
		// 	} else if existsRes.StatusCode == 404 {
		// 		log.Printf("Index %s does not exist", indexName)
		// 	}
		// }
	})
	return client
}
