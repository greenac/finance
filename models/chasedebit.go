package models

import (
	"errors"
	"github.com/greenac/finance/logger"
	"strconv"
	"strings"
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

func (cd *ChaseDebit) SetValues(ent []string) error {
	entries := cleanEntry(ent, numCDEntries)
	if len(entries) != numCDEntries {
		logger.Error("`ChaseDebit::SetValues` Invalid number of entries:", len(entries), "should be:", numCDEntries, entries)
		for i, e := range entries {
			logger.Warn(i, e)
		}
		return errors.New("INVALID_NUM_ENTRIES")
	}

	pd, err := time.Parse(DateLayout, entries[chaseDebitPostDateIndex])
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing date:", err)
		return errors.New("INVALID_ENTRY")
	}

	amt, err := handleParseFloat(entries[chaseDebitAmountIndex])
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing amount:", err)
		return errors.New("INVALID_ENTRY")
	}

	bal, err := handleParseFloat(entries[chaseDebitBalanceIndex])
	if err != nil {
		logger.Error("`ChaseDebit::SetValues` parsing balance:", err)
		return errors.New("INVALID_ENTRY")
	}

	cn, err := strconv.ParseInt(entries[chaseDebitCheckNumberIndex], 10, 64)
	if err != nil {
		if entries[chaseDebitCheckNumberIndex] == "" {
			cn = -1
		} else {
			logger.Error("`ChaseDebit::SetValues` parsing check number entry:", entries[chaseDebitCheckNumberIndex], err)
			return errors.New("INVALID_ENTRY")
		}
	}

	cd.Details = entries[chaseDebitDetailsIndex]
	cd.PostDate = pd
	cd.Description = entries[chaseDebitDescriptionIndex]
	cd.Amount = amt
	cd.Type = entries[chaseDebitTypeIndex]
	cd.Balance = bal
	cd.CheckNumber = cn

	return nil
}

func (cd *ChaseDebit) DebitedAmount() float64 {
	return -cd.Amount
}

func (cd *ChaseDebit) Desc() string {
	return cd.Description
}

func (cd *ChaseDebit) TransType() TransType {
	var t TransType
	if strings.ToLower(cd.Details) == "deposit" {
		t = Deposit
	} else {
		t = Withdrawal
	}

	return t
}

func (cd *ChaseDebit) TransDate() time.Time {
	return cd.PostDate
}
