package analysis

import (
	"github.com/greenac/finance/csv"
	"github.com/greenac/finance/json"
	"strings"
)

const Random = "random"

type Summation map[string]float64

type Analyzer struct {
	Datums *[]csv.Datum
	Groups *json.Group
	bins   map[string][]csv.Datum
}

func (a *Analyzer) GroupByType() {
	a.bins = make(map[string][]csv.Datum)
	for _, dat := range *a.Datums {
		if dat.Type == "payment" {
			continue
		}

		has := false
		for group, entries := range *a.Groups {
			for _, ent := range entries {
				if strings.Contains(strings.ToLower(dat.Description), strings.ToLower(ent)) {
					a.add(dat, group)
					has = true
				}
			}
		}

		if !has {
			a.add(dat, Random)
		}
	}
}

func (a *Analyzer) add(d csv.Datum, group string) {
	var dts []csv.Datum
	if dats, has := a.bins[group]; has {
		dts = append(dats, d)
	} else {
		dts = []csv.Datum{d}
	}

	a.bins[group] = dts
}

func (a *Analyzer) Sum() *Summation {
	sum := make(Summation)
	for grp, dts := range a.bins {
		var t float64 = 0
		for _, d := range dts {
			t += d.Amount
		}

		sum[grp] = t
	}

	return &sum
}

func (a *Analyzer) Bin(b string) *[]csv.Datum {
	dats, has := a.bins[b]
	if !has {
		return nil
	}

	return &dats
}
