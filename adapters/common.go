package adapters

import "github.com/owm-inc/weatherchecker-go/common"

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

func NewMeasurementArray() MeasurementArray {
	return make(MeasurementArray, 0)
}

// AdaptStub is an adapter for MeasurementArray constructor
func AdaptStub(s string) MeasurementArray { return NewMeasurementArray() }

type AdapterCollection struct {
	data    map[string](map[string]func(string) (MeasurementArray, error))
	dataSem chan struct{}
}

func (c *AdapterCollection) Add(source string, wt string, fn func(string) (MeasurementArray, error)) {
	<-c.dataSem

	if _, ok := c.data[source]; ok == false {
		c.data[source] = make(map[string]func(string) (MeasurementArray, error))
	}
	c.data[source][wt] = fn

	c.dataSem <- struct{}{}
}

func (c *AdapterCollection) Retrieve(source, wt string) (adaptFunc func(string) (MeasurementArray, error)) {
	<-c.dataSem

	sourceFuncs, aExists := c.data[source]
	if aExists == true {
		storedFunc, bExists := sourceFuncs[wt]
		if bExists == true {
			adaptFunc = storedFunc
		}
	}

	c.dataSem <- struct{}{}

	return adaptFunc
}

func MakeAdapterCollection() AdapterCollection {
	c := AdapterCollection{}
	c.data = make(map[string](map[string]func(string) (MeasurementArray, error)))
	c.dataSem = make(chan struct{}, 1)
	c.dataSem <- struct{}{}

	return c
}

func GetAdaptFunc(sourceName string, wtypeName string) (adaptFunc func(string) (MeasurementArray, error), err error) {
	fnColl := MakeAdapterCollection()
	fnColl.Add("owm", "current", OwmAdaptCurrentWeather)
	fnColl.Add("owm", "forecast", OwmAdaptForecast)
	fnColl.Add("wunderground", "current", WundergroundAdaptCurrentWeather)
	fnColl.Add("myweather2", "current", Myweather2AdaptCurrentWeather)
	fnColl.Add("forecast.io", "current", ForecastioAdaptCurrentWeather)
	fnColl.Add("forecast.io", "forecast", ForecastioAdaptForecast)
	fnColl.Add("worldweatheronline", "current", WorldweatheronlineAdaptCurrentWeather)
	fnColl.Add("worldweatheronline", "forecast", WorldweatheronlineAdaptForecast)

	adaptFunc = fnColl.Retrieve(sourceName, wtypeName)

	if adaptFunc == nil {
		err = common.NoAdaptFunc
	}

	return adaptFunc, err
}
