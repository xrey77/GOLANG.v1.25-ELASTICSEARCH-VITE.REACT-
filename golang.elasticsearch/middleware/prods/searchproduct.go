package middleware

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Products Wild Cards Search
// @Description Product Search with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param page path int true "Page number"
// @Param key path string true "Key string"
// @Success 200 {array} dto.Products
// @Router /products/search/{page}/{key} [get]
func ProductSearch(c *gin.Context) {
	pageStr := c.Param("page")
	param1 := c.Param("key")
	key := "*" + strings.ToLower(param1) + "*"

	// 1. Setup Pagination variables
	perPage := 5

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	from := (page - 1) * perPage

	// 2. Initialize Elasticsearch client
	esClient := dbconfig.Connection()

	indexName := "products"

	// 3. Define the Search Query (Match query for descriptions)
	query := map[string]interface{}{
		"from": from,
		"size": perPage,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"descriptions": key,
				// "case_insensitive": true,
			},
		},
	}

	// Encode query to JSON
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(500, gin.H{"message": "Error encoding query"})
		return
	}

	// 4. Execute Search
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(strings.NewReader(buf.String())),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer res.Body.Close()

	// 5. Parse Response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(500, gin.H{"message": "Error parsing response"})
		return
	}

	// Extract hits and total records
	hits := r["hits"].(map[string]interface{})
	totalRecords := int64(hits["total"].(map[string]interface{})["value"].(float64))
	totalPages := math.Ceil(float64(totalRecords) / float64(perPage))

	// 6. Map results to your DTO
	var prods []dto.Products
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]

		// Convert map to your DTO struct
		var prod dto.Products
		byteData, _ := json.Marshal(source)
		json.Unmarshal(byteData, &prod)

		prods = append(prods, prod)
	}

	if len(prods) == 0 {
		c.JSON(404, gin.H{"message": "products not found."})
		return
	}

	c.JSON(200, gin.H{
		"page":         page,
		"totpage":      totalPages,
		"totalrecords": totalRecords,
		"products":     prods,
	})
}
