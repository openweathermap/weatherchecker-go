package adapters

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/owm-inc/weatherchecker-go/common"
)

type WundergroundResponseStruct struct {
	Version        string `json:"version"`
	TermsOfService string `json:"termsofService"`
	Features       map[string]int
}

type WundergroundImageStruct struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type WundergroundLocationStruct struct {
	Full           string `json:"full"`
	City           string `json:"city"`
	State          string `json:"state"`
	StateName      string `json:"state_name"`
	Country        string `json:"country"`
	CountryIso3166 string `json:"country_iso3166"`
	Zip            string `json:"zip"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Elevation      string `json:"elevation"`
}

type WundergroundCurrentObservationStruct struct {
	Image                 WundergroundImageStruct    `json:"image"`
	DisplayLocation       WundergroundLocationStruct `json:"display_location"`
	ObservationLocation   WundergroundLocationStruct `json:"observation_location"`
	StationId             string                     `json:"station_id"`
	ObservationTime       string                     `json:"observation_time"`
	ObservationTimeRfc822 string                     `json:"observation_time_tfc822"`
	ObservationEpoch      string                     `json:"observation_epoch"`
	LocalTimeRfc822       string                     `json:"local_time_rfc822"`
	LocalEpoch            string                     `json:"local_epoch"`
	LocalTzShort          string                     `json:"local_tz_short"`
	LocalTzLong           string                     `json:"local_tz_long"`
	LocalTzOffset         string                     `json:"local_tz_offset"`
	Weather               string                     `json:"weather"`
	TemperatureString     string                     `json:"temperature_string"`
	TempF                 float64                    `json:"temp_f"`
	TempC                 float64                    `json:"temp_c"`
	RelativeHumidity      string                     `json:"relative_humidity"`
	WindString            string                     `json:"wind_string"`
	WindDir               string                     `json:"wind_dir"`
	WindDegrees           int                        `json:"wind_degrees"`
	WindMph               float64                    `json:"wind_mph"`
	WindGustMph           string                     `json:"wind_gust_mph"`
	WindKph               float64                    `json:"wind_kph"`
	WindGustKph           string                     `json:"wind_gust_kph"`
	PressureMb            string                     `json:"pressure_mb"`
	PressureIn            string                     `json:"pressure_in"`
	PressureTrend         string                     `json:"pressure_trend"`
	DewpointString        string                     `json:"dewpoint_string"`
	DewpointF             float64                    `json:"dewpoint_f"`
	DewpointC             float64                    `json:"dewpoint_c"`
	HeatIndexString       string                     `json:"heat_index_string"`
	HeatIndexF            string                     `json:"heat_index_f"`
	HeatIndexC            string                     `json:"heat_index_c"`
	WindChillString       string                     `json:"windchill_string"`
	WindChillF            string                     `json:"windchill_f"`
	WindChillC            string                     `json:"windchill_c"`
	FeelsLikeString       string                     `json:"feelslike_string"`
	FeelsLikeF            string                     `json:"feelslike_f"`
	FeelsLikeC            string                     `json:"feelslike_c"`
	VisibilityMi          string                     `json:"visibility_mi"`
	VisibilityKm          string                     `json:"visibility_km"`
	SolarRadiation        string                     `json:"solarradiation"`
	UV                    string                     `json:"UV"`
	Precip1hrString       string                     `json:"precip_1hr_string"`
	Precip1hrIn           string                     `json:"precip_1hr_in"`
	Precip1hrMetric       string                     `json:"precip_1hr_metric"`
	PrecipTodayString     string                     `json:"precip_today_string"`
	PrecipTodayIn         string                     `json:"precip_today_in"`
	PrecipTodayMetric     string                     `json:"precip_today_metric"`
	Icon                  string                     `json:"icon"`
}

type WundergroundWeatherStruct struct {
	Response           WundergroundResponseStruct           `json:"response"`
	CurrentObservation WundergroundCurrentObservationStruct `json:"current_observation"`
}

func wundergroundDecode(s string) (data WundergroundWeatherStruct, err error) {
	byteString := []byte(s)

	err = json.Unmarshal(byteString, &data)

	return data, err
}

// WundergroundAdaptCurrentWeather normalizes Weather Underground current response.
func WundergroundAdaptCurrentWeather(jsonString string) (measurements MeasurementArray, err error) {
	defer func() {
		if r := recover(); r != nil {
			measurements = AdaptStub(jsonString)
			err = common.AdapterPanicErr
		}
	}()
	data, decodeErr := wundergroundDecode(jsonString)

	if decodeErr != nil {
		return AdaptStub(jsonString), decodeErr
	}

	dt, _ := strconv.ParseInt(data.CurrentObservation.ObservationEpoch, 10, 64)

	if dt == 0 {
		return nil, common.MalformedResponse
	}

	humidityRaw := strings.TrimRight(strings.TrimSpace(data.CurrentObservation.RelativeHumidity), "%")
	pressureRaw := strings.TrimSpace(data.CurrentObservation.PressureMb)
	precipitationRaw := strings.TrimSpace(data.CurrentObservation.PrecipTodayMetric)
	tempRaw := data.CurrentObservation.TempC
	windRaw := data.CurrentObservation.WindKph

	humidity, _ := strconv.ParseFloat(humidityRaw, 64)
	pressure, _ := strconv.ParseFloat(pressureRaw, 64)
	precipitation, _ := strconv.ParseFloat(precipitationRaw, 64)
	temp := float64(tempRaw)
	wind, _ := convertUnits(float64(windRaw), "kph")

	measurements = append(measurements, MeasurementSchema{Data: Measurement{Humidity: humidity, Precipitation: precipitation, Pressure: pressure, Temp: temp, Wind: wind}, Timestamp: dt})

	return measurements, err
}
