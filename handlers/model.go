package handlers

import (
	"errors"
	"github.com/greenac/finance/logger"
	"github.com/greenac/finance/models"
)

func CreateModel(l *[]string, mn models.Name) (*models.CsvModel, error) {
	var m models.CsvModel
	switch mn {
	case models.ChaseCreditName:
		m = &models.ChaseCredit{}
	case models.ChaseDebitName:
		m = &models.ChaseDebit{}
	case models.CapOneCreditName:
		m = &models.CapOneCredit{}
	default:
		logger.Error("`CreateModel` unhandled model name:", mn)
		return &m, errors.New("UNHANDLED_SWITCH_CASE")
	}

	err := m.SetValues(*l)
	if err != nil {
		logger.Error("`CreateModel` setting value for model:", mn, "line:", *l, "with length:", len(*l))
		for i, ll := range *l {
			logger.Warn(i, ll)
		}

		return nil, err
	}

	return &m, nil
}
