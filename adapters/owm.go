package adapters

import (
        "encoding/json"
        )

type OwmLocationCoords struct {
    Longitude float64 `json:"lon"`
    Latitude float64 `json:"lat"`
}

type OwmSystemInfo struct {
    Type int `json:"type"`
    Id int `json:"id"`
    Message float64 `json:"message"`
    Country string `json:"country"`
    Sunrise int `json:"sunrise"`
    Sunset int `json:"sunset"`
}

type OwmWeatherInfo struct {
    Id int `json:"id"`
    Main string `json:"main"`
    Description string `json:"description"`
    Icon string `json:"icon"`
}

type OwmMainInfo struct {
    Temp float64 `json:"temp"`
    Pressure float64 `json:"pressure"`
    Humidity int `json:"humidity"`
    TempMin float64 `json:"temp_min"`
    TempMax float64 `json:"temp_max"`
}

type OwmWindInfo struct {
    Speed float64 `json:"speed"`
    Degree float64 `json:"deg"`
}

type OwmCloudInfo struct {
    All int `json:"all"`
}

type OwmWeatherStruct struct {
    Coord OwmLocationCoords `json:"coord"`
    Code int `json:"cod"`
    Name string `json:"name"`
    Sys OwmSystemInfo `json:"sys"`
    Weather []OwmWeatherInfo `json:"weather"`
    Base string `json:"base"`
    Main OwmMainInfo `json:"main"`
    Visibility int `json:"visibility"`
    Wind OwmWindInfo `json:"wind"`
    Clouds OwmCloudInfo `json:"clouds"`
    Timestamp int `json:"dt"`
    Id int `json:"id"`
}

func owmDecode (s string) OwmWeatherStruct {
    var data OwmWeatherStruct

    var byteString = []byte(s)

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
