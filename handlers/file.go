package handlers

import (
	"github.com/greenac/finance/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func PathsInDir(dp string) ([]string, error) {
	var paths []string
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("`PathsInDir` failed to get working dir:", err)
		return paths, err
	}

	fs, err := ioutil.ReadDir(filepath.Join(wd, dp))
	if err != nil {
		logger.Error("`PathsInDir` reading dir path:", dp, err)
		return paths, err
	}

	paths = make([]string, 0)
	for _, fp := range fs {
		if strings.Contains(strings.ToLower(fp.Name()), ".csv") {
			paths = append(paths, filepath.Join(dp, fp.Name()))
		}
	}

	return paths, nil
}
