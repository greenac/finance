package csv

import (
	"errors"
	"github.com/greenac/finance/logger"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)


type Parser struct {
	FilePaths []string
}

func (p *Parser) Read(path string) (*[][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		logger.Error("Failed to open name file at path:", path, err)
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
				logger.Error("Reading file at path:", path, err)
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

func (p *Parser) Parse(path string, low bool) (*[]Datum, error) {
	lines, err := read(path)
	if err != nil {
		logger.Error("Reading lines from:", path)
		return nil, err
	}

	dats := make([]Datum, 0)
	for i, l := range *lines {
		if i == 0 {
			continue
		}

		pts := make([]string, 0)
		pt := make([]byte, 0)
		for _, b := range l {
			if b == ',' {
				var txt string
				if low {
					txt = strings.ToLower(string(pt))
				} else {
					txt = string(pt)
				}

				pts = append(pts, txt)
				pt = []byte{}
			} else {
				pt = append(pt, b)
			}
		}

		ln, err := createLine(&pts)
		if err != nil {
			continue
		}

		dats = append(dats, *ln)
	}

	return &dats, nil
}

func (p *Parser) CreateModel(l *[]string) (*Datum, error) {
	ll := *l
	if len(ll) != ENTRIES_IN_LINE {
		return nil, errors.New("MALFORMED_LINE")
	}

	tType := ll[TypeIndex]
	des := ll[descriptionIndex]

	tDate, err := time.Parse(DATE_LAYOUT, ll[DateIndex])
	if err != nil {
		logger.Error("`createLine` parsing transaction date:", err)
		return nil, err
	}

	pDate, err := time.Parse(DATE_LAYOUT, ll[postDateIndex])
	if err != nil {
		logger.Error("`createLine` parsing post date:", err)
		return nil, err
	}

	amt, err := strconv.ParseFloat(ll[amountIndex], 64)
	if err != nil {
		logger.Error("`createLine` parsing amount:", err)
		return nil, err
	}

	ln := Datum{
		Type:        tType,
		Date:        tDate,
		PostDate:    pDate,
		Description: des,
		Amount:      -amt,
	}

	return &ln, nil
}
