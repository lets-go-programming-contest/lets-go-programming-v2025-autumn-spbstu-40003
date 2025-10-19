package sorter

import (
	"sort"

	"github.com/wedwincode/task-3/internal/reader"
)

func Sort(valutes []reader.Valute) []reader.Valute {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	return valutes
}
