package analysis

import (
	"errors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/finance/json"
	"github.com/greenac/finance/models"
	"github.com/greenac/finance/utils"
	"strings"
	"time"
)

const Random = "random"

type Summation map[string]float64
type bins map[string][]models.CsvModel

type BinsByDate struct {
	Date   time.Time
	Models []models.CsvModel
}

type BinsWithCountByDate struct {
	Date  time.Time
	Count int
}

type Analyzer struct {
	Models       *models.ModelsByType
	Groups       *json.Group
	bins         bins
	binsAndDates *[]BinsByDate
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

func (a *Analyzer) GroupByDate() {
	a.bins = make(bins)
	binsAndDates := make([]BinsByDate, 0)

	for _, mods := range *a.Models {
		for _, m := range *mods {
			if m.TransType() == models.Deposit {
				continue
			}

			td := m.TransDate()

			if len(binsAndDates) == 0 {
				binsAndDates = append(binsAndDates, BinsByDate{utils.StartOfDay(td), []models.CsvModel{m}})
				continue
			}

			add := false
			for i := 0; i < len(binsAndDates); i++ {
				bbd := &(binsAndDates[i])
				ed := utils.EndOfDay(bbd.Date)
				if (td.After(bbd.Date) && td.Before(ed)) || td.Equal(bbd.Date) {
					bbd.Models = append(bbd.Models, m)
					logger.Log("Appening model", bbd.Date, "td:", td, "end date:", ed, "model count", len(bbd.Models))
					break
				}

				if i == len(binsAndDates)-1 {
					add = true
					break
				}
			}

			if add {
				logger.Error("Adding new model start:", utils.StartOfDay(td), "td:", td)
				binsAndDates = append(binsAndDates, BinsByDate{utils.StartOfDay(td), []models.CsvModel{m}})
			}
		}
	}

	a.binsAndDates = &binsAndDates
}

func (a *Analyzer) CountsByDate() *[]BinsWithCountByDate {
	a.GroupByDate()
	bins := *(a.binsAndDates)
	binsWithCounts := make([]BinsWithCountByDate, len(bins))

	for i, b := range bins {
		bwc := BinsWithCountByDate{Date: b.Date, Count: len(b.Models)}
		binsWithCounts[i] = bwc
	}

	return &binsWithCounts
}

func (a *Analyzer) BinsCountedByDate() (*[]BinsByDate, error) {
	if a.binsAndDates == nil {
		return nil, errors.New("BINS_AND_DATES_NOT_SET")
	}

	return a.binsAndDates, nil
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
