package adapters

import (
	"encoding/json"
	"strconv"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/owm-inc/weatherchecker-go/common"
)

type GismeteoMinMaxAvg struct {
	Avg float64 `json:"avg"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type GismeteoWindData struct {
	Speed float64 `json:"speed"`
	Trend float64 `json:"trend"`
}

type GismeteoWeatherDataCommon struct {
	AirTemp    float64          `json:"air"`
	WaterTemp  float64          `json:"water"`
	TimeString string           `json:"time"`
	Pressure   float64          `json:"pressure"`
	Wind       GismeteoWindData `json:"wind"`
}

type GismeteoWeatherData struct {
	GismeteoWeatherDataCommon
	Humidity float64 `json:"humidity"`
}

type GismeteoWeatherDataAggregate struct {
	GismeteoWeatherDataCommon
	Humidity GismeteoMinMaxAvg `json:"humidity"`
}

type GismeteoMoonData struct {
	Luminosity float64 `json:"luminosity"`
}

type GismeteoDayData struct {
	Alias    string                       `json:"alias"`
	Date     string                       `json:"date"`
	FA       []GismeteoWeatherData        `json:"fa"`
	FO       []GismeteoWeatherData        `json:"fo"`
	MI       GismeteoWeatherDataAggregate `json:"mi"`
	MA       GismeteoWeatherDataAggregate `json:"ma"`
	MoonData GismeteoMoonData             `json:"moon"`
}

type GismeteoCityData struct {
	UTCOffset string            `json:"gmt"`
	CityName  string            `json:"name"`
	DayData   []GismeteoDayData `json:"days"`
}

type GismeteoInformerPartnerInfo struct {
	CityIds []int64 `json:"order"`
}

type GismeteoInformerData struct {
	WeatherData map[string]GismeteoCityData `json:"weather"`
	PartnerInfo GismeteoInformerPartnerInfo `json:"partner"`
}

type GismeteoInformerJson struct {
	Data []interface{} `json:"data"`
}

func gismeteoDecode(s string) (GismeteoInformerData, error) {
	var dataA GismeteoInformerJson
	var dataB GismeteoInformerData
	var errA error
	var errB error

	var byteString = []byte(strings.Trim(strings.Replace(strings.Replace(s, "parseResponse(", "", -1), ");", "", -1), "\n"))

	errA = json.Unmarshal(byteString, &dataA)

	if errA == nil {
		weatherDataMap := dataA.Data[1]

		temp1, _ := json.Marshal(&weatherDataMap)

		errB = json.Unmarshal(temp1, &dataB)

	}

	return dataB, errB
}

func GismeteoParseDateString(timeString, utcOffsetString string) (date time.Time, dateErr error) {
	timeOffset, timeOffsetErr := strconv.ParseInt(utcOffsetString, 10, 64)

	if timeOffsetErr != nil {
		return date, common.InvalidTimeOffsetString
	}

	secondsOffset := timeOffset * 3600

	timeArray := strings.Split(timeString, " ")
	if len(timeArray) != 6 {
		return date, common.InvalidTimeString
	}
	var timeParseErr error
	var y, mo, d, h, mi, s int64
	y, timeParseErr = strconv.ParseInt(timeArray[0], 10, 64)
	mo, timeParseErr = strconv.ParseInt(timeArray[1], 10, 64)
	d, timeParseErr = strconv.ParseInt(timeArray[2], 10, 64)
	h, timeParseErr = strconv.ParseInt(timeArray[3], 10, 64)
	mi, timeParseErr = strconv.ParseInt(timeArray[4], 10, 64)
	s, timeParseErr = strconv.ParseInt(timeArray[5], 10, 64)

	if timeParseErr != nil {
		return date, timeParseErr
	}

	date = time.Date(int(y), time.Month(mo), int(d), int(h), int(mi), int(s), 0, time.FixedZone("", int(secondsOffset)))
	return date, dateErr
}

func GismeteoAdaptCurrentWeather(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = common.AdapterPanicErr
		}
	}()

	data, decodeErr := gismeteoDecode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	cityId := data.PartnerInfo.CityIds[0]
	cityData := data.WeatherData[strconv.FormatInt(cityId, 10)]
	currentData := cityData.DayData[0].FA[0]

	timeOffsetRaw := cityData.UTCOffset
	date, dateErr := GismeteoParseDateString(currentData.TimeString, timeOffsetRaw)

	if dateErr != nil {
		return nil, dateErr
	}

	tempC := currentData.AirTemp
	pressure, _ := convertUnits(currentData.Pressure, "mmHg")
	windSpeed := currentData.Wind.Speed
	humidity := currentData.Humidity

	measurements = append(measurements, MeasurementSchema{Timestamp: date.Unix(), Data: Measurement{Temp: tempC, Pressure: pressure, Wind: windSpeed, Humidity: humidity}})

	return measurements, err
}

func GismeteoAdaptCurrentWeatherHtml(htmlString string) (measurements MeasurementArray, err error) {
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
