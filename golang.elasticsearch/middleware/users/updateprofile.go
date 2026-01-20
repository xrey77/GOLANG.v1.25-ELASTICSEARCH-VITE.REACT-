package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
	utils "golang.elasticsearch/utils"

	"github.com/gin-gonic/gin"
)

// @Summary User Profile Update“
// @Description This will update user profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Param body body dto.ProfileData true "New Profile Details"
// @Success 200 {array} dto.ProfileData
// @Router /api/updateprofile/{id} [patch]
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var userDto dto.ProfileData

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	// 1. Verify user exists (using your existing utility)
	user, err := utils.GetUserid(id)
	if err != nil || len(user) == 0 {
		c.JSON(404, gin.H{"message": "User ID not found."})
		return
	}

	// 3. Prepare the update payload (Partial document update)
	updateData := map[string]interface{}{
		"doc": map[string]interface{}{
			"firstname": userDto.Firstname,
			"lastname":  userDto.Lastname,
			"mobile":    userDto.Mobile,
		},
	}
	payload, _ := json.Marshal(updateData)

	// 4. Execute the Elasticsearch Update request
	client := dbconfig.Connection()
	res, err := client.Update(
		"users", // Index name
		id,      // Document ID
		bytes.NewReader(payload),
		client.Update.WithContext(context.Background()),
	)

	if err != nil {
		c.JSON(500, gin.H{"message": "Error updating database"})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(500, gin.H{"message": fmt.Sprintf("Elasticsearch error: %s", res.String())})
		return
	}

	c.JSON(200, gin.H{"message": "Your profile has been updated successfully."})
}
