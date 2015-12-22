package adapters

import (
	"encoding/json"
)

type ForecastioWeatherBase struct {
	Time              int     `json:"time"`
	Summary           string  `json:"summary"`
	Icon              string  `json:"icon"`
	PrecipIntensity   int     `json:"precipIntensity"`
	PrecipProbability int     `json:"precipProbability"`
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

type ForecastioDailyWeatherDataEntry struct {
	ForecastioWeatherBase
	MoonPhase                  float64 `json:"moonPhase"`
	PrecipIntensity            float64 `json:"precipIntensity"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime     int     `json:"precipIntensityMaxTime"`
	PrecipType                 string  `json:"precipType"`
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
	Daily     ForecastioDailyWeather   `json:"daily"`
}

func forecastioDecode(s string) ForecastioWeatherResponse {
	var data ForecastioWeatherResponse

	var byteString = []byte(s)

	json.Unmarshal(byteString, &data)

	return data
}

func ForecastioAdaptCurrentWeather(jsonString string) (measurements MeasurementArray) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
		}
	}()
	var data = forecastioDecode(jsonString)

	dt := int64(data.Current.Time)

	humidity_raw := data.Current.Humidity
	pressure_raw := data.Current.Pressure
	precipitation_raw := data.Current.PrecipIntensity
	temp_raw := data.Current.Temperature
	wind_raw := data.Current.WindSpeed

	humidity := float64(humidity_raw)
	pressure := float64(pressure_raw)
	precipitation := float64(precipitation_raw)
	temp := float64((temp_raw - 32) * 5 / 9)
	wind := float64(wind_raw / 2.23)

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements
}
