package adapters

import (
        "encoding/json"
        )

type LocationCoords struct {
    Longitude float64 `json:"lon"`
    Latitude float64 `json:"lat"`
}

type SystemInfo struct {
    Type int `json:"type"`
    Id int `json:"id"`
    Message float64 `json:"message"`
    Country string `json:"country"`
    Sunrise int `json:"sunrise"`
    Sunset int `json:"sunset"`
}

type WeatherInfo struct {
    Id int `json:"id"`
    Main string `json:"main"`
    Description string `json:"description"`
    Icon string `json:"icon"`
}

type MainInfo struct {
    Temp float64 `json:"temp"`
    Pressure float64 `json:"pressure"`
    Humidity int `json:"humidity"`
    TempMin float64 `json:"temp_min"`
    TempMax float64 `json:"temp_max"`
}

type WindInfo struct {
    Speed float64 `json:"speed"`
    Degree float64 `json:"deg"`
}

type CloudInfo struct {
    All int `json:"all"`
}

type WeatherHash struct {
    Coord LocationCoords `json:"coord"`
    Code int `json:"cod"`
    Name string `json:"name"`
    Sys SystemInfo `json:"sys"`
    Weather []WeatherInfo `json:"weather"`
    Base string `json:"base"`
    Main MainInfo `json:"main"`
    Visibility int `json:"visibility"`
    Wind WindInfo `json:"wind"`
    Clouds CloudInfo `json:"clouds"`
    Timestamp int `json:"dt"`
    Id int `json:"id"`
}

func owmDecode (s string) WeatherHash {
    var byteString = []byte(s)
    var data = WeatherHash{}

    json.Unmarshal(byteString, &data)

    return data
}

func OwmAdaptCurrentWeather(jsonString string) MeasurementArray {
    var data = owmDecode(jsonString)
    var measurements MeasurementArray

    temp := float64(data.Main.Temp)
    pressure := float64(data.Main.Pressure)
    wind := float64(data.Wind.Speed)
    humidity := float64(0)
    precipitation := float64(0)

    measurements = append(measurements, MeasurementSchema{Humidity:humidity, Precipitation:precipitation, Pressure:pressure, Temp:temp, Wind:wind})

    return measurements
}
