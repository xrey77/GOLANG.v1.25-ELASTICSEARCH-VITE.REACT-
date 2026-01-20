package middleware

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
	"golang.elasticsearch/models"
)

func AddSalesData(c *gin.Context) {
	var salesDto dto.Sales

	if err := c.ShouldBindJSON(&salesDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}
	esClient := dbconfig.Connection()
	indexName := "sales"

	layout := "2006-01-02"
	xdate, err := time.Parse(layout, salesDto.Salesdate)
	if err != nil {
		return
	}

	saleModel := &models.Sale{
		Amount:    salesDto.Amount,
		Salesdate: xdate,
	}

	data, err := json.Marshal(saleModel)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error marshaling data"})
		return
	}

	// 3. Insert the document
	res, err := esClient.Index(
		indexName,
		bytes.NewReader(data),
		esClient.Index.WithRefresh("wait_for"), // Optional: ensures data is searchable immediately
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

	c.JSON(201, gin.H{
		"message": "New Sales has been added successfully. ",
	})

}
