package adapters

type MeasurementSchema struct {
    Humidity float32
    Pressure float32
    Precipitation float32
    Temp float32
    Wind float32
}

type MeasurementArray []MeasurementSchema

func AdaptStub (s string) MeasurementArray {return MeasurementArray{MeasurementSchema{}}}

func AdaptWeather(sourceName string, wtypeName string, data string) MeasurementArray {
    var adaptFunc func(string)MeasurementArray
    var fnTable = make(map[string](map[string]func(string)MeasurementArray))

    fnTable["OpenWeatherMap"] = make(map[string]func(string)MeasurementArray)
    fnTable["OpenWeatherMap"]["current"] = OwmAdaptCurrentWeather

    fnTable["Weather Underground"] = make(map[string]func(string)MeasurementArray)
    fnTable["Weather Underground"]["current"] = WundergroundAdaptCurrentWeather

    fnTable["MyWeather2"] = make(map[string]func(string)MeasurementArray)
    fnTable["MyWeather2"]["current"] = Myweather2AdaptCurrentWeather

    adaptFunc = AdaptStub

    _, p_ok := fnTable[sourceName]
    if p_ok == true {
        storedFunc, f_ok := fnTable[sourceName][wtypeName]
        if f_ok == true {
            adaptFunc = storedFunc
        }
    }

    return adaptFunc(data)
}
