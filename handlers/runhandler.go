package handlers

import (
	"errors"
	"github.com/greenac/finance/analysis"
	"github.com/greenac/finance/logger"
	"github.com/greenac/finance/models"
)

type RunHandler struct {
	analyzer     *analysis.Analyzer
	DirPaths     models.CSVDirPaths
	modelsByType models.ModelsByType
}

func (rh *RunHandler) Setup() {
	if len(rh.DirPaths) == 0 {
		logger.Error("`RunHandler::Setup` `Paths` not set")
		panic(errors.New("UNSET_REQUIRED_VAR"))
	}

	if rh.analyzer == nil {
		fps := make(models.CSVFilePaths)
		for mn, dp := range rh.DirPaths {
			fns, err := PathsInDir(dp)
			if err != nil {
				logger.Warn("`RunHandler::Setup` reading paths in directory:", dp, err)
				continue
			}

			fps[mn] = fns
		}

		mbt := make(models.ModelsByType)
		for mn, paths := range fps {
			parser := Parser{FilePaths: paths}
			mods, err := parser.Parse(mn, true)
			if err != nil {
				logger.Warn("`RunHandler::Setup` creating models:", mn, err)
				continue
			}

			mbt[mn] = mods
		}

		rh.modelsByType = mbt
	}
}

func (rh *RunHandler) Fill() {

}

func (rh *RunHandler) AllModels() *[]models.CsvModel {
	ms := make([]models.CsvModel, 0)
	for _, mods := range rh.modelsByType {
		for _, m := range *mods {
			ms = append(ms, m)
		}
	}

	return &ms
}
