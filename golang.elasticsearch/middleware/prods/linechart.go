package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin" // Ensure you have this imported
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/models"
)

func GetLineChart(c *gin.Context) {
	esClient := dbconfig.Connection()
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("sales"),
		esClient.Search.WithSort("salesdate:asc"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ES search failed"})
		return
	}
	defer res.Body.Close()

	var response struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	monthKeys := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	salesMap := make(map[string]float64)

	for _, hit := range response.Hits.Hits {
		var sale models.Sale
		if err := json.Unmarshal(hit.Source, &sale); err == nil {
			month := sale.Salesdate.Format("Jan")
			salesMap[month] += sale.Amount

		}
	}

	var chartValues []chart.Value
	for _, m := range monthKeys {
		if val, ok := salesMap[m]; ok && val > 0 {
			chartValues = append(chartValues, chart.Value{
				Label: fmt.Sprintf("%s: %s", m, humanize.Commaf(val)),
				Value: val,
			})
		}
	}

	// Initialize PieChart instead of BarChart
	graph := chart.PieChart{
		Title:  "Monthly Sales Distribution",
		Height: 600,
		Width:  512, // Pie charts usually look better in square aspect ratios
		Values: chartValues,
		SliceStyle: chart.Style{
			FontColor:   drawing.ColorWhite, // Sets the label/amount text color
			FontSize:    10.0,               // Optional: adjust size as needed
			StrokeWidth: 2,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	// Render process remains identical
	if err := graph.Render(chart.PNG, buffer); err != nil {
		c.JSON(500, gin.H{"error": "Failed to render chart"})
		return
	}

	// ... [Keep the image composition/logo logic the same] ...

	// Note: Adjust combinedHeight to use graph.Height (512)
	chartImg, _, _ := image.Decode(buffer)
	logoFile, _ := os.Open("assets/images/logo.png")
	defer logoFile.Close()
	logoImg, _, _ := image.Decode(logoFile)

	logoHeight := 60
	gap := 10
	combinedHeight := logoHeight + gap + 600
	canvas := image.NewRGBA(image.Rect(0, 0, 512, combinedHeight))

	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	logoX := (512 - logoImg.Bounds().Dx()) / 2
	draw.Draw(canvas, image.Rect(logoX, 0, logoX+logoImg.Bounds().Dx(), logoHeight), logoImg, image.Point{}, draw.Over)
	draw.Draw(canvas, image.Rect(0, logoHeight+gap, 512, combinedHeight), chartImg, image.Point{}, draw.Over)

	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, canvas)

}

// func commaFormatter(v interface{}) string {
// 	p := message.NewPrinter(language.English)
// 	return p.Sprintf("%.2f", v)
// }
