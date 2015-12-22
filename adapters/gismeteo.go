package adapters

import (
	"errors"
	"strconv"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"
)

func normalize_pressure(pressure float64, unit string) (float64, error) {
	unitNotFoundError := errors.New("Unit not found")

	unitTable := make(map[string]float64)
	unitTable["mmHg"] = 1013.25 / 760

	rate, u_ok := unitTable[unit]
	if u_ok == false {
		return float64(0), unitNotFoundError
	}

	result := pressure * rate

	return result, nil
}

func GismeteoAdaptCurrentWeather(htmlString string) (measurements MeasurementArray) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(htmlString)
		}
	}()

	dt := time.Now().Unix()

	htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if htmlErr != nil {
		return MeasurementArray{}
	}

	labelLookupTable := make(map[string]string)
	labelLookupTable["Ветер"] = "wind_speed"
	labelLookupTable["Давление"] = "pressure"
	labelLookupTable["Влажность"] = "humidity"

	infoItems := htmlDoc.Find(`div.info_item.clearfix`).Nodes

	var temp_raw string

	nodeGroup1 := htmlDoc.Find(`div.ii._temp`).Find(`span.value.js_value.val_to_convert`).Contents().Nodes
	nodeIndex1 := 0
	if len(nodeGroup1) < nodeIndex1+1 {
	} else {
		node1 := nodeGroup1[nodeIndex1]
		temp_raw = node1.Data
	}

	itemMap := make(map[string]string)

	for _, item := range infoItems {
		nodeGroup1 := goquery.NewDocumentFromNode(item).Find(`div.ii.info_label`).Contents().Nodes
		nodeIndex1 := 0
		if len(nodeGroup1) < nodeIndex1+1 {
			continue
		}
		node1 := nodeGroup1[nodeIndex1]

		lookupKey := string(node1.Data)
		mapKey, h_ok := labelLookupTable[lookupKey]
		if h_ok == false {
			continue
		}

		nodeGroup2 := goquery.NewDocumentFromNode(item).Find(`div.ii.info_value`).Contents().Nodes
		nodeIndex2 := 0
		if len(nodeGroup2) < nodeIndex2+1 {
			continue
		}
		node2 := nodeGroup2[nodeIndex2]

		value := string(node2.FirstChild.Data)

		itemMap[mapKey] = value
	}

	entry := Measurement{}

	temp, err := strconv.ParseFloat(temp_raw, 64)
	if err == nil {
		entry.Temp = temp
	}

	humidity_raw, im_ok := itemMap["humidity"]
	if im_ok == true {
		humidity, err := strconv.ParseFloat(humidity_raw, 64)
		if err == nil {
			entry.Humidity = humidity
		}
	}
	wind_speed_raw, im_ok := itemMap["wind_speed"]
	if im_ok == true {
		wind, err := strconv.ParseFloat(wind_speed_raw, 64)
		if err == nil {
			entry.Wind = wind
		}
	}
	pressure_raw, im_ok := itemMap["pressure"]
	if im_ok == true {
		pressure, err := strconv.ParseFloat(pressure_raw, 64)
		if err == nil {
			entry.Pressure = pressure
		}
	}

	measurements = append(measurements, MeasurementSchema{Data: entry, Timestamp: dt})

	return measurements
}
