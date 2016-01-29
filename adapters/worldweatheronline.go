package adapters

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/owm-inc/weatherchecker-go/common"
)

type WorldweatheronlineConditionBase struct {
	CloudCover          string `json:"cloudcover"`
	FeelsLikeC          string `json:"FeelsLikeC"`
	FeelsLikeF          string `json:"FeelsLikeF"`
	Humidity            string `json:"humidity"`
	ObservationTime     string `json:"observation_time"`
	PrecipMM            string `json:"precipMM"`
	Pressure            string `json:"pressure"`
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

type WorldweatheronlineCurrentMeasurement struct {
	WorldweatheronlineConditionBase
	TempC string `json:"temp_C"`
	TempF string `json:"temp_F"`
}

type WorldweatheronlineForecastMeasurement struct {
	WorldweatheronlineConditionBase
	TempC         string `json:"tempC"`
	TempF         string `json:"tempF"`
	UTCdate       string `json:"UTCdate"`
	UTCtime       string `json:"UTCtime"`
	WindGustKmph  string `json:"WindGustKmph"`
	WindGustMiles string `json:"WindGustMiles"`
}

type WorldweatheronlineForecast struct {
	Date         string                                  `json:"date"`
	Measurements []WorldweatheronlineForecastMeasurement `json:"hourly"`
}

type WorldweatheronlineResponseData struct {
	CurrentCondition []WorldweatheronlineCurrentMeasurement `json:"current_condition"`
	WeatherForecast  []WorldweatheronlineForecast           `json:"weather"`
}

type WorldweatheronlineResponse struct {
	Data WorldweatheronlineResponseData `json:"data"`
}

func worldweatheronlineDecode(s string) (data WorldweatheronlineResponse, err error) {
	byteString := []byte(s)

	err = json.Unmarshal(byteString, &data)

	return data, err
}

func WorldweatheronlineAdaptCurrentWeather(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = common.AdapterPanicErr
		}
	}()

	data, decodeErr := worldweatheronlineDecode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	dt, _ := strconv.ParseInt(data.Data.CurrentCondition[0].ObservationTime, 10, 64)

	humidity_raw := strings.TrimSpace(data.Data.CurrentCondition[0].Humidity)
	pressure_raw := strings.TrimSpace(data.Data.CurrentCondition[0].Pressure)
	precipitation_raw := strings.TrimSpace(data.Data.CurrentCondition[0].PrecipMM)
	temp_raw := strings.TrimSpace(data.Data.CurrentCondition[0].TempC)
	wind_raw := strings.TrimSpace(data.Data.CurrentCondition[0].WindSpeedKmph)

	humidity, _ := strconv.ParseFloat(humidity_raw, 64)
	pressure, _ := strconv.ParseFloat(pressure_raw, 64)
	precipitation, _ := strconv.ParseFloat(precipitation_raw, 64)
	temp, _ := strconv.ParseFloat(temp_raw, 64)
	wind_kph, _ := strconv.ParseFloat(wind_raw, 64)
	wind, _ := convertUnits(float64(wind_kph), "kph")

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}

func WorldweatheronlineAdaptForecast(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = common.AdapterPanicErr
		}
	}()
	data, decodeErr := worldweatheronlineDecode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	for _, day_entry := range data.Data.WeatherForecast {
		for _, entry := range day_entry.Measurements {
			dateSplit := strings.Split(entry.UTCdate, "-")
			y, _ := strconv.ParseInt(dateSplit[0], 10, 64)
			m, _ := strconv.ParseInt(dateSplit[1], 10, 64)
			d, _ := strconv.ParseInt(dateSplit[2], 10, 64)
			h_military, _ := strconv.ParseFloat(strings.TrimSpace(entry.UTCtime), 64)
			h := int(math.Floor(h_military / 100))

			dtt := time.Date(int(y), time.Month(m), int(d), int(h), 0, 0, 0, time.UTC)

			dt := dtt.Unix()

			humidity_raw := strings.TrimSpace(entry.Humidity)
			pressure_raw := strings.TrimSpace(entry.Pressure)
			precipitation_raw := strings.TrimSpace(entry.PrecipMM)
			temp_raw := strings.TrimSpace(entry.TempC)
			wind_raw := strings.TrimSpace(entry.WindSpeedKmph)

			humidity, _ := strconv.ParseFloat(humidity_raw, 64)
			pressure, _ := strconv.ParseFloat(pressure_raw, 64)
			precipitation, _ := strconv.ParseFloat(precipitation_raw, 64)
			temp, _ := strconv.ParseFloat(temp_raw, 64)
			wind_kph, _ := strconv.ParseFloat(wind_raw, 64)
			wind, _ := convertUnits(float64(wind_kph), "kph")

			measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})
		}
	}

	return measurements, err
}
