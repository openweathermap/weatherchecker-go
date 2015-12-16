package adapters

import (
	"time"
)

type Measurement struct {
	Humidity      float64
	Pressure      float64
	Precipitation float64
	Temp          float64
	Wind          float64
}

type MeasurementSchema struct {
	Data      Measurement
	Timestamp time.Time
}

type MeasurementArray []MeasurementSchema

func AdaptStub(s string) MeasurementArray { return MeasurementArray{} }

func AdaptWeather(sourceName string, wtypeName string, data string) (measurements MeasurementArray) {
	var adaptFunc func(string) MeasurementArray
	var fnTable = make(map[string](map[string]func(string) MeasurementArray))

	for _, provider := range []string{"OpenWeatherMap", "Weather Underground", "MyWeather2", "Forecast.io", "WorldWeatherOnline", "AccuWeather", "Gismeteo"} {
		fnTable[provider] = make(map[string]func(string) MeasurementArray)
	}

	fnTable["OpenWeatherMap"]["current"] = OwmAdaptCurrentWeather
	fnTable["OpenWeatherMap"]["forecast"] = OwmAdaptForecast
	fnTable["Weather Underground"]["current"] = WundergroundAdaptCurrentWeather
	fnTable["MyWeather2"]["current"] = Myweather2AdaptCurrentWeather
	fnTable["Forecast.io"]["current"] = ForecastioAdaptCurrentWeather
	fnTable["WorldWeatherOnline"]["current"] = WorldweatheronlineAdaptCurrentWeather
	fnTable["AccuWeather"]["current"] = AccuweatherAdaptCurrentWeather
	fnTable["Gismeteo"]["current"] = GismeteoAdaptCurrentWeather

	adaptFunc = AdaptStub

	_, p_ok := fnTable[sourceName]
	if p_ok == true {
		storedFunc, f_ok := fnTable[sourceName][wtypeName]
		if f_ok == true {
			adaptFunc = storedFunc
		}
	}

	measurements = adaptFunc(data)

	return measurements
}
