package middleware

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
	"golang.elasticsearch/models"
)

// @Summary Add New Product
// @Description Create a new product in the system
// @Tags Products
// @Accept json
// @Produce json
// @Param product body dto.Products true "Product object data"
// @Success 201 {object} map[string]interface{} "Successfully created"
// @Failure 400 {object} map[string]interface{} "Invalid request format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/addproduct [post]
func AddProduct(c *gin.Context) {
	var productDto dto.Products

	if err := c.ShouldBindJSON(&productDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	esClient := dbconfig.Connection()
	indexName := "products"

	productModel := &models.Product{
		Category:       productDto.Category,
		Descriptions:   productDto.Descriptions,
		Qty:            productDto.Qty,
		Unit:           productDto.Unit,
		Costprice:      productDto.Costprice,
		Sellprice:      productDto.Sellprice,
		Saleprice:      productDto.Saleprice,
		Productpicture: productDto.Productpicture,
		Alertstocks:    productDto.Alertstocks,
		Criticalstocks: productDto.Criticalstocks,
	}

	data, err := json.Marshal(productModel)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error marshaling data"})
		return
	}

	// Perform the Indexing
	res, err := esClient.Index(
		indexName,
		bytes.NewReader(data),
		esClient.Index.WithRefresh("wait_for"),
	)

	// Handle Elasticsearch Errors
	if err != nil {
		c.JSON(500, gin.H{"message": "Elasticsearch connection error"})
		return
	}

	// Always close the body when res is not nil
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(500, gin.H{"message": "Failed to index product in Elasticsearch"})
		return
	}

	// Only send the success response ONCE at the very end
	c.JSON(201, gin.H{
		"message": "New product has been added successfully.",
	})
}
