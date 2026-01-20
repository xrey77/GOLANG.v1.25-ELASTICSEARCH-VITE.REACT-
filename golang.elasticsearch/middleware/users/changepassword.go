package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	utils "golang.elasticsearch/utils"

	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
)

// @Summary Change User Password
// @Description User Change Password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Param body body dto.ChangePassword true "New Password Details"
// @Success 200 {object} dto.ChangePassword
// @Router /api/changepassword/{id} [patch]
func ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var userDto dto.ChangePassword

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

	// 2. Hash the new password
	hash, _ := utils.HashPassword(userDto.Password)

	// 3. Prepare the update payload (Partial document update)
	updateData := map[string]interface{}{
		"doc": map[string]interface{}{
			"password":    hash,
			"isactivated": true,
			"userpicture": "pix.png",
			"mailtoken":   "0",
		},
	}
	payload, _ := json.Marshal(updateData)

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

	c.JSON(200, gin.H{"message": "Password has been changed."})
}
