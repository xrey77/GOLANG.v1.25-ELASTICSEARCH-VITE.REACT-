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
	pageStr := c.DefaultQuery("page", "1")
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

	totalVal := hitsObj["total"].(map[string]interface{})["value"].(float64)
	totalPages := math.Ceil(totalVal / float64(perPage))

	c.JSON(200, gin.H{
		"page":         pg,
		"totpage":      totalPages,
		"totalrecords": totalVal,
		"products":     hitsObj["hits"],
	})
}
