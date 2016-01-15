package adapters

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func AccuweatherIsMetric(htmlString string) (isMetric bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			isMetric = false
			err = AdapterPanicErr
		}
	}()
	htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if htmlErr != nil {
		return false, htmlErr
	}
	settingsString := htmlDoc.Find(`div#control-panel`).Find(`a#bt-menu-settings`).Find(`span.menu-arrow`).Contents().Nodes[1].Data
	if strings.Contains(settingsString, "°F") {
		return false, nil
	} else if strings.Contains(settingsString, "°C") {
		return true, nil
	} else {
		return false, errors.New("HTML parser error")
	}
}

func AccuweatherAdaptCurrentHumidity(htmlString string) (humidity float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			humidity = float64(0)
			err = AdapterPanicErr
		}
	}()
	htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if htmlErr != nil {
		return float64(0), htmlErr
	}

	node1 := htmlDoc.Find(`tr`).Nodes[5]
	node2 := goquery.NewDocumentFromNode(node1).Contents().Nodes[3]

	finalData := node2.FirstChild.Data

	humidityRaw := finalData
	humidityString := strings.TrimRight(strings.TrimSpace(humidityRaw), "%")

	humidity, convErr := strconv.ParseFloat(humidityString, 64)
	if convErr != nil {
		return float64(0), convErr
	}

	return humidity, nil
}

func AccuweatherAdaptCurrentTemp(htmlString string) (float64, error) {
	htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if htmlErr != nil {
		return float64(0), htmlErr
	}

	nodeGroup1 := htmlDoc.Find(`tr.temp`).Find(`td.first-col`).Contents().Nodes
	nodeIndex1 := 0
	if len(nodeGroup1)-1 < nodeIndex1 {
		return float64(0), nodeErr
	}
	node1 := nodeGroup1[nodeIndex1]

	finalData := node1.Data

	tempRaw := finalData
	tempString := strings.TrimRight(strings.TrimSpace(tempRaw), "°")
	temp, convErr := strconv.ParseFloat(tempString, 64)

	if convErr != nil {
		return float64(0), convErr
	}

	return temp, nil
}

func AccuweatherAdaptCurrentWeather(htmlString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(htmlString)
			err = AdapterPanicErr
		}
	}()

	dt := time.Now().Unix()

	isMetric, _ := AccuweatherIsMetric(htmlString)

	humidity, _ := AccuweatherAdaptCurrentHumidity(htmlString)
	precipitation := float64(0)
	pressure := float64(0)
	tempRaw, _ := AccuweatherAdaptCurrentTemp(htmlString)
	wind := float64(0)

	var temp float64
	if isMetric == true {
		temp = tempRaw
	} else {
		temp, _ = convertUnits(tempRaw, "F")
	}

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}
