package models

import (
	"errors"
	"github.com/greenac/finance/logger"
	"strings"
	"time"
)

const (
	chaseCreditDateIndex        = 0
	chaseCreditPostDateIndex    = 1
	chaseCreditDescriptionIndex = 2
	chaseCreditCategoryIndex    = 3
	chaseCreditTypeIndex        = 4
	chaseCreditAmountIndex      = 5
)

const NumChaseCreditEntries = 6

type ChaseCredit struct {
	Date        time.Time
	PostDate    time.Time
	Description string
	Category    string
	Type        string
	Amount      float64
}

func (cc *ChaseCredit) SetValues(ents []string) error {
	entries := cleanEntry(ents, NumChaseCreditEntries)
	if len(entries) != NumChaseCreditEntries {
		logger.Error("`ChaseCredit::SetValues` Invalid number of entries:", len(entries), "should be:", NumChaseCreditEntries)
		return errors.New("INVALID_NUM_ENTRIES")
	}

	d, err := time.Parse(DateLayout, entries[chaseCreditDateIndex])
	if err != nil {
		logger.Error("`ChaseCredit::SetValues` parsing date:", err)
		return errors.New("INVALID_ENTRY")
	}

	pDate, err := time.Parse(DateLayout, entries[chaseCreditPostDateIndex])
	if err != nil {
		logger.Error("`ChaseCredit::SetValues` parsing post date:", err)
		return errors.New("INVALID_ENTRY")
	}

	amt, err := handleParseFloat(entries[chaseCreditAmountIndex])
	if err != nil {
		logger.Error("`ChaseCredit::SetValues` parsing amount:", err)
		return errors.New("INVALID_ENTRY")
	}

	cc.Date = d
	cc.PostDate = pDate
	cc.Description = entries[chaseCreditDescriptionIndex]
	cc.Category = entries[chaseCreditCategoryIndex]
	cc.Type = entries[chaseCreditTypeIndex]
	cc.Amount = amt

	return nil
}

func (cc *ChaseCredit) DebitedAmount() float64 {
	return -cc.Amount
}

func (cc *ChaseCredit) Desc() string {
	return cc.Description
}

func (cc *ChaseCredit) TransType() TransType {
	var t TransType
	tr := strings.ToLower(cc.Type)
	if tr == "payment" || tr == "return" {
		t = Withdrawal
	} else {
		t = Deposit
	}

	return t
}

func (cc *ChaseCredit) TransDate() time.Time {
	return cc.PostDate
}
