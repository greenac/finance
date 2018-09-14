package models

import (
	"errors"
	"github.com/greenac/finance/logger"
	"strconv"
	"time"
)

const (
	capOneCreditStageIndex     = 0
	capOneCreditDateIndex    = 1
	capOneCreditPostDateIndex    = 2
	capOneCreditCardNumIndex    = 3
	capOneCreditDescriptionIndex = 4
	capOneCreditCategoryIndex        = 5
	capOneCreditDebitIndex     = 6
	capOneCreditCreditIndex = 7
)

const numOfCapOneCreditEntries = 8

type CapOneCredit struct {
	Stage       string
	Date   time.Time
	PostDate    time.Time
	CardNum     int64
	Description string
	Category    string
	Debit       float64
	Credit      float64
}

func (cc *CapOneCredit) SetValues(entries []string) error {
	if len(entries) != numOfCapOneCreditEntries {
		logger.Error("`CapOneCredit::SetValues` Invalid number of entries:", len(entries), "should be:", numOfCapOneCreditEntries)
		return errors.New("INVALID_NUM_ENTRIES")
	}

	d, err := time.Parse(DateLayout, entries[capOneCreditDateIndex])
	if err != nil {
		logger.Error("`CapOneCredit::SetValues` parsing  date:", err)
		return errors.New("INVALID_ENTRY")
	}

	pd, err := time.Parse(DateLayout, entries[capOneCreditPostDateIndex])
	if err != nil {
		logger.Error("`CapOneCredit::SetValues` parsing post date:", err)
		return errors.New("INVALID_ENTRY")
	}

	cn, err := strconv.ParseInt(entries[capOneCreditCardNumIndex], 10, 64)
	if err != nil {
		logger.Error("`CapOneCredit::SetValues` parsing credit card number:", err)
		return errors.New("INVALID_ENTRY")
	}

	db, err := strconv.ParseFloat(entries[capOneCreditDebitIndex], 64)
	if err != nil {
		logger.Error("`CapOneCredit::SetValues` parsing debit:", err)
		return errors.New("INVALID_ENTRY")
	}

	cr, err := strconv.ParseFloat(entries[capOneCreditCreditIndex], 64)
	if err != nil {
		logger.Error("`CapOneCredit::SetValues` parsing credit:", err)
		return errors.New("INVALID_ENTRY")
	}


	cc.Stage = entries[capOneCreditStageIndex]
	cc.Date = d
	cc.PostDate = pd
	cc.CardNum = cn
	cc.Description = entries[capOneCreditDescriptionIndex]
	cc.Category = entries[capOneCreditCategoryIndex]
	cc.Debit = db
	cc.Credit= cr

	return nil
}

func (cc *CapOneCredit) DebitedAmount () float64 {
	return cc.Debit
}
