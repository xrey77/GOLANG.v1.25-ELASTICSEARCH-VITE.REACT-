package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	dbconfig "golang.elasticsearch/dbconfig"
	utils "golang.elasticsearch/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Update user profile picture
// @Description Upload user picture
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Param userpic formData file true "New Profile Picture"
// @Success 200 {object} map[string]interface{}
// @Router /api/uploadpicture/{id} [patch]
func UploadPicture(c *gin.Context) {
	id := c.Param("id")
	user, err := utils.GetUserid(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if len(user) > 0 {

		file, err := c.FormFile("userpic") // "file" is the key for the form data
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		newfile := "00" + id + ext
		dst := filepath.Join("./assets/users/", newfile) // Destination path

		// Save the uploaded file to the specified destination
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		updateData := map[string]interface{}{
			"doc": map[string]interface{}{
				"userpicture": newfile,
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

		c.JSON(200, gin.H{
			"userpic": newfile,
			"message": "Profile picture has been changed."})

	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}

}
