package adapters

import (
	"strconv"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/owm-inc/weatherchecker-go/common"
)

func GismeteoAdaptCurrentWeather(htmlString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(htmlString)
			err = common.AdapterPanicErr
		}
	}()

	dt := time.Now().Unix()

	htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if htmlErr != nil {
		return AdaptStub(htmlString), htmlErr
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
		temp_raw = strings.Replace(strings.Trim(node1.Data, `"`), "−", "-", -1)
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

		itemMap[mapKey] = strings.Trim(value, `"`)
	}

	entry := Measurement{}

	temp, tempConvErr := strconv.ParseInt(strings.Trim(temp_raw, `"`), 10, 64)
	if tempConvErr == nil {
		entry.Temp = float64(temp)
	} else {
		err = tempConvErr
		return measurements, err
	}

	humidity_raw, im_ok := itemMap["humidity"]
	if im_ok == true {
		humidity, humidityConvErr := strconv.ParseFloat(humidity_raw, 64)
		if humidityConvErr == nil {
			entry.Humidity = humidity
		} else {
			err = humidityConvErr
			return measurements, err
		}
	}
	wind_speed_raw, im_ok := itemMap["wind_speed"]
	if im_ok == true {
		windSpeed, windSpeedErr := strconv.ParseFloat(wind_speed_raw, 64)
		if windSpeedErr == nil {
			entry.Wind = windSpeed
		} else {
			err = windSpeedErr
			return measurements, err
		}
	}
	pressure_raw, im_ok := itemMap["pressure"]
	if im_ok == true {
		pressureUnconv, pressureUnconvErr := strconv.ParseFloat(pressure_raw, 64)
		if pressureUnconvErr == nil {
			pressure, pressureErr := convertUnits(pressureUnconv, "mmHg")
			if pressureErr == nil {
				entry.Pressure = pressure
			} else {
				err = pressureErr
				return measurements, err
			}
		} else {
			err = pressureUnconvErr
			return measurements, err
		}
	}

	measurements = append(measurements, MeasurementSchema{Data: entry, Timestamp: dt})

	return measurements, err
}
