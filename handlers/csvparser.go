package handlers

import (
	"github.com/greenac/finance/logger"
	"github.com/greenac/finance/models"
	"io"
	"os"
	"strings"
)

type Parser struct {
	FilePaths []string
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
		} else {
			if i == len(lines) {
				lines = append(lines, []byte{b})
			} else {
				lines[i] = append(lines[i], b)
			}
		}
	}

	return &lines, nil
}

func (p *Parser) Parse(n models.Name, useLowerCase bool) (*[]models.CsvModel, error) {
	mods := make([]models.CsvModel, 0)
	for _, path := range p.FilePaths {
		lines, err := p.read(path)
		if err != nil {
			logger.Error("`Parser::Parse` reading lines from:", path, err)
			return nil, err
		}

		for i, l := range *lines {
			if i == 0 {
				continue
			}

			inQuotes := false
			pts := make([]string, 0)
			pt := make([]byte, 0)
			for j, b := range l {
				if b == '"' {
					inQuotes = !inQuotes
				}

				if b == ',' && !inQuotes {
					var txt string
					if useLowerCase {
						txt = strings.ToLower(string(pt))
					} else {
						txt = string(pt)
					}

					pts = append(pts, txt)

					if j == len(l) - 1 || l[j + 1] == ',' {
						pts = append(pts, "")
					}

					pt = []byte{}
				} else {
					pt = append(pt, b)
				}
			}

			ln, err := CreateModel(&pts, n)
			if err != nil {
				continue
			}

			mods = append(mods, *ln)
		}
	}


	return &mods, nil
}

func (p *Parser) Clean() {

}
