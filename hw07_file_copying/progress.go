package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Progress struct {
	total   int64
	current int64
	percent int64
	output  io.Writer
	width   uint8
}

func (p *Progress) Add(chunk int64) {
	if p.total <= 0 || chunk <= 0 {
		return
	}

	if p.current == 0 {
		fmt.Fprint(p.output, p.getProgressBar(0))
	}

	p.current += chunk

	c := float64(p.current)
	t := float64(p.total)

	percent := int64(c / t * 100)

	if percent > 0 && percent <= 100 && percent != p.percent {
		p.percent = percent

		fmt.Fprint(p.output, "\r"+p.getProgressBar(percent))
	}
}

func (p *Progress) Finish() {
	if p.current > 0 {
		fmt.Fprintln(p.output)
	}
}

func (p *Progress) getProgressBar(percent int64) string {
	var bar string

	switch percent {
	case 0:
		bar = ">" + strings.Repeat(" ", int(p.width-1))
	case 100:
		bar = strings.Repeat("=", int(p.width))
	default:
		chunks := int(float64(percent) / (100 / float64(p.width)))

		bar = strings.Repeat("=", chunks) + ">" + strings.Repeat(" ", int(p.width)-chunks-1)
	}

	return fmt.Sprintf(" %3d%% [%s]", percent, bar)
}

func NewProgress(total int64) *Progress {
	return &Progress{
		total:  total,
		output: os.Stdout,
		width:  50,
	}
}
