package handlers

import (
	"errors"
	"github.com/greenac/finance/logger"
	"io"
	"os"
	"strings"
)

type LinesForPath struct {
	path  string
	Lines *[][]string
}

type FilePath struct {
	Entries int
	Path    string
}

type Parser struct {
	FilePaths []FilePath
}

func (p *Parser) read(path string) (*[][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		logger.Error("`Parser::read` failed to open name file at path:", path, err)
		return nil, err
	}

	const buffer = 1000
	var offset int64 = 0
	cont := true
	data := make([]byte, 0)
	for cont {
		d := make([]byte, buffer)
		n, err := f.ReadAt(d, offset)
		if err != nil {
			if err == io.EOF {
				cont = false
			} else {
				logger.Error("`Parser::read` Reading file at path:", path, err)
				return nil, err
			}
		}

		offset += int64(n)
		data = append(data, d[:n]...)
	}

	lines := make([][]byte, 0)
	i := 0
	for _, b := range data {
		if b == '\n' {
			i += 1
		} else if i == len(lines) {
			lines = append(lines, []byte{b})
		} else {
			lines[i] = append(lines[i], b)
		}
	}

	return &lines, nil
}

func (p *Parser) parseFileAtPath(path string) (*[][]byte, error) {
	lines, err := p.read(path)
	if err != nil {
		logger.Error("`Parser::parseFileAtPath` reading lines from:", path, err)
		return nil, err
	}

	return lines, nil
}

func (p *Parser) LinesForFile(path FilePath) (*[][]string, error) {
	lines := make([][]string, 0)
	ls, err := p.parseFileAtPath(path.Path)
	if err != nil {
		logger.Error("`Parser::LinesForFile` parsing file:", path)
		return nil, err
	}

	for i, l := range *ls {
		// Don't read the header of the csv file
		if i == 0 {
			continue
		}

		pts := strings.Split(string(l), ",")
		if len(pts) != path.Entries {
			logger.Error("`Parser::LinesForFile` has:", len(pts), "entries. ", string(l), "expected:", path.Entries)
			return nil, errors.New("INVALID_CSV_FORMAT")
		}

		lines = append(lines, pts)
	}

	return &lines, nil
}

func (p *Parser) Parse() (*[]LinesForPath, error) {
	lfp := make([]LinesForPath, 0)
	for _, path := range p.FilePaths {
		lines, err := p.LinesForFile(path)
		if err != nil {
			logger.Error("`Parser::Parse` reading lines from:", path, err)
			return nil, err
		}

		lfp = append(lfp, LinesForPath{path.Path, lines})
	}

	return &lfp, nil
}
