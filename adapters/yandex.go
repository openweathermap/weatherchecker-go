package adapters

import (
	"encoding/xml"
	"time"
)

type YandexWeatherCondition struct {
	Code string `xml:"code,attr"`
}

type YandexWeatherStation struct {
	Lang     string `xml:"lang,attr"`
	Distance string `xml:"distance,attr"`
	Name     string `xml:",innerxml"`
}

type YandexWeatherData struct {
	Date             string                 `xml:"date,attr"`
	Station          []YandexWeatherStation `xml:"station"`
	ObservationTime  string                 `xml:"observation_time"`
	Uptime           string                 `xml:"uptime"`
	Temperature      int64                  `xml:"temperature"`
	WeatherCondition string                 `xml:"weather_condition"`
	WeatherType      string                 `xml:"weather_type"`
	WindDirection    string                 `xml:"wind_direction"`
	WindSpeed        float64                `xml:"wind_speed"`
	Humidity         int64                  `xml:"humidity"`
	PressureMbar     int64                  `xml:"mslp_pressure"`
}

type YandexWeatherExport struct {
	City       string `xml:"city,attr"`
	CityId     string `xml:"id,attr"`
	Country    string `xml:"country,attr"`
	CountryId  string `xml:"country_id,attr"`
	RegionCode string `xml:"region,attr"`
	GeoId      string `xml:"geoid,attr"`
	Longitude  string `xml:"lon,attr"`
	Latitude   string `xml:"lat,attr"`
	Zoom       string `xml:"zoom,attr"`
	SourceType string `xml:"source,attr"`

	Current   YandexWeatherData   `xml:"fact"`
	Yesterday YandexWeatherData   `xml:"yesterday"`
	Forecast  []YandexWeatherData `xml:"day"`
}

func yandexDecode(s string) (data YandexWeatherExport, err error) {
	byteString := []byte(s)

	err = xml.Unmarshal(byteString, &data)

	return data, err
}

func YandexAdaptCurrentWeather(xmlString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(xmlString)
			err = AdapterPanicErr
		}
	}()

	data, decodeErr := yandexDecode(xmlString)

	if decodeErr != nil {
		return AdaptStub(xmlString), decodeErr
	}

	// Yandex does not provide time zone information for its timestamps - Geocoding service is required
	/*
		timeValue, timeErr := time.Parse("RFC3339", data.Current.ObservationTime)

		if timeErr != nil {
			return AdaptStub(xmlString), timeErr
		}

		dt := timeValue.Unix()
	*/

	dt := time.Now().Unix()

	temp := float64(data.Current.Temperature)
	pressure := float64(data.Current.PressureMbar)
	wind := float64(data.Current.WindSpeed)
	humidity := float64(data.Current.Humidity)
	precipitation := float64(0)

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}
