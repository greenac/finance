package json

import (
	"encoding/json"
	"errors"
	"github.com/greenac/finance/logger"
	"io/ioutil"
	"strings"
)

type Group map[string][]string

func Read(path string) (*Group, error) {
	if path == "" {
		logger.Error("`Reader::Fill` path is not set")
		return nil, errors.New("PATH_NOT_SET")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("`Reader::Fill` failed to read file at path:", path, err)
		return nil, err
	}

	group := make(Group, 0)
	err = json.Unmarshal(data, &group)
	if err != nil {
		logger.Error("`Reader::Fill` failed to unmarshal json file:", path, err)
		return nil, err
	}

	for _, names := range group {
		for i, n := range names {
			names[i] = strings.ToLower(n)
		}
	}

	return &group, nil
}
