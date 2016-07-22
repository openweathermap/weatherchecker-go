package main

import (
	"errors"
	"sync"
)

type WeatherType int

const (
	WTUnknown WeatherType = iota
	WTCurrent
	WTForecast
)

func (t WeatherType) GetString() (string, error) {
	switch t {
	case WTCurrent:
		return "current", nil
	case WTForecast:
		return "forecast", nil
	default:
		return "", errors.New("WeatherType unknown")
	}
}

func StringToWT(s string) (WeatherType, error) {
	switch s {
	case "current":
		return WTCurrent, nil
	case "forecast":
		return WTForecast, nil
	default:
		return WTUnknown, errors.New("WeatherType unknown")
	}
}

type AdaptFunc func(string) (MeasurementArray, error)

// Measurement represents the data extracted from provider data.
type Measurement struct {
	Humidity      float64 `json:"humidity"`
	Pressure      float64 `json:"pressure"`
	Precipitation float64 `json:"precipitation"`
	Temp          float64 `json:"temp"`
	Wind          float64 `json:"wind"`
}

// MeasurementSchema is the holding structure for provider response.
type MeasurementSchema struct {
	Data      Measurement `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// MeasurementArray is a collection of provider responses
type MeasurementArray []MeasurementSchema

// AdaptStub is an adapter for MeasurementArray constructor
func AdaptStub(s string) MeasurementArray { return MeasurementArray{} }

type AdapterCollection struct {
	sync.Mutex
	data map[string](map[WeatherType]AdaptFunc)
}

func (c *AdapterCollection) exec(fn func()) {
	c.Lock()
	defer c.Unlock()

	fn()
}

func (c *AdapterCollection) Set(source string, wt WeatherType, fn AdaptFunc) {
	c.exec(func() {
		if _, ok := c.data[source]; ok == false {
			c.data[source] = map[WeatherType]AdaptFunc{}
		}
		c.data[source][wt] = fn
	})
}

func (c *AdapterCollection) Get(source string, wt WeatherType) (AdaptFunc, bool) {
	var adaptFunc AdaptFunc
	var exists = false
	c.exec(func() {
		sourceFuncs, aExists := c.data[source]
		if aExists == true {
			storedFunc, bExists := sourceFuncs[wt]
			if bExists == true {
				adaptFunc = storedFunc
				exists = true
			}
		}
	})

	return adaptFunc, exists
}

func MakeAdapterCollection() *AdapterCollection {
	return &AdapterCollection{data: map[string](map[WeatherType]AdaptFunc){}}
}

func GetAdaptFunc(sourceName string, wtypeName string) (AdaptFunc, error) {
	var wt, wtErr = StringToWT(wtypeName)
	if wtErr != nil {
		return nil, wtErr
	}

	fnColl := MakeAdapterCollection()
	fnColl.Set("owm", WTCurrent, OwmAdaptCurrentWeather)
	fnColl.Set("owm", WTForecast, OwmAdaptForecast)
	fnColl.Set("wunderground", WTCurrent, WundergroundAdaptCurrentWeather)
	fnColl.Set("myweather2", WTCurrent, Myweather2AdaptCurrentWeather)
	fnColl.Set("forecast.io", WTCurrent, ForecastioAdaptCurrentWeather)
	fnColl.Set("forecast.io", WTForecast, ForecastioAdaptForecast)
	fnColl.Set("worldweatheronline", WTCurrent, WorldweatheronlineAdaptCurrentWeather)
	fnColl.Set("worldweatheronline", WTForecast, WorldweatheronlineAdaptForecast)

	var adaptFunc, exists = fnColl.Get(sourceName, wt)

	if !exists {
		return nil, NoAdaptFunc
	}

	return adaptFunc, nil
}
