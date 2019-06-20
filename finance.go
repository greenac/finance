package main

import (
	"errors"
	"github.com/greenac/finance/env"
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

	capCreditPath := os.Getenv(env.CapOneCreditPath)
	chaseCreditPath := os.Getenv(env.ChaseCreditPath)
	chaseDebitPath := os.Getenv(env.ChaseDebitPath)
	groupPath := os.Getenv(env.GroupPath)

	if capCreditPath == "" || chaseCreditPath == "" || chaseDebitPath == "" || groupPath == "" {
		logger.Error("`main` One or more env vars has not been set")
		panic(errors.New("ENV_VARS_NOT_SET"))
	}

	dps := models.CSVDirPaths{
		models.ChaseCreditName:  chaseCreditPath,
		models.ChaseDebitName:   chaseDebitPath,
		models.CapOneCreditName: capCreditPath,
	}

	rh := handlers.RunHandler{DirPaths: dps}
	rh.Fill()
	err = rh.AddGroups(groupPath)
	if err != nil {
		logger.Error("`main` Adding groups:", err)
		panic(err)
	}

	err = rh.Analyze()
	if err != nil {
		logger.Error("`main` running analysis:", err)
		panic(err)
	}
}
