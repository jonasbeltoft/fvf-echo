package main

type ListPageData struct {
	IsHtmx  bool
	Page    string
	Records [][]string
}

type Data struct {
	LineCharts []LineChart
	IsHtmx     bool
	Page       string
}

type LineChart struct {
	Title string
	Lines DualLines
	ID    string
}

type DualLines struct {
	Prediction LineData
	Actual     LineData
	Marker     Marker
}

type LineData struct {
	X []string
	Y []float64
}

type Marker struct {
	X string
	Y float64
}
