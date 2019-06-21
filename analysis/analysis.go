package analysis

import (
	"github.com/greenac/finance/json"
	"github.com/greenac/finance/models"
	"github.com/greenac/finance/utils"
	"strings"
	"time"
)

const Random = "random"

type Summation map[string]float64
type bins map[string][]models.CsvModel

type binsByDate struct {
	date time.Time
	models []models.CsvModel
}

type Analyzer struct {
	Models *models.ModelsByType
	Groups *json.Group
	bins   bins
	binsAndDates binsByDate
}

func (a *Analyzer) GroupByType() {
	a.bins = make(bins)
	for _, mods := range *a.Models {
		for _, m := range *mods {
			if m.TransType() == models.Deposit {
				continue
			}

			has := false
			for group, entries := range *a.Groups {
				for _, ent := range entries {
					if strings.Contains(strings.ToLower(m.Desc()), strings.ToLower(ent)) {
						a.add(m, group)
						has = true
					}
				}
			}

			if !has {
				a.add(m, Random)
			}
		}
	}
}

func (a *Analyzer) GroupByDate(binSize int) {
	a.bins = make(bins)
	binsAndDates := make([]binsByDate, 0)

	for _, mods := range *a.Models {
		for _, m := range *mods {
			if m.TransType() == models.Deposit {
				continue
			}

			sd := utils.StartOfDay(m.)
			if ()
		}
	}
}

func (a *Analyzer) add(m models.CsvModel, group string) {
	var mods []models.CsvModel
	if ms, has := a.bins[group]; has {
		mods = append(ms, m)
	} else {
		mods = []models.CsvModel{m}
	}

	a.bins[group] = mods
}

func (a *Analyzer) Sum() *Summation {
	sum := make(Summation)
	for grp, dts := range a.bins {
		var t float64 = 0
		for _, d := range dts {
			t += d.DebitedAmount()
		}

		sum[grp] = t
	}

	return &sum
}

func (a *Analyzer) Bin(b string) *[]models.CsvModel {
	mods, has := a.bins[b]
	if !has {
		return nil
	}

	return &mods
}
