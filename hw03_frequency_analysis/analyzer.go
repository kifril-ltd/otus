package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type FrequencyAnalyser struct {
	NormalizeFuncs []NormalizeFunc

	frequencyMap  FrequencyMap
	frequencyList FrequencyList
}

func NewFrequencyAnalyser() *FrequencyAnalyser {
	return &FrequencyAnalyser{}
}

func (fa *FrequencyAnalyser) AddNormalizeFunc(fn ...NormalizeFunc) {
	fa.NormalizeFuncs = append(fa.NormalizeFuncs, fn...)
}

func (fa *FrequencyAnalyser) Analyse(in string) {
	fa.buildMap(in)
	fa.buildList()
}

func (fa *FrequencyAnalyser) Top(limit int) []string {
	fa.sortWords()

	l := min(limit, len(fa.frequencyList))
	top := make([]string, l)

	for i, item := range fa.frequencyList[:l] {
		top[i] = item.word
	}

	return top
}

func (fa *FrequencyAnalyser) buildMap(in string) {
	res := make(FrequencyMap)

	words := strings.Fields(in)
	for _, word := range words {
		normalized, err := fa.normalize(word)
		if normalized == "" || err != nil {
			continue
		}
		res[normalized]++
	}

	fa.frequencyMap = res
}

func (fa *FrequencyAnalyser) buildList() {
	res := make(FrequencyList, len(fa.frequencyMap))
	i := 0

	for word, freq := range fa.frequencyMap {
		res[i] = FrequencyListItem{word: word, frequency: freq}
		i++
	}

	fa.frequencyList = res
}

func (fa *FrequencyAnalyser) normalize(word string) (string, error) {
	res := word

	for _, norm := range fa.NormalizeFuncs {
		var err error
		res, err = norm(res)
		if err != nil {
			return "", err
		}
	}

	return res, nil
}

func (fa *FrequencyAnalyser) sortWords() {
	if !sort.IsSorted(sort.Reverse(fa.frequencyList)) {
		sort.Stable(sort.Reverse(&fa.frequencyList))
	}
}

type (
	NormalizeFunc func(string) (string, error)
	FrequencyMap  map[string]int
)

type FrequencyListItem struct {
	word      string
	frequency int
}

type FrequencyList []FrequencyListItem

func (fl FrequencyList) Len() int {
	return len(fl)
}

func (fl FrequencyList) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}

func (fl FrequencyList) Less(i, j int) bool {
	if fl[i].frequency != fl[j].frequency {
		return fl[i].frequency < fl[j].frequency
	}
	return fl[i].word > fl[j].word
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
