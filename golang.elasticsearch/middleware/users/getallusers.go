package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve users
// @Description Display all users
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Users
// @Router /api/getallusers [get]
func GetAllUsers(c *gin.Context) {
	client := dbconfig.Connection()

	// Define a "match all" query for Elasticsearch
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 10000, // Important: Set a size limit to retrieve more than the default 10 records
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		// Handle the error appropriately, e.g., return from the function
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding query"})
		return
	}

	// Perform the search request
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("users"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
	)

	// Note: You removed the `Do(context.Background())` call as it was a syntax error.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search request failed", "details": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		// Read the response body for a better error message if needed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Elasticsearch error", "status": res.Status()})
		return
	}

	// Decode the response body
	var response struct {
		Hits struct {
			Hits []struct {
				ID      string          `json:"_id"`
				Source_ json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding Elasticsearch response", "details": err.Error()})
		return
	}

	var users []dto.Users
	for _, hit := range response.Hits.Hits {
		var user dto.Users
		if err := json.Unmarshal(hit.Source_, &user); err == nil {
			user.Id = hit.ID
			users = append(users, user)
		} else {
			// Log the unmarshal error if needed
			log.Printf("Error unmarshaling user data: %s\n", err)
		}
	}

	c.JSON(http.StatusOK, users)
}
