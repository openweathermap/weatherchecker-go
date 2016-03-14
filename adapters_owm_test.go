package main

import (
	"errors"
	"testing"
)

func TestOwmAdaptCurrentWeather(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{"coord":{"lon":37.62,"lat":55.75},"sys":{"type":1,"id":7323,"message":0.01,"country":"RU","sunrise":1450158786,"sunset":1450184179},"weather":[{"id":800,"main":"Clear","description":"Sky is Clear","icon":"01n"}],"base":"cmc stations","main":{"temp":-4.65,"pressure":1018,"humidity":92,"temp_min":-6,"temp_max":-3},"wind":{"speed":4,"deg":240},"clouds":{"all":0},"dt":1450220400,"id":524901,"name":"Moscow","cod":200}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(92), Precipitation: float64(0), Pressure: float64(1018), Temp: float64(-4.65), Wind: float64(4)}, Timestamp: int64(1450220400)}}
	result, resultErr := OwmAdaptCurrentWeather(s)

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

func TestOwmAdaptForecast(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{"cod":"200","message":0.0167,"city":{"id":524901,"name":"Moscow","coord":{"lon":37.615555,"lat":55.75222},"country":"RU","population":0,"sys":{"population":0}},"cnt":42,"list":[{"dt":1450213200,"main":{"temp":-3.98,"temp_min":-4.59,"temp_max":-3.98,"pressure":1012.72,"sea_level":1033.89,"grnd_level":1012.72,"humidity":86,"temp_kf":0.61},"weather":[{"id":600,"main":"Snow","description":"light snow","icon":"13n"}],"clouds":{"all":48},"wind":{"speed":4.01,"deg":302.501},"snow":{"3h":0.033},"sys":{"pod":"n"},"dt_txt":"2015-12-15 21:00:00"},{"dt":1450224000,"main":{"temp":-4.26,"temp_min":-4.75,"temp_max":-4.26,"pressure":1012.18,"sea_level":1033.39,"grnd_level":1012.18,"humidity":87,"temp_kf":0.49},"weather":[{"id":600,"main":"Snow","description":"light snow","icon":"13n"}],"clouds":{"all":76},"wind":{"speed":3.77,"deg":280.502},"snow":{"3h":0.058},"sys":{"pod":"n"},"dt_txt":"2015-12-16 00:00:00"}]}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(86), Precipitation: float64(0.033), Pressure: float64(1012.72), Temp: float64(-3.98), Wind: 4.01}, Timestamp: int64(1450213200)}, MeasurementSchema{Data: Measurement{Humidity: float64(87), Precipitation: float64(0.058), Pressure: float64(1012.18), Temp: float64(-4.26), Wind: float64(3.77)}, Timestamp: int64(1450224000)}}
	result, resultErr := OwmAdaptForecast(s)

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
