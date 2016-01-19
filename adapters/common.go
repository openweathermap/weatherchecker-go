package adapters

import (
	"errors"
)

type Measurement struct {
	Humidity      float64 `json:"humidity"`
	Pressure      float64 `json:"pressure"`
	Precipitation float64 `json:"precipitation"`
	Temp          float64 `json:"temp"`
	Wind          float64 `json:"wind"`
}

type MeasurementSchema struct {
	Data      Measurement `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type MeasurementArray []MeasurementSchema

var AdapterPanicErr = errors.New("Adapter panicking")
var nodeErr = errors.New(`Node not found`)

func AdaptStub(s string) MeasurementArray { return make(MeasurementArray, 0) }
func AdaptNull(s string) (measurements MeasurementArray, err error) {
	return AdaptStub(s), errors.New("No adapt function")
}

func AdaptWeather(sourceName string, wtypeName string, data string) (measurements MeasurementArray, err error) {
	var adaptFunc func(string) (MeasurementArray, error)
	var fnTable = make(map[string](map[string]func(string) (MeasurementArray, error)))

	for _, provider := range []string{"owm", "wunderground", "myweather2", "forecast.io", "worldweatheronline", "yandex", "accuweather", "gismeteo"} {
		fnTable[provider] = make(map[string]func(string) (MeasurementArray, error))
	}

	fnTable["owm"]["current"] = OwmAdaptCurrentWeather
	fnTable["owm"]["forecast"] = OwmAdaptForecast
	fnTable["wunderground"]["current"] = WundergroundAdaptCurrentWeather
	fnTable["myweather2"]["current"] = Myweather2AdaptCurrentWeather
	fnTable["forecast.io"]["current"] = ForecastioAdaptCurrentWeather
	fnTable["forecast.io"]["forecast"] = ForecastioAdaptForecast
	fnTable["worldweatheronline"]["current"] = WorldweatheronlineAdaptCurrentWeather
	fnTable["worldweatheronline"]["forecast"] = WorldweatheronlineAdaptForecast
	fnTable["yandex"]["current"] = YandexAdaptCurrentWeather
	fnTable["accuweather"]["current"] = AccuweatherAdaptCurrentWeather
	fnTable["gismeteo"]["current"] = GismeteoAdaptCurrentWeather

	adaptFunc = AdaptNull

	_, p_ok := fnTable[sourceName]
	if p_ok == true {
		storedFunc, f_ok := fnTable[sourceName][wtypeName]
		if f_ok == true {
			adaptFunc = storedFunc
		}
	}

	measurements, err = adaptFunc(data)

	return measurements, err
}
