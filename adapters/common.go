package adapters

type MeasurementSchema struct {
    Humidity float32
    Pressure float32
    Precipitation float32
    Temp float32
    Wind float32
}

type MeasurementArray []MeasurementSchema

func AdaptWeather(sourceName string, wtypeName string, data string) MeasurementArray {
    var fnTable = make(map[string](map[string]func(string)MeasurementArray))

    fnTable["OpenWeatherMap"] = make(map[string]func(string)MeasurementArray)
    fnTable["OpenWeatherMap"]["current"] = OwmAdaptCurrentWeather
    fnTable["OpenWeatherMap"]["forecast"] = OwmAdaptForecastWeather

    fnTable["Weather Underground"] = make(map[string]func(string)MeasurementArray)
    fnTable["Weather Underground"]["current"] = func(data string)MeasurementArray {return MeasurementArray{MeasurementSchema{}}}
    fnTable["Weather Underground"]["forecast"] = func(data string)MeasurementArray {return MeasurementArray{MeasurementSchema{}}}

    fnTable["MyWeather2"] = make(map[string]func(string)MeasurementArray)
    fnTable["MyWeather2"]["current"] = func(data string)MeasurementArray {return MeasurementArray{MeasurementSchema{}}}
    fnTable["MyWeather2"]["forecast"] = func(data string)MeasurementArray {return MeasurementArray{MeasurementSchema{}}}

    return fnTable[sourceName][wtypeName](data)
}
