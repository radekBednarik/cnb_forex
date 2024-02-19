package parser

import (
	"strconv"
	"strings"
)

type ForexDataForDate struct {
	Date      string
	ForexData []singleCurrData
}

type singleCurrData struct {
	Country string
	Name    string
	Symbol  string
	Value   float64
}

func (d *ForexDataForDate) ParseFromText(inputText string) *ForexDataForDate {
	rows := strings.Split(inputText, "\n")

	d.Date = strings.Split(rows[0], " ")[0]

	for _, row := range rows[2 : len(rows)-1] {
		items := strings.Split(row, "|")
		fVal, err := strconv.ParseFloat(strings.Replace(items[4], ",", ".", 1), 64)
		if err != nil {
			fVal = 0.0
		}

		data := singleCurrData{Country: items[0], Name: items[1], Symbol: items[3], Value: fVal}

		d.ForexData = append(d.ForexData, data)

	}

	return d
}
