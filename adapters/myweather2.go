package adapters

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/owm-inc/weatherchecker-go/common"
)

type Myweather2WindInfo struct {
	Speed string `json:"speed"`
}

type Myweather2WeatherData struct {
	Temp     string               `json:"temp"`
	Humidity string               `json:"humidity"`
	Pressure string               `json:"pressure"`
	Wind     []Myweather2WindInfo `json:"wind"`
}

type Myweather2Weather struct {
	CurrentWeather []Myweather2WeatherData `json:"curren_weather"`
}

type Myweather2Response struct {
	Weather Myweather2Weather `json:"weather"`
}

func myweather2Decode(s string) (data Myweather2Response, err error) {
	byteString := []byte(s)

	err = json.Unmarshal(byteString, &data)

	return data, err
}

func Myweather2AdaptCurrentWeather(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = common.AdapterPanicErr
		}
	}()
	data, decodeErr := myweather2Decode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	dt := time.Now().Unix()

	humidity_raw := strings.TrimSpace(data.Weather.CurrentWeather[0].Humidity)
	pressure_raw := strings.TrimSpace(data.Weather.CurrentWeather[0].Pressure)
	temp_raw := strings.TrimSpace(data.Weather.CurrentWeather[0].Temp)
	wind_raw := strings.TrimSpace(data.Weather.CurrentWeather[0].Wind[0].Speed)

	humidity, _ := strconv.ParseFloat(humidity_raw, 64)
	pressure, _ := strconv.ParseFloat(pressure_raw, 64)
	temp, _ := strconv.ParseFloat(temp_raw, 64)
	wind, _ := strconv.ParseFloat(wind_raw, 64)

	precipitation := float64(0)

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}
