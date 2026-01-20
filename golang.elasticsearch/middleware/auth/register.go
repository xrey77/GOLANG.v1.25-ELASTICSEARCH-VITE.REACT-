package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
	"golang.elasticsearch/models"
	"golang.elasticsearch/utils"
)

// @Summary User Registration
// @Description Create User Account
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.UserRegister true "Account Registration"
// @Success 200 {array} dto.UserRegister
// @Router /auth/signup [post]
func Register(c *gin.Context) {
	var userDto dto.UserRegister
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	hashPwd, _ := utils.HashPassword(userDto.Password)
	client := dbconfig.Connection()
	indexName := "users"

	// 1. Check if index exists; if not, create it with your mapping
	exists, _ := client.Indices.Exists([]string{indexName})
	if exists.StatusCode == 404 {
		mapping := `{ "mappings": { "properties": { ... } } }` // your mapping here
		client.Indices.Create(indexName, client.Indices.Create.WithBody(strings.NewReader(mapping)))
	}

	userDto.Email = strings.ToLower(userDto.Email)
	userEmail, _ := SearchByEmail(userDto.Email)
	if len(userEmail) > 0 {
		c.JSON(400, gin.H{
			"message": "Email Address is already taken."})
		return
	}

	userName, _ := SearchByUsername(userDto.Username)
	if len(userName) > 0 {
		c.JSON(400, gin.H{"message": "Username is already taken."})
		return
	}

	// 2. Prepare userModel and marshal to JSON
	userModel := &models.User{
		Firstname:   userDto.Firstname,
		Lastname:    userDto.Lastname,
		Email:       userDto.Email,
		Mobile:      userDto.Mobile,
		Username:    userDto.Username,
		Password:    hashPwd,
		Roles:       "ROLE_USER",
		Isactivated: true,
		Userpicture: "pix.png",
		Mailtoken:   0,
		Secret:      nil,
		Qrcodeurl:   nil,
	}

	data, err := json.Marshal(userModel)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error marshaling data"})
		return
	}

	// 3. Insert the document
	res, err := client.Index(
		indexName,
		bytes.NewReader(data),
		client.Index.WithRefresh("wait_for"), // Optional: ensures data is searchable immediately
	)

	if err != nil || res.IsError() {
		c.JSON(500, gin.H{"message": "Failed to index user"})
		return
	}

	var esResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		c.JSON(500, gin.H{"message": "Error parsing Elasticsearch response"})
		return
	}

	// Safely extract the ID
	createdID := esResult["_id"].(string)

	c.JSON(201, gin.H{
		"message": "You have registered successfully, your user ID Is " + createdID,
	})

}

func SearchByEmail(email string) ([]models.User, error) {
	var users []models.User
	client := dbconfig.Connection()

	// FIX 1: Ensure you are searching the SAME index used in Register
	indexName := "users"

	// FIX 2: Use term query with .keyword for exact matches
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"email.keyword": email,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{indexName}, // Updated index name
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// FIX 3: Safe extraction of hits
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return users, nil
	}

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		var user models.User
		sourceBytes, _ := json.Marshal(source)
		if err := json.Unmarshal(sourceBytes, &user); err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}

//

func SearchByUsername(username string) ([]models.User, error) {
	var users []models.User
	client := dbconfig.Connection()

	// FIX 1: Ensure you are searching the SAME index used in Register
	indexName := "users"

	// FIX 2: Use term query with .keyword for exact matches
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"username.keyword": username,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{indexName}, // Updated index name
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// FIX 3: Safe extraction of hits
	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return users, nil
	}

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		var user models.User
		sourceBytes, _ := json.Marshal(source)
		if err := json.Unmarshal(sourceBytes, &user); err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}
