package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	config "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
)

func GetUserid(id string) ([]dto.Users, error) {
	esClient := config.Connection()

	// 1. Define the query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": []string{id},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// 2. Execute Search
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("users"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from ES: %s", res.String())
	}

	// 3. Define a temporary structure to parse the ES response
	var r struct {
		Hits struct {
			Hits []struct {
				ID     string    `json:"_id"`
				Source dto.Users `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	// 4. Decode directly into the struct
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// 5. Build the result slice
	var users []dto.Users
	for _, hit := range r.Hits.Hits {
		user := hit.Source
		user.Id = hit.ID // Assign the document ID to the struct field
		users = append(users, user)
	}

	return users, nil
}
