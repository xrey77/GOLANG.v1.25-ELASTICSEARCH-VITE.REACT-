package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Get user by ID
// @Description Retrieve a single user's details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Success 200 {object} dto.Users
// @Router /api/getuserid/{id} [get]
func GetUserid(c *gin.Context) {
	id := c.Param("id")
	client := dbconfig.Connection()

	// 1. Build Query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": id, // Use _id to filter by Elasticsearch document ID
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encoding error"})
		return
	}

	// 2. Execute Search
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("users"),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ES search error"})
		return
	}

	// 3. Parse Response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decoding error"})
		return
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// 4. Extract Data
	hit := hits[0].(map[string]interface{})
	source := hit["_source"].(map[string]interface{})

	// Convert map to struct
	sourceData, _ := json.Marshal(source)
	var user dto.Users
	if err := json.Unmarshal(sourceData, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Mapping error"})
		return
	}

	// Assign the ES document ID to the struct
	user.Id = hit["_id"].(string)

	// 5. Send Response
	c.JSON(http.StatusOK, user)
}
