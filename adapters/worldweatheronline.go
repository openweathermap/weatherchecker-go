package adapters

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type WorldweatheronlineCondition struct {
	CloudCover          string `json:"cloudcover"`
	FeelsLikeC          string `json:"FeelsLikeC"`
	FeelsLikeF          string `json:"FeelsLikeF"`
	Humidity            string `json:"humidity"`
	ObservationTime     string `json:"observation_time"`
	PrecipMM            string `json:"precipMM"`
	Pressure            string `json:"pressure"`
	TempC               string `json:"temp_C"`
	TempF               string `json:"temp_F"`
	Visibility          string `json:"visibility"`
	WeatherCode         string `json:"weatherCode"`
	WeatherDescriptions []struct {
		Value string `json:"value"`
	} `json:"weatherDesc"`
	WeatherIcons []struct {
		Value string `json:"value"`
	} `json:"weatherIconUrl"`
	WindDir16Point string `json:"winddir16Point"`
	WindDirDegree  string `json:"winddirDegree"`
	WindSpeedKmph  string `json:"windspeedKmph"`
	WindSpeedMiles string `json:"windspeedMiles"`
}

type WorldweatheronlineResponseData struct {
	CurrentCondition []WorldweatheronlineCondition `json:"current_condition"`
}

type WorldweatheronlineResponse struct {
	Data WorldweatheronlineResponseData `json:"data"`
}

func worldweatheronlineDecode(s string) WorldweatheronlineResponse {
	var data WorldweatheronlineResponse

	var byteString = []byte(s)

	json.Unmarshal(byteString, &data)

	return data
}

func WorldweatheronlineAdaptCurrentWeather(jsonString string) MeasurementArray {
	var data = worldweatheronlineDecode(jsonString)
	var measurements MeasurementArray

	dt := time.Now()

	humidity_raw := strings.TrimSpace(data.Data.CurrentCondition[0].Humidity)
	pressure_raw := strings.TrimSpace(data.Data.CurrentCondition[0].Pressure)
	precipitation_raw := strings.TrimSpace(data.Data.CurrentCondition[0].PrecipMM)
	temp_raw := strings.TrimSpace(data.Data.CurrentCondition[0].TempC)
	wind_raw := strings.TrimSpace(data.Data.CurrentCondition[0].WindSpeedKmph)

	humidity, _ := strconv.ParseFloat(humidity_raw, 64)
	pressure, _ := strconv.ParseFloat(pressure_raw, 64)
	precipitation, _ := strconv.ParseFloat(precipitation_raw, 64)
	temp, _ := strconv.ParseFloat(temp_raw, 64)
	wind, _ := strconv.ParseFloat(wind_raw, 64)

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements
}
