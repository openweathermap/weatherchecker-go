package adapters

import (
        "encoding/json"
        )

type LocationCoords struct {
    Longitude float32 `json:"lon"`
    Latitude float32 `json:"lat"`
}

type SystemInfo struct {
    Type int `json:"type"`
    Id int `json:"id"`
    Message float32 `json:"message"`
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
    Temp float32 `json:"temp"`
    Pressure float32 `json:"pressure"`
    Humidity int `json:"humidity"`
    TempMin float32 `json:"temp_min"`
    TempMax float32 `json:"temp_max"`
}

type WindInfo struct {
    Speed float32 `json:"speed"`
    Degree float32 `json:"deg"`
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

func decode (s string) WeatherHash {
    var byteString = []byte(s)
    var data = WeatherHash{}

    json.Unmarshal(byteString, &data)

    return data
}

func OwmAdaptCurrentWeather(jsonString string) MeasurementArray {
    var data = decode(jsonString)
    var measurements MeasurementArray

    measurements = append(measurements, MeasurementSchema{Temp:data.Main.Temp, Pressure:data.Main.Pressure, Wind:data.Wind.Speed})

    return measurements
}

func OwmAdaptForecastWeather(json_string string) MeasurementArray {
    var data = decode(json_string)
    var measurements MeasurementArray

    measurements = append(measurements, MeasurementSchema{Temp:data.Main.Temp, Pressure:data.Main.Pressure, Wind:data.Wind.Speed})

    return measurements
}
