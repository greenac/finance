package models

import (
	"time"
)

const (
	ChaseCreditTypeIndex        = 0
	ChaseCreditDateIndex        = 1
	ChaseCreditPostDateIndex    = 2
	ChaseCreditDescriptionIndex = 3
	ChaseCreditAmountIndex      = 4
)

type ChaseCredit struct {
	Type        string
	Date        time.Time
	PostDate    time.Time
	Description string
	Amount      float64
}
