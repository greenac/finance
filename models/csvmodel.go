package models

import (
	"strconv"
	"time"
)

const DateLayout = "01/02/2006"

type Name string

const (
	ChaseCreditName  Name = "chaseCredit"
	ChaseDebitName   Name = "chaseDebit"
	CapOneCreditName Name = "capOneCredit"
)

type TransType string

const (
	Withdrawal TransType = "withdrawal"
	Deposit    TransType = "deposit"
)

type CsvModel interface {
	SetValues([]string) error
	DebitedAmount() float64
	Desc() string
	TransType() TransType
	TransDate() time.Time
}

type CSVDirPaths map[Name]string
type ModelsByType map[Name]*[]CsvModel
type CSVFilePaths map[Name][]string

func cleanEntry(ln []string, num int) []string {
	if len(ln) <= num {
		return ln
	}

	newLn := make([]string, num)
	for i := 0; i < num; i++ {
		newLn[i] = ln[i]
	}

	return newLn
}

func handleParseFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if s == "" || s == " " {
			v = 0
		} else {
			return 0, err
		}
	}

	return v, nil
}

func EntriesForAccount(n Name) int {
	var ent int
	switch n {
	case ChaseCreditName:
		ent = NumChaseCreditEntries
	case ChaseDebitName:
		ent = NumChaseDebitEntries
	case CapOneCreditName:
		ent = NumCapOneCreditEntries
	}

	return ent
}
