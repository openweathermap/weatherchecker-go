package adapters

type MeasurementSchema struct {
    Humidity float32
    Pressure float32
    Precipitation float32
    Temp float32
    Wind float32
}

func Adapt_weather(source_name string, wtype_name string, data string) MeasurementSchema {
    var fnTable = make(map[string](map[string]func(string)MeasurementSchema))

    fnTable["OpenWeatherMap"] = make(map[string]func(string)MeasurementSchema)
    fnTable["OpenWeatherMap"]["current"] = Owm_adapt_current_weather
    fnTable["OpenWeatherMap"]["forecast"] = Owm_adapt_forecast_weather

    fnTable["Weather Underground"] = make(map[string]func(string)MeasurementSchema)
    fnTable["Weather Underground"]["current"] = func(data string)MeasurementSchema {return MeasurementSchema{}}
    fnTable["Weather Underground"]["forecast"] = func(data string)MeasurementSchema {return MeasurementSchema{}}

    return fnTable[source_name][wtype_name](data)
}
