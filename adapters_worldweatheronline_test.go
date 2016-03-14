package main

import (
	"errors"
	"testing"
)

func TestWorldweatheronlineAdaptCurrentWeather(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{ "data": { "current_condition": [ {"cloudcover": "75", "FeelsLikeC": "0", "FeelsLikeF": "32", "humidity": "81", "observation_time": "1450955339", "precipMM": "2.5", "pressure": "1003", "temp_C": "5", "temp_F": "41", "visibility": "10", "weatherCode": "116",  "weatherDesc": [ {"value": "Partly Cloudy" } ],  "weatherIconUrl": [ {"value": "http:\/\/cdn.worldweatheronline.net\/images\/wsymbols01_png_64\/wsymbol_0002_sunny_intervals.png" } ], "winddir16Point": "W", "winddirDegree": "270", "windspeedKmph": "30", "windspeedMiles": "19" } ],  "request": [ {"query": "Lat 55.75 and Lon 37.62", "type": "LatLon" } ] }}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(81), Precipitation: float64(2.5), Pressure: float64(1003), Temp: float64(5), Wind: float64(8.333333333333334)}, Timestamp: int64(1450955339)}}
	result, resultErr := WorldweatheronlineAdaptCurrentWeather(s)

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

func TestWorldweatheronlineAdaptForecast(t *testing.T) {
	var err error
	compErr := errors.New("Mismatch.")
	s := `{ "data": { "weather": [ {"astronomy":[{"moonrise":"04:38 PM","moonset":"08:11 AM","sunrise":"08:59 AM","sunset":"04:00 PM"}],"date":"2015-12-25","hourly":[{"chanceoffog":"0","chanceoffrost":"89","chanceofhightemp":"0","chanceofovercast":"90","chanceofrain":"76","chanceofremdry":"0","chanceofsnow":"76","chanceofsunshine":"0","chanceofthunder":"0","chanceofwindy":"0","cloudcover":"90","DewPointC":"-1","DewPointF":"31","FeelsLikeC":"-4","FeelsLikeF":"25","HeatIndexC":"0","HeatIndexF":"33","humidity":"93","precipMM":"0.5","pressure":"1021","tempC":"5","tempF":"33","time":"200","UTCdate":"2015-12-25","UTCtime":"0","visibility":"10","weatherCode":"122","weatherDesc":[{"value":"Overcast"}],"weatherIconUrl":[{"value":"http://cdn.worldweatheronline.net/images/wsymbols01_png_64/wsymbol_0004_black_low_cloud.png"}],"WindChillC":"-4","WindChillF":"25","winddir16Point":"NW","winddirDegree":"321","WindGustKmph":"24","WindGustMiles":"15","windspeedKmph":"14","windspeedMiles":"9"},{"chanceoffog":"0","chanceoffrost":"89","chanceofhightemp":"0","chanceofovercast":"53","chanceofrain":"1","chanceofremdry":"0","chanceofsnow":"0","chanceofsunshine":"0","chanceofthunder":"0","chanceofwindy":"0","cloudcover":"92","DewPointC":"-1","DewPointF":"30","FeelsLikeC":"-4","FeelsLikeF":"24","HeatIndexC":"0","HeatIndexF":"32","humidity":"93","precipMM":"0.3","pressure":"1022","tempC":"4","tempF":"32","time":"500","UTCdate":"2015-12-25","UTCtime":"300","visibility":"10","weatherCode":"122","weatherDesc":[{"value":"Overcast"}],"weatherIconUrl":[{"value":"http://cdn.worldweatheronline.net/images/wsymbols01_png_64/wsymbol_0004_black_low_cloud.png"}],"WindChillC":"-4","WindChillF":"24","winddir16Point":"NW","winddirDegree":"310","WindGustKmph":"25","WindGustMiles":"16","windspeedKmph":"19","windspeedMiles":"9"}],"maxtempC":"3","maxtempF":"37","mintempC":"-1","mintempF":"30","uvIndex":"0"} ] }}`
	expectation := MeasurementArray{MeasurementSchema{Data: Measurement{Humidity: float64(93), Precipitation: float64(0.5), Pressure: float64(1021), Temp: float64(5), Wind: 3.888888888888889}, Timestamp: int64(1451001600)}, MeasurementSchema{Data: Measurement{Humidity: float64(93), Precipitation: float64(0.3), Pressure: float64(1022), Temp: float64(4), Wind: float64(5.277777777777778)}, Timestamp: int64(1451012400)}}
	result, resultErr := WorldweatheronlineAdaptForecast(s)

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
