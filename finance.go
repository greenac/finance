package main

import (
	"fmt"
	"github.com/greenac/finance/analysis"
	"github.com/greenac/finance/csv"
	"github.com/greenac/finance/json"
	"github.com/greenac/finance/logger"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("`main` loading .env file")
	}

	csvPath := os.Getenv("CSV_PATH")
	groupPath := os.Getenv("GROUP_PATH")

	if csvPath == "" || groupPath == "" {
		logger.Error("One or more env vars has not been set")
		panic("ENV_VARS_NOT_SET")
	}

	dats, err := csv.Parse(csvPath, true)
	if err != nil {
		logger.Error("`main` failed to read csv file at path:", csvPath)
	}

	group, err := json.Read(groupPath)
	if err != nil {
		logger.Error("`main` failed to read json file at path:", groupPath)
	}

	a := analysis.Analyzer{Datums: dats, Groups: group}
	a.GroupByType()

	s := a.Sum()
	out := "\n"
	for k, v := range *s {
		out += fmt.Sprint(k, ": ", v, "\n")
	}

	logger.Log(out)

	rnd := a.Bin(analysis.Random)
	if rnd == nil {
		logger.Error("No bin with key:", analysis.Random)
	} else {
		for i, d := range *rnd {
			logger.Log(i, d)
		}
	}
}
