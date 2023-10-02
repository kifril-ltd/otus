package progressbar

import (
	"fmt"
	"math"
)

type ProgressBar struct {
	percent  int64
	cur      int64
	total    int64
	progress string
	symbol   string
	width    int64
}

type BarOption func(bar *ProgressBar)

func WithSymbol(symbol string) BarOption {
	return func(bar *ProgressBar) {
		bar.symbol = symbol
	}
}

func WithWidth(width int64) BarOption {
	return func(bar *ProgressBar) {
		bar.width = width
	}
}

func NewBar(total int64, options ...BarOption) *ProgressBar {
	bar := &ProgressBar{
		percent:  0,
		cur:      0,
		total:    total,
		progress: "",
		symbol:   "#",
		width:    50,
	}

	for _, opt := range options {
		opt(bar)
	}

	return bar
}

func (bar *ProgressBar) Draw(cur int64) {
	bar.cur = cur
	bar.percent = bar.calcPercent()

	for len(bar.progress) < bar.calcProgressLen() {
		bar.progress += bar.symbol
	}

	fmt.Printf("\r[%-*s] %3d%% %8d/%d", bar.width, bar.progress, bar.percent, bar.cur, bar.total)
}

func (bar *ProgressBar) calcPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *ProgressBar) calcProgressLen() int {
	return int(math.Floor(float64(bar.percent) * float64(bar.width) / 100))
}

func (bar *ProgressBar) Finish() {
	fmt.Println()
}
