package main

import (
	"encoding/csv"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setup for the template renderer
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Mock data
func newData() Data {
	return Data{
		LineCharts: []LineChart{
			{
				Title: "First Chart",
				Lines: DualLines{
					Prediction: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00", "2024-03-15 05:00", "2024-03-15 06:00", "2024-03-15 07:00", "2024-03-15 08:00", "2024-03-15 09:00"},
						Y: []float64{50, 52, 55, 58, 60, 62, 65, 68, 70, 72},
					},
					Actual: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00"},
						Y: []float64{46, 51, 55, 56, 62},
					},
				},
				ID: "first-chart",
			},
			{
				Title: "Second Chart",
				Lines: DualLines{
					Prediction: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00", "2024-03-15 05:00", "2024-03-15 06:00", "2024-03-15 07:00", "2024-03-15 08:00", "2024-03-15 09:00"},
						Y: []float64{50, 55, 51, 56, 53, 51, 55, 59, 54, 58},
					},
					Actual: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00"},
						Y: []float64{46, 51, 55, 54, 51},
					},
				},
				ID: "second-chart",
			},
			{
				Title: "Third Chart",
				Lines: DualLines{
					Prediction: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00", "2024-03-15 05:00", "2024-03-15 06:00", "2024-03-15 07:00", "2024-03-15 08:00", "2024-03-15 09:00"},
						Y: []float64{40, 37, 39, 42, 38, 43, 47, 45, 41, 44},
					},
					Actual: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00"},
						Y: []float64{46, 42, 40, 41, 36},
					},
				},
				ID: "third-chart",
			},
			{
				Title: "Fourth Chart",
				Lines: DualLines{
					Prediction: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00", "2024-03-15 05:00", "2024-03-15 06:00", "2024-03-15 07:00", "2024-03-15 08:00", "2024-03-15 09:00"},
						Y: []float64{80, 77, 79, 83, 85, 82, 79, 76, 78, 81},
					},
					Actual: LineData{
						X: []string{"2024-03-15 00:00", "2024-03-15 01:00", "2024-03-15 02:00", "2024-03-15 03:00", "2024-03-15 04:00"},
						Y: []float64{74, 78, 81, 88, 92},
					},
				},
				ID: "fourth-chart",
			},
		},
	}
}

// Returns array of arrays of strings from the CSV file
// Aray of arrays of strings is used to represent the CSV data, line by line, cell by cell
// Each line is an array of strings, each cell is a string. [4][2] is fifth line, third cell
func readDatabaseData() ([][]string, error) {
	// Open data/test_data.csv
	file, err := os.Open("data/test_data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	// Create a new Echo instance
	e := echo.New()
	// Use the Logger middleware
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	// Serve all static files in the "static" folder to '/assets'
	e.Static("/static", "static")
	e.File("/favicon.png", "images/favicon.png")

	// Define routes
	e.GET("/", root)
	e.GET("/dashboard", dashboard)
	e.GET("/list", list)

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8000"))
}

// Route handlers

func root(c echo.Context) error {
	data := newData()
	data.IsHtmx = c.Request().Header.Get("HX-Request") == "true"
	return c.Render(200, "index", data)
}

func dashboard(c echo.Context) error {
	data := newData()
	data.IsHtmx = c.Request().Header.Get("HX-Request") == "true"
	code := 200
	for i := range data.LineCharts {
		marker := Marker{
			X: data.LineCharts[i].Lines.Actual.X[len(data.LineCharts[i].Lines.Actual.X)-1],
			Y: data.LineCharts[i].Lines.Actual.Y[len(data.LineCharts[i].Lines.Actual.Y)-1],
		}
		data.LineCharts[i].Lines.Marker = marker
	}
	return c.Render(code, "dashboard", data)
}

func list(c echo.Context) error {
	data, err := readDatabaseData()
	if err != nil {
		return c.JSON(500, "Error reading data")
	}
	listPageData := ListPageData{
		IsHtmx:  c.Request().Header.Get("HX-Request") == "true",
		Page:    "List",
		Records: data,
	}
	return c.Render(200, "listPage", listPageData)
}
