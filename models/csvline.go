package models

const DateLayout = "01/02/2006"

type CsvLine interface {
	SetValues([]string) error
	DebitedAmount() float64
}
