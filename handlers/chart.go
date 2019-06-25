package handlers

import (
	"bytes"
	"github.com/greenac/finance/analysis"
	"github.com/greenac/finance/logger"
	"github.com/wcharczuk/go-chart"
	"sort"
)


func DrawSingleAxisBinWithAmounts(bins *[]analysis.BinWithAmount) error {
	dBins := *(bins)

	sort.Slice(dBins, func(i, j int) bool {
		return dBins[i].Date.Before(dBins[j].Date)
	})

	b0 := dBins[0]
	xs := make([]float64, len(dBins))
	ys := make([]float64, len(dBins))

	for i, b := range dBins {
		s := b0.Date.Unix() - b.Date.Unix()
		xs[i] = float64(s)
		ys[i] = b.Amount
	}

	c := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xs,
				YValues: ys,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := c.Render(chart.PNG, buffer)
	if err != nil {
		logger.Log("`DrawSingleAxisBinWithAmounts` rendering chart:", err)
		return err
	}

	return nil
}
