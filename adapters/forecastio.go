package adapters

import (
	"encoding/json"
)

type ForecastioWeatherBase struct {
	Time              int     `json:"time"`
	Summary           string  `json:"summary"`
	Icon              string  `json:"icon"`
	PrecipIntensity   float64 `json:"precipIntensity"`
	PrecipProbability float64 `json:"precipProbability"`
	PrecipType        string  `json:"precipType"`
	DewPoint          float64 `json:"dewPoint"`
	Humidity          float64 `json:"humidity"`
	WindSpeed         float64 `json:"windSpeed"`
	WindBearing       float64 `json:"windBearing"`
	CloudCover        float64 `json:"cloudCover"`
	Pressure          float64 `json:"pressure"`
	Ozone             float64 `json:"ozone"`
}

type ForecastioCurrentWeather struct {
	ForecastioWeatherBase
	Temperature         float64 `json:"temperature"`
	ApparentTemperature float64 `json:"apparentTemperature"`
}

type ForecastioHourlyWeatherDataEntry struct {
	ForecastioWeatherBase
	Temperature         float64 `json:"temperature"`
	ApparentTemperature float64 `json:"apparentTemperature"`
}

type ForecastioHourlyWeather struct {
	Summary string                             `json:"summary"`
	Icon    string                             `json:"icon"`
	Data    []ForecastioHourlyWeatherDataEntry `json:"data"`
}

type ForecastioDailyWeatherDataEntry struct {
	ForecastioWeatherBase
	MoonPhase                  float64 `json:"moonPhase"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime     int     `json:"precipIntensityMaxTime"`
	PrecipAccumulation         float64 `json:"precipAccumulation"`
	SunriseTime                int     `json:"sunriseTime"`
	SunsetTime                 int     `json:"sunsetTime"`
	TemperatureMin             float64 `json:"temperatureMin"`
	TemperatureMax             float64 `json:"temperatureMax"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime int     `json:"apparentTemperatureMinTime"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime int     `json:"apparentTemperatureMaxTime"`
}

type ForecastioDailyWeather struct {
	Summary string                            `json:"summary"`
	Icon    string                            `json:"icon"`
	Data    []ForecastioDailyWeatherDataEntry `json:"data"`
}

type ForecastioWeatherResponse struct {
	Latitude  float64                  `json:"latitude"`
	Longitude float64                  `json:"longitude"`
	Timezone  string                   `json:"timezone"`
	Offset    int                      `json:"offset"`
	Current   ForecastioCurrentWeather `json:"currently"`
	Hourly    ForecastioHourlyWeather  `json:"hourly"`
	Daily     ForecastioDailyWeather   `json:"daily"`
}

func forecastioDecode(s string) (data ForecastioWeatherResponse, err error) {
	var byteString = []byte(s)

	err = json.Unmarshal(byteString, &data)

	return data, err
}

func ForecastioAdaptCurrentWeather(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = AdapterPanicErr
		}
	}()
	data, decodeErr := forecastioDecode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	dt := int64(data.Current.Time)

	humidity_raw := data.Current.Humidity
	pressure_raw := data.Current.Pressure
	precipitation_raw := data.Current.PrecipIntensity
	temp_raw := data.Current.Temperature
	wind_raw := data.Current.WindSpeed

	humidity := float64(humidity_raw) * 100
	pressure := float64(pressure_raw)
	precipitation := float64(precipitation_raw)
	temp := temp_raw
	wind := wind_raw

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}

func ForecastioAdaptForecast(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
		}
	}()
	data, decodeErr := forecastioDecode(jsonString)

	if decodeErr != nil {
		panic(decodeErr.Error())
	}

	for _, entry := range data.Hourly.Data {
		dt := int64(entry.Time)

		humidity_raw := entry.Humidity
		pressure_raw := entry.Pressure
		precipitation_raw := entry.PrecipIntensity
		temp_raw := entry.Temperature
		wind_raw := entry.WindSpeed

		humidity := float64(humidity_raw) * 100
		pressure := float64(pressure_raw)
		precipitation := float64(precipitation_raw)
		temp := temp_raw
		wind := wind_raw

		measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})
	}

	return measurements, err
}
