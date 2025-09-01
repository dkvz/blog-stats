package stats

import (
	"io"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Functions here are more generic versions
// of what runtime will be using.

// We need to be able to generate scatter plots
// And the runtime decides whether to generate
// images or run a local HTTP server.

func GenerateScatterPlot(data map[float64]float64, title string) components.Charter {
	scatter := charts.NewScatter()

	scatter.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	chartData := make([]opts.ScatterData, 0, len(data))
	for x, y := range data {
		chartData = append(chartData, opts.ScatterData{
			Value:      []float64{x, y},
			SymbolSize: 10,
		})
	}

	scatter.AddSeries("data", chartData).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show: opts.Bool(false),
			}))

	return scatter
}

func GeneratePlotPage(w io.Writer, charts ...components.Charter) {
	page := components.NewPage()
	page.AddCharts(charts...)
	page.Render(w)
}
