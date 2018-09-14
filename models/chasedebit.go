package models

import (
	"errors"
	"github.com/greenac/finance/logger"
	"strconv"
	"time"
)

const (
	chaseDebitDetailsIndex     = 0
	chaseDebitPostDateIndex    = 1
	chaseDebitDescriptionIndex = 2
	chaseDebitAmountIndex      = 3
	chaseDebitTypeIndex        = 4
	chaseDebitBalanceIndex     = 5
	chaseDebitCheckNumberIndex = 6
)

const numCDEntries = 7

type ChaseDebit struct {
	Details     string
	PostDate    time.Time
	Description string
	Amount      float64
	Type        string
	Balance     float64
	CheckNumber int64
}

func (cd *ChaseDebit) SetValues(entries []string) error {
	if len(entries) != numCDEntries {
		logger.Error("`ChaseDebit::SetValues` Invalid number of entries:", len(entries), "should be:", numCDEntries)
		return errors.New("INVALID_NUM_ENTRIES")
	}

	pd, err := time.Parse(DateLayout, entries[chaseDebitPostDateIndex])
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing date:", err)
		return errors.New("INVALID_ENTRY")
	}

	pDate, err := time.Parse(DateLayout, entries[chaseDebitPostDateIndex])
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing post date:", err)
		return errors.New("INVALID_ENTRY")
	}

	amt, err := strconv.ParseFloat(entries[chaseDebitAmountIndex], 64)
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing amount:", err)
		return errors.New("INVALID_ENTRY")
	}

	bal, err := strconv.ParseFloat(entries[chaseDebitBalanceIndex], 64)
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing balance:", err)
		return errors.New("INVALID_ENTRY")
	}

	cn, err := strconv.ParseInt(entries[chaseDebitCheckNumberIndex], 10, 64)
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing check number:", err)
		return errors.New("INVALID_ENTRY")
	}

	cd.Details = entries[chaseDebitDetailsIndex]
	cd.PostDate = pd
	cd.Description = entries[chaseDebitDescriptionIndex]
	cd.Amount = amt
	cd.Type = entries[chaseDebitTypeIndex]
	cd.Balance = bal
	cd.PostDate = pDate
	cd.CheckNumber = cn

	return nil
}

func (cd *ChaseDebit) DebitedAmount () float64 {
	return -cd.Amount
}
