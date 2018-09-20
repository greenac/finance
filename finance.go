package main

import (
	"errors"
	"github.com/greenac/finance/handlers"
	"github.com/greenac/finance/logger"
	"github.com/greenac/finance/models"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("`main` loading .env file")
		panic(errors.New("NO_DOT_ENV"))
	}

	capCreditPath := os.Getenv("CAPONE_CREDIT_PATH")
	chaseCreditPath := os.Getenv("CHASE_CREDIT_PATH")
	chaseDebitPath := os.Getenv("CHASE_DEBIT_PATH")
	groupPath := os.Getenv("GROUP_PATH")

	if capCreditPath == "" || chaseCreditPath == "" || chaseDebitPath == "" || groupPath == "" {
		logger.Error("`main` One or more env vars has not been set")
		panic(errors.New("ENV_VARS_NOT_SET"))
	}

	//group, err := json.Read(groupPath)
	//if err != nil {
	//	logger.Error("`main` failed to read json file at path:", groupPath)
	//}
	//
	//a := analysis.Analyzer{Datums: dats, Groups: group}
	//a.GroupByType()
	//
	//s := a.Sum()
	//out := "\n"
	//for k, v := range *s {
	//	out += fmt.Sprint(k, ": ", v, "\n")
	//}
	//
	//logger.Log(out)
	//
	//rnd := a.Bin(analysis.Random)
	//if rnd == nil {
	//	logger.Error("No bin with key:", analysis.Random)
	//} else {
	//	for i, d := range *rnd {
	//		logger.Log(i, d)
	//	}
	//}

	dps := models.CSVDirPaths{
		models.ChaseCreditName:  chaseCreditPath,
		models.ChaseDebitName:   chaseDebitPath,
		models.CapOneCreditName: capCreditPath,
	}

	rh := handlers.RunHandler{DirPaths: dps}
	rh.Setup()
}
