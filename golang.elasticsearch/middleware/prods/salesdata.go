package middleware

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart/v2"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetSalesChart(c *gin.Context) {
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
		chartValues = append(chartValues, chart.Value{
			Label: m,
			Value: salesMap[m], // Defaults to 0.0 if month is missing in salesMap
		})
	}

	graph := chart.BarChart{
		Title: "Monthly Sales Report",
		Background: chart.Style{
			Padding: chart.Box{Top: 20, Bottom: 20},
		},
		Height:   512,
		Width:    1280,
		BarWidth: 60,
		Bars:     chartValues,
		YAxis: chart.YAxis{
			// Style: chart.Style{
			// 	TextRotationDegrees: -90.0, // Rotates labels sideways
			// },
			ValueFormatter: commaFormatter,
		},
	}

	// 2. Render chart to a temporary buffer
	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		c.JSON(500, gin.H{"error": "Failed to render chart"})
		return
	}
	chartImg, _, _ := image.Decode(buffer)

	// 3. Load your Logo
	logoFile, _ := os.Open("assets/images/logo.png") // Ensure path is correct
	defer logoFile.Close()
	logoImg, _, _ := image.Decode(logoFile)

	// 4. Create a combined canvas
	// Height = Logo Height + Gap + Chart Height
	logoHeight := 60
	gap := 10
	combinedHeight := logoHeight + gap + graph.Height
	canvas := image.NewRGBA(image.Rect(0, 0, graph.Width, combinedHeight))

	// Draw Background (White)
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 5. Place Logo (centered at the top)
	logoX := (graph.Width - logoImg.Bounds().Dx()) / 2
	draw.Draw(canvas, image.Rect(logoX, 0, logoX+logoImg.Bounds().Dx(), logoHeight), logoImg, image.Point{}, draw.Over)

	// 6. Place Chart (below the logo)
	draw.Draw(canvas, image.Rect(0, logoHeight+gap, graph.Width, combinedHeight), chartImg, image.Point{}, draw.Over)

	// 7. Write final image to response
	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, canvas)
}

func commaFormatter(v interface{}) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", v)
}
