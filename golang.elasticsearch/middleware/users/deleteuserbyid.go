package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
)

// DeleteUserid godoc
// @Summary Delete user by ID
// @Description Find and Delete a single user's details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Success 200 {object} DeleteResponse
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/deleteuserbyid/{id} [delete]
func DeleteUserid(c *gin.Context) {
	id := c.Param("id")
	esClient := dbconfig.Connection()

	// 1. Direct Delete call using ID
	// No search query is needed for a simple ID-based deletion.
	res, err := esClient.Delete(
		"users",
		id,
		esClient.Delete.WithContext(context.Background()),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute delete"})
		return
	}
	defer res.Body.Close()

	// 2. Check for HTTP errors (e.g., 404 Not Found)
	if res.IsError() {
		if res.StatusCode == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Elasticsearch error"})
		return
	}

	// 3. Optional: Parse response for confirmation
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing ES response"})
		return
	}

	// Elasticsearch returns "result": "deleted" on success
	c.JSON(http.StatusOK, gin.H{
		"message": "User has been deleted successfully",
		"result":  r["result"],
	})
}
