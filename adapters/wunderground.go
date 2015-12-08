package adapters

import (
        "encoding/json"
        "strconv"
        "strings"
        )

type WundergroundResponseStruct struct {
    Version string `json:"version"`
    TermsOfService string `json:"termsofService"`
    Features map[string]int
}

type WundergroundImageStruct struct {
    Url string `json:"url"`
    Title string `json:"title"`
    Link string `json:"link"`
}

type WundergroundLocationStruct struct {
    Full string `json:"full"`
    City string `json:"city"`
    State string `json:"state"`
    StateName string `json:"state_name"`
    Country string `json:"country"`
    CountryIso3166 string `json:"country_iso3166"`
    Zip string `json:"zip"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Elevation string `json:"elevation"`
}

type WundergroundCurrentObservationStruct struct {
    Image WundergroundImageStruct                   `json:"image"`
    DisplayLocation WundergroundLocationStruct      `json:"display_location"`
    ObservationLocation WundergroundLocationStruct  `json:"observation_location"`
    StationId string                                `json:"station_id"`
    ObservationTime string                          `json:"observation_time"`
    ObservationTimeRfc822 string                    `json:"observation_time_tfc822"`
    ObservationEpoch string                         `json:"observation_epoch"`
    LocalTimeRfc822 string                          `json:"local_time_rfc822"`
    LocalEpoch string                               `json:"local_epoch"`
    LocalTzShort string                             `json:"local_tz_short"`
    LocalTzLong string                              `json:"local_tz_long"`
    LocalTzOffset string                            `json:"local_tz_offset"`
    Weather string                                  `json:"weather"`
    TemperatureString string                        `json:"temperature_string"`
    TempF float32                                   `json:"temp_f"`
    TempC float32                                   `json:"temp_c"`
    RelativeHumidity string                         `json:"relative_humidity"`
    WindString string                               `json:"wind_string"`
    WindDir string                                  `json:"wind_dir"`
    WindDegrees int                                 `json:"wind_degrees"`
    WindMph float32                                 `json:"wind_mph"`
    WindGustMph string                              `json:"wind_gust_mph"`
    WindKph float32                                 `json:"wind_kph"`
    WindGustKph string                              `json:"wind_gust_kph"`
    PressureMb string                               `json:"pressure_mb"`
    PressureIn string                               `json:"pressure_in"`
    PressureTrend string                            `json:"pressure_trend"`
    DewpointString string                           `json:"dewpoint_string"`
    DewpointF int                                   `json:"dewpoint_f"`
    DewpointC int                                   `json:"dewpoint_c"`
    HeatIndexString string                          `json:"heat_index_string"`
    HeatIndexF string                               `json:"heat_index_f"`
    HeatIndexC string                               `json:"heat_index_c"`
    WindChillString string                          `json:"windchill_string"`
    WindChillF string                               `json:"windchill_f"`
    WindChillC string                               `json:"windchill_c"`
    FeelsLikeString string                          `json:"feelslike_string"`
    FeelsLikeF string                               `json:"feelslike_f"`
    FeelsLikeC string                               `json:"feelslike_c"`
    VisibilityMi string                             `json:"visibility_mi"`
    VisibilityKm string                             `json:"visibility_km"`
    SolarRadiation string                           `json:"solarradiation"`
    UV string                                       `json:"UV"`
    Precip1hrString string                          `json:"precip_1hr_string"`
    Precip1hrIn string                              `json:"precip_1hr_in"`
    Precip1hrMetric string                          `json:"precip_1hr_metric"`
    PrecipTodayString string                        `json:"precip_today_string"`
    PrecipTodayIn string                            `json:"precip_today_in"`
    PrecipTodayMetric string                        `json:"precip_today_metric"`
    Icon string                                     `json:"icon"`
    IconUrl string                                  `json:"icon_url"`
    ForecastUrl string                              `json:"forecast_url"`
    HistoryUrl string                               `json:"history_url"`
    ObUrl string                                    `json:"ob_url"`
}

type WundergroundWeatherStruct struct {
    Response WundergroundResponseStruct `json:"response"`
    CurrentObservation WundergroundCurrentObservationStruct `json:"current_observation"`
}

func wundergroundDecode (s string) WundergroundWeatherStruct {
    var byteString = []byte(s)
    var data = WundergroundWeatherStruct{}

    json.Unmarshal(byteString, &data)

    return data
}

func WundergroundAdaptCurrentWeather(jsonString string) MeasurementArray {
    var data = wundergroundDecode(jsonString)
    var measurements MeasurementArray

    humidity_raw := strings.TrimRight(strings.TrimSpace(data.CurrentObservation.RelativeHumidity), "%")
    pressure_raw := strings.TrimSpace(data.CurrentObservation.PressureMb)
    precipitation_raw := strings.TrimSpace(data.CurrentObservation.PrecipTodayMetric)
    temp_raw := data.CurrentObservation.TempC
    wind_raw := data.CurrentObservation.WindKph

    humidityF64, _ := strconv.ParseFloat(humidity_raw, 32)
    pressureF64, _ := strconv.ParseFloat(pressure_raw, 32)
    precipitationF64, _ := strconv.ParseFloat(precipitation_raw, 32)

    pressure := float32(pressureF64)
    humidity := float32(humidityF64)
    precipitation := float32(precipitationF64)
    temp := float32(temp_raw)
    wind := float32(wind_raw)

    measurements = append(measurements, MeasurementSchema{Humidity:humidity, Precipitation:precipitation, Pressure:pressure, Temp:temp, Wind:wind})

    return measurements
}