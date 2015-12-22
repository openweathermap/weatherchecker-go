package adapters

import (
	"encoding/json"
)

type OwmLocationCoords struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type OwmWeatherInfo struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type OwmMainInfo struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}

type OwmWindInfo struct {
	Speed  float64 `json:"speed"`
	Degree float64 `json:"deg"`
}

type OwmPrecipInfo struct {
	Expect3h float64 `json:"3h"`
}

type OwmCloudInfo struct {
	All int `json:"all"`
}

type OwmWeatherStruct struct {
	Name       string           `json:"name"`
	Weather    []OwmWeatherInfo `json:"weather"`
	Main       OwmMainInfo      `json:"main"`
	Visibility int              `json:"visibility"`
	Wind       OwmWindInfo      `json:"wind"`
	Snow       OwmPrecipInfo    `json:"snow"`
	Rain       OwmPrecipInfo    `json:"rain"`
	Clouds     OwmCloudInfo     `json:"clouds"`
	Timestamp  int64            `json:"dt"`
	Id         int              `json:"id"`
}

type OwmCityInfo struct {
	OwmId      int               `json:"id"`
	Name       string            `json:"name"`
	Coord      OwmLocationCoords `json:"coord"`
	CountryISO string            `json:"country"`
	Population int               `json:"population"`
}

type OwmCurrentStruct struct {
	OwmWeatherStruct
	Coord OwmLocationCoords `json:"coord"`
	Code  int               `json:"cod"`
}

type OwmForecastStruct struct {
	Code    string             `json:"cod"`
	Message float64            `json:"message"`
	City    OwmCityInfo        `json:"city"`
	Count   int                `json:"cnt"`
	Data    []OwmWeatherStruct `json:"list"`
}

func owmCurrentDecode(s string) (data OwmCurrentStruct) {
	var byteString = []byte(s)

	json.Unmarshal(byteString, &data)

	return data
}

func owmForecastDecode(s string) (data OwmForecastStruct) {
	var byteString = []byte(s)

	json.Unmarshal(byteString, &data)

	return data
}

func OwmAdaptCurrentWeather(jsonString string) (measurements MeasurementArray) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
		}
	}()

	var data = owmCurrentDecode(jsonString)

	dt := int64(data.Timestamp)

	temp := float64(data.Main.Temp)
	pressure := float64(data.Main.Pressure)
	wind := float64(data.Wind.Speed)
	humidity := float64(data.Main.Humidity)
	precipitation := float64(0)

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements
}

func OwmAdaptForecast(jsonString string) (measurements MeasurementArray) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
		}
	}()
	var data = owmForecastDecode(jsonString)

	for _, entry := range data.Data {
		dt := int64(entry.Timestamp)

		temp := float64(entry.Main.Temp)
		pressure := float64(entry.Main.Pressure)
		wind := float64(entry.Wind.Speed)
		humidity := float64(entry.Main.Humidity)
		precipitation := float64(entry.Snow.Expect3h) + float64(entry.Rain.Expect3h)

		measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})
	}

	return measurements
}
