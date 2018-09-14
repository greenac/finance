package models

import (
	"time"
)

const (
	ChaseDebitDetailsIndex     = 0
	ChaseDebitPostDateIndex    = 1
	ChaseDebitDescriptionIndex = 2
	ChaseDebitAmountIndex      = 3
	ChaseDebitTypeIndex        = 4
	ChaseDebitBalanceIndex     = 5
	ChaseDebitCheckNumberIndex = 6
)

type ChaseDebit struct {
	Details     string
	PostDate    time.Time
	Description string
	Amount      float64
	Type        string
	Balance     float64
	CheckNumber int
}
