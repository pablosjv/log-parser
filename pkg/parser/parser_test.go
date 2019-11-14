package parser

import (
	"testing"
	"time"
)

func TestParseLine(test *testing.T) {
	line := "1565647246869 Jeremyah Morrigan"
	time := time.Unix(1565647246869, 0)
	origin := "Jeremyah"
	dest := "Morrigan"
	t, o, d := parseLine(line)

	if time != t || dest != d || origin != o {
		test.Error()
	}
}
