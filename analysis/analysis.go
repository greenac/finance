package analysis

import (
	"errors"
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

func (bbd *BinsByDate) totalAmount() float64 {
	amount := 0.0
	for _, b := range bbd.Models {
		a := b.DebitedAmount()
		if a <= 0 {
			continue
		}

		amount += b.DebitedAmount()
	}

	return amount
}


type BinsWithCountByDate struct {
	Date  time.Time
	Count int
}


type BinWithAmount struct {
	Date time.Time
	Amount float64
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

			for i := 0; i < len(binsAndDates); i+=1 {
				bbd := &(binsAndDates[i])
				ed := utils.EndOfDay(bbd.Date)
				if (td.After(bbd.Date) && td.Before(ed)) || td.Equal(bbd.Date) {
					bbd.Models = append(bbd.Models, m)
					break
				}

				if i == len(binsAndDates)-1 {
					binsAndDates = append(binsAndDates, BinsByDate{utils.StartOfDay(td), []models.CsvModel{m}})
					break
				}
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

func (a *Analyzer) BinsWithAmounts() *[]BinWithAmount {
	bins := make([]BinWithAmount, len(*(a.binsAndDates)))
	for i, b := range *(a.binsAndDates) {
		bins[i] = BinWithAmount{Date: b.Date, Amount: b.totalAmount()}
	}

	return &bins
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
