package handlers

import (
	"errors"
	"fmt"
	"github.com/greenac/finance/analysis"
	"github.com/greenac/finance/json"
	"github.com/greenac/finance/logger"
	"github.com/greenac/finance/models"
)

type RunHandler struct {
	analyzer     *analysis.Analyzer
	DirPaths     models.CSVDirPaths
	modelsByType models.ModelsByType
	groups       *json.Group
}

func (rh *RunHandler) Fill() {
	if len(rh.DirPaths) == 0 {
		logger.Error("`RunHandler::Setup` `Paths` not set")
		panic(errors.New("UNSET_REQUIRED_VAR"))
	}

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

func (rh *RunHandler) AddGroups(gp string) error {
	grp, err := json.Read(gp)
	if err != nil {
		logger.Error("`RunHandler::SetupGroups` failed to read from path:", gp)
		return err
	}

	rh.groups = grp
	return nil
}

func (rh *RunHandler) Analyze() error {
	if rh.groups == nil {
		logger.Error("`RunHandler::Analyze` no groups set")
		return errors.New("UNSET_VAR")
	}

	if rh.analyzer == nil {
		a := analysis.Analyzer{Models: &rh.modelsByType, Groups: rh.groups}
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
			logger.Log("Number of unknown entries:", len(*rnd))
			for i, m := range *rnd {
				logger.Log(i, m.Desc(), m.DebitedAmount())
			}
		}
	}

	return nil
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

func (rh *RunHandler) BinDeltaT() error {
	var a analysis.Analyzer
	if rh.analyzer == nil {
		a = analysis.Analyzer{Models: &rh.modelsByType, Groups: rh.groups}
		rh.analyzer = &a
	}

	rh.analyzer.GroupByDate()
	bwas := rh.analyzer.BinsWithAmounts()

	err := DrawSingleAxisBinWithAmounts(bwas)
	if err != nil {
		logger.Error("`RunHandler::BinDeltaT` drawing chart:", err)
		return err
	}

	return nil
}
