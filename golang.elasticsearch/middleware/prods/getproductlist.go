package middleware

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	dbconfig "golang.elasticsearch/dbconfig"
)

// @Summary Product Listings
// @Description Products Pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param page path int true "Page number"
// @Success 200 {array} []dto.Products
// @Router /products/list/:page [get]
func GetProductList(c *gin.Context) {
	//	pageStr := c.DefaultQuery("page", "1")
	pageStr := c.Param("page")

	pg, _ := strconv.Atoi(pageStr)
	if pg < 1 {
		pg = 1
	}

	perPage := 5
	offset := (pg - 1) * perPage

	esClient := dbconfig.Connection()

	res, err := esClient.Search(
		esClient.Search.WithIndex("products"),
		esClient.Search.WithFrom(offset),
		esClient.Search.WithSize(perPage),
		esClient.Search.WithContext(c.Request.Context()),
		esClient.Search.WithSort("descriptions.keyword:asc"),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "ES Search failed"})
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(500, gin.H{"error": "Decode failed"})
		return
	}

	hitsObj, ok := r["hits"].(map[string]interface{})
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid ES response structure"})
		return
	}

	// 1. Extract the raw hits slice
	hitsList, ok := hitsObj["hits"].([]interface{})
	if !ok {
		c.JSON(500, gin.H{"error": "No hits found"})
		return
	}

	counter := offset + 1

	// 2. Create a new slice to store only the _source data
	var products []interface{}
	for _, hit := range hitsList {
		if h, ok := hit.(map[string]interface{}); ok {
			if source, exists := h["_source"].(map[string]interface{}); exists {
				// if source, exists := h["_source"]; exists {
				source["id"] = counter
				products = append(products, source)
				counter++

			}
		}
	}

	totalVal := hitsObj["total"].(map[string]interface{})["value"].(float64)
	totalPages := math.Ceil(totalVal / float64(perPage))

	c.JSON(200, gin.H{
		"page":         pg,
		"totpage":      totalPages,
		"totalrecords": totalVal,
		"products":     products, // Now contains only the _source objects
	})

}
