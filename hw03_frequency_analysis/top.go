package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const tokenLimit = 10

var rgx = regexp.MustCompile(`(?:\pP*((?:\pP*\pL+)+)\pP*)|(\pP{2,})`)

type token struct {
	str   string
	count int
}

type tokenList []*token

func (tl tokenList) Len() int {
	return len(tl)
}

func (tl tokenList) Less(i, j int) bool {
	a := tl[i]
	b := tl[j]

	if a.count == b.count {
		return a.str < b.str
	}

	return a.count > b.count
}

func (tl tokenList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func Top10(s string) []string {
	tokens := make(tokenList, 0)
	pointers := make(map[string]*token)

	for _, group := range rgx.FindAllStringSubmatch(s, -1) {
		sub := group[2]

		if sub == "" {
			sub = strings.ToLower(group[1])
		}

		if _, ok := pointers[sub]; !ok {
			t := token{sub, 1}

			tokens = append(tokens, &t)
			pointers[sub] = &t

			continue
		}

		pointers[sub].count++
	}

	sort.Sort(tokens)

	sorted := make([]string, 0)

	for i, p := range tokens {
		if i == tokenLimit {
			break
		}

		sorted = append(sorted, p.str)
	}

	return sorted
}
