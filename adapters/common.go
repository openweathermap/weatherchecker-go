package adapters

type Measurement struct {
	Humidity      float64
	Pressure      float64
	Precipitation float64
	Temp          float64
	Wind          float64
}

type MeasurementSchema struct {
	Data      Measurement
	Timestamp int64
}

type MeasurementArray []MeasurementSchema

func AdaptStub(s string) MeasurementArray { return MeasurementArray{} }

func AdaptWeather(sourceName string, wtypeName string, data string) (measurements MeasurementArray) {
	var adaptFunc func(string) MeasurementArray
	var fnTable = make(map[string](map[string]func(string) MeasurementArray))

	for _, provider := range []string{"owm", "wunderground", "myweather2", "forecast.io", "worldweatheronline", "accuweather", "gismeteo"} {
		fnTable[provider] = make(map[string]func(string) MeasurementArray)
	}

	fnTable["owm"]["current"] = OwmAdaptCurrentWeather
	fnTable["owm"]["forecast"] = OwmAdaptForecast
	fnTable["wunderground"]["current"] = WundergroundAdaptCurrentWeather
	fnTable["myweather2"]["current"] = Myweather2AdaptCurrentWeather
	fnTable["forecast.io"]["current"] = ForecastioAdaptCurrentWeather
	fnTable["forecast.io"]["forecast"] = ForecastioAdaptForecast
	fnTable["worldweatheronline"]["current"] = WorldweatheronlineAdaptCurrentWeather
	fnTable["worldweatheronline"]["forecast"] = WorldweatheronlineAdaptForecast
	fnTable["accuweather"]["current"] = AccuweatherAdaptCurrentWeather
	fnTable["gismeteo"]["current"] = GismeteoAdaptCurrentWeather

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
