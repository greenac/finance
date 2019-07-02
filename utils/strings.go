package utils

import (
	"strings"
)

func CombineQuotes(pts *[]string) *[]string {
	cpts := make([]string, 0)
	chk := false
	cur := ""
	for _, p := range *pts {
		if p == "" {
			cpts = append(cpts, p)
			continue
		}

		runes := []rune(p)
		if chk {
			// The end of the quoted comma entry
			if runes[len(runes)-1] == '"' {
				cur += strings.ReplaceAll(p, "\"", "")
				cpts = append(cpts, cur)
				cur = ""
				chk = false
			} else {
				cpts = append(cpts, cur)
			}
		} else if strings.Contains(p, "\"") {
			chk = true
			cur += strings.ReplaceAll(p, "\"", "")
		} else {
			cpts = append(cpts, p)
		}
	}

	return &cpts
}
