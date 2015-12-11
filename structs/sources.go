package structs

import (
        "os"
        )

type SourceEntry struct {
    Name string `json:"name"`
    Urls map[string]string `json:"urls"`
    Keys Keyring
}

func CreateSources() []SourceEntry {
    var sources []SourceEntry
    var keys Keyring
    var urls map[string]string
    var entry SourceEntry

    keys = Keyring{Key:os.Getenv("OWM_KEY")}
    urls = map[string]string {"current":`http://api.openweathermap.org/data/2.5/weather?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&units=metric`,
                              "forecast":`http://api.openweathermap.org/data/2.5/forecast/daily?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&cnt=10&mode=json&units=metric`}
    entry = SourceEntry{Name:"OpenWeatherMap", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("WUNDERGROUND_KEY")}
    urls = map[string]string {"current": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/conditions/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`,
                              "forecast": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/forecast10day/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`}
    entry = SourceEntry{Name:"Weather Underground", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("MYWEATHER2_KEY"), Uref:os.Getenv("MYWEATHER2_UREF")}
    urls = map[string]string {"current": `http://www.myweather2.com/developer/forecast.ashx?uac={{.Source.Keys.Key}}&query={{.Location.Latitude}},{{.Location.Longitude}}&temp_unit=c&output=json&ws_unit=kph`,
                              "forecast": `http://www.myweather2.com/developer/weather.ashx?uac={{.Source.Keys.Key}}&uref={{.Source.Keys.Uref}}&query={{.Location.Latitude}},{{.Location.Longitude}}&output=json&temp_unit=c&ws_unit=kph`}
    entry = SourceEntry{Name:"MyWeather2", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("FORECASTIO_KEY")}
    urls = map[string]string {"current": `https://api.forecast.io/forecast/{{.Source.Keys.Key}}/{{.Location.Latitude}},{{.Location.Longitude}}`,
                              "forecast": `https://api.forecast.io/forecast/{{.Source.Keys.Key}}/{{.Location.Latitude}},{{.Location.Longitude}}`}
    entry = SourceEntry{Name:"Forecast.io", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("WORLDWEATHERONLINE_KEY")}
    urls = map[string]string {"current": `http://api.worldweatheronline.com/free/v2/weather.ashx?key={{.Source.Keys.Key}}&q={{.Location.Latitude}},{{.Location.Longitude}}&format=json&fx=no`,
                              "forecast": `http://api.worldweatheronline.com/free/v2/weather.ashx?key={{.Source.Keys.Key}}&q={{.Location.Latitude}},{{.Location.Longitude}}&format=json&fx=yes&num_of_days=5&tp=24`}
    entry = SourceEntry{Name:"WorldWeatherOnline", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{}
    urls = map[string]string {"current": `http://www.accuweather.com/ru/{{.Location.Iso_country}}/{{.Location.Accuweather_city_name}}/{{.Location.Accuweather_id}}/hourly-weather-forecast/{{.Location.Accuweather_id}}`,
                              "forecast": `http://www.accuweather.com/ru/{{.Location.Iso_country}}/{{.Location.Accuweather_city_name}}/{{.Location.Accuweather_id}}/hourly-weather-forecast/{{.Location.Accuweather_id}}`}
    entry = SourceEntry{Name:"AccuWeather", Urls:urls, Keys:keys}
    sources = append(sources, entry)


    keys = Keyring{}
    urls = map[string]string {"current": `http://beta.gismeteo.ru/weather-{{.Location.Gismeteo_city_name}}-{{.Location.Gismeteo_id}}/`,
                              "forecast": `http://beta.gismeteo.ru/weather-{{.Location.Gismeteo_city_name}}-{{.Location.Gismeteo_id}}/`}
    entry = SourceEntry{Name:"Gismeteo", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    return sources
}
