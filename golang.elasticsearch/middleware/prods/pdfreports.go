package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/dto"
)

func ProductPDFReport(c *gin.Context) {
	cfg := config.NewBuilder().
		WithPageNumber(props.PageNumber{
			Pattern: "Page {current} of {total}",
			Place:   props.Bottom,
		}).
		WithLeftMargin(20).
		WithRightMargin(20).
		Build()

	m := maroto.New(cfg)

	err := m.RegisterFooter(
		text.NewRow(10, "Prepared by: Rey Gragasin", props.Text{
			Size:  6,
			Align: align.Left,
			Top:   5,
		}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Footer registration failed"})
		return
	}

	esClient := dbconfig.Connection()

	// Elasticsearch Query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 10000, // Max limit for standard search
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding query"})
		return
	}

	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("products"), // Ensure this matches your actual index
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search request failed", "details": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Elasticsearch error", "status": res.Status()})
		return
	}

	var response struct {
		Hits struct {
			Hits []struct {
				ID      string          `json:"_id"`
				Source_ json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decoding error", "details": err.Error()})
		return
	}

	var products []dto.Products
	for _, hit := range response.Hits.Hits {
		var product dto.Products
		if err := json.Unmarshal(hit.Source_, &product); err == nil {
			product.Id = hit.ID
			products = append(products, product)
		}
	}

	now := time.Now()
	imgBytes, _ := os.ReadFile("assets/images/logo.png")

	// Header Rows
	m.AddRows(
		row.New(20).Add(
			image.NewFromBytesCol(12, imgBytes, extension.Png, props.Rect{
				Center:  true,
				Percent: 80,
			}),
		),
		row.New(15).Add(
			text.NewCol(12, "Product List Report", props.Text{
				Top:   5,
				Size:  20,
				Style: fontstyle.Bold,
				Align: align.Center,
			}),
		),
		row.New(12).Add(
			text.NewCol(12, "As of "+now.Format("January 02, 2006"), props.Text{
				Size:  10,
				Style: fontstyle.Normal,
				Align: align.Center,
			}),
		),
	)

	// Table Header
	m.AddRows(
		row.New(10).Add(
			text.NewCol(1, "ID", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(5, "Product Descriptions", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(2, "Stocks", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(2, "Cost", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(2, "Sell", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		),
	)

	// Data Rows
	ln := 0
	for _, p := range products {
		ln++
		m.AddRows(row.New(8).Add(
			text.NewCol(1, fmt.Sprintf("%d", ln), props.Text{Align: align.Center}),
			text.NewCol(5, p.Descriptions),
			text.NewCol(2, fmt.Sprintf("%.0f", p.Qty), props.Text{Align: align.Center}),
			text.NewCol(2, fmt.Sprintf("%.2f", p.Costprice), props.Text{Align: align.Center}),
			text.NewCol(2, fmt.Sprintf("%.2f", p.Sellprice), props.Text{Align: align.Center}),
		))
	}

	// Document Generation
	doc, err := m.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed"})
		return
	}

	pdfBytes := doc.GetBytes()
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=products_report.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
