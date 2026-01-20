package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dbconfig "golang.elasticsearch/dbconfig"
	utils "golang.elasticsearch/utils"

	"github.com/gin-gonic/gin"
	"golang.elasticsearch/dto"
	"golang.org/x/crypto/bcrypt"
)

// @Summary User Login
// @Description Authenticat User
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.UserLogin true "User Login Credentials"
// @Success 200 {array} dto.UserLogin
// @Router /auth/signin [post]
func Login(c *gin.Context) {
	var userDto dto.UserLogin

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}
	plainPwd := userDto.Password
	user, err := GetUserInfo(userDto.Username)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if user == nil {

		c.JSON(404, gin.H{"message": "Username not found, please register."})
		return
	} else {

		hashPwd := user.Password
		err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plainPwd))
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid Password."})
			return
		} else {

			token, _ := utils.GenerateJWT(user.Email)

			c.JSON(200, gin.H{
				"id":          user.Id,
				"firstname":   user.Firstname,
				"lastname":    user.Lastname,
				"email":       user.Email,
				"mobile":      user.Mobile,
				"username":    user.Username,
				"roles":       user.Roles,
				"isactivated": user.Isactivated,
				"isblocked":   user.Isblocked,
				"userpicture": user.Userpicture,
				"qrcodeurl":   user.Qrcodeurl,
				"token":       token,
				"message":     "Login Successfull."})
		}

	}
}

func GetUserInfo(userName string) (*dto.Users, error) {
	client := dbconfig.Connection()

	// 1. Build the search query
	// Use "username.keyword" if your field is type 'text' to ensure an exact match
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"term": map[string]interface{}{
	// 			"username.keyword": userName,
	// 		},
	// 	},
	// }

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"username": userName, // Use the base field, not .keyword
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// 2. Execute the Search
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("users"),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching user: %s", res.String())
	}

	// 3. Parse the Response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// 4. Extract hits
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

	if len(hits) == 0 {
		return nil, nil // User not found
	}

	hit := hits[0].(map[string]interface{})

	// 2. Extract the system ID (which is at the top level of the hit)
	id := hit["_id"].(string)

	// 5. Convert the first hit to your DTO
	firstHit := hits[0].(map[string]interface{})["_source"]
	sourceData, _ := json.Marshal(firstHit)

	var user dto.Users
	if err := json.Unmarshal(sourceData, &user); err != nil {
		return nil, err
	}
	user.Id = id

	return &user, nil
}
