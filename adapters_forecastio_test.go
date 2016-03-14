package main

import (
	"errors"
	"testing"
)

func TestForecastioAdaptCurrentWeather(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{"latitude":55.75,"longitude":37.62,"timezone":"Europe/Moscow","offset":3,"currently":{"time":1450793454,"summary":"Mostly Cloudy","icon":"partly-cloudy-night","precipIntensity":0.0051,"precipProbability":0.21,"precipType":"rain","temperature":39.98,"apparentTemperature":33.77,"dewPoint":37.77,"humidity":0.92,"windSpeed":9.65,"windBearing":276,"cloudCover":0.78,"pressure":1018,"ozone":303.91}}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(92), Precipitation: float64(0.0051), Pressure: float64(1018), Temp: float64(39.98), Wind: float64(9.65)}, Timestamp: int64(1450793454)}}
	result, resultErr := ForecastioAdaptCurrentWeather(s)

	if resultErr != nil {
		t.Errorf(resultErr.Error())
	}

	if len(expectation) != len(result) {
		err = compErr
	} else {
		for i, _ := range expectation {
			if result[i] != expectation[i] {
				err = compErr
			}
		}
	}

	if err != nil {
		t.Errorf(ErrorOut(expectation, result))
	}
}

func TestForecastioAdaptForecast(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{"latitude":55.75,"longitude":37.62,"timezone":"Europe/Moscow","offset":3,"hourly":{"summary":"Drizzle later this evening.","icon":"rain","data":[{"time":1450792800,"summary":"Drizzle","icon":"rain","precipIntensity":0.0051,"precipProbability":0.21,"precipType":"rain","temperature":40.05,"apparentTemperature":33.84,"dewPoint":37.71,"humidity":0.91,"windSpeed":9.67,"windBearing":276,"cloudCover":0.76,"pressure":1001.93,"ozone":303.78},{"time":1450796400,"summary":"Drizzle","icon":"rain","precipIntensity":0.0053,"precipProbability":0.22,"precipType":"rain","temperature":39.68,"apparentTemperature":33.44,"dewPoint":38.03,"humidity":0.94,"windSpeed":9.57,"windBearing":276,"cloudCover":0.87,"pressure":1002.25,"ozone":304.47}]}}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(91), Precipitation: float64(0.0051), Pressure: float64(1001.93), Temp: float64(40.05), Wind: 9.67}, Timestamp: int64(1450792800)}, MeasurementSchema{Data: Measurement{Humidity: float64(94), Precipitation: float64(0.0053), Pressure: float64(1002.25), Temp: float64(39.68), Wind: float64(9.57)}, Timestamp: int64(1450796400)}}
	result, resultErr := ForecastioAdaptForecast(s)

	if resultErr != nil {
		t.Errorf(ErrorOut(expectation, result))
	}

	if len(expectation) != len(result) {
		err = compErr
	} else {
		for i, _ := range expectation {
			if result[i] != expectation[i] {
				err = compErr
			}
		}
	}

	if err != nil {
		t.Errorf(ErrorOut(expectation, result))
	}
}
