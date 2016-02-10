package structs

import (
	"os"
)

type SourceEntry struct {
	Name string            `json:"name"`
	Urls map[string]string `json:"urls"`
	Keys Keyring           `json:"keys"`
}

func CreateSources() (sources []SourceEntry) {
	sources = append(sources, SourceEntry{Name: "owm", Urls: map[string]string{"current": `http://api.openweathermap.org/data/2.5/weather?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&units=metric`,
		"forecast": `http://pro.openweathermap.org/data/2.5/forecast?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&cnt=10&mode=json&units=metric`}, Keys: Keyring{Key: os.Getenv("OWM_KEY")}})

	sources = append(sources, SourceEntry{Name: "wunderground", Urls: map[string]string{"current": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/conditions/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`,
		"forecast": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/forecast10day/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`}, Keys: Keyring{Key: os.Getenv("WUNDERGROUND_KEY")}})

	sources = append(sources, SourceEntry{Name: "myweather2", Urls: map[string]string{"current": `http://www.myweather2.com/developer/forecast.ashx?uac={{.Source.Keys.Key}}&query={{.Location.Latitude}},{{.Location.Longitude}}&temp_unit=c&output=json&ws_unit=kph`,
		"forecast": `http://www.myweather2.com/developer/weather.ashx?uac={{.Source.Keys.Key}}&uref={{.Source.Keys.Uref}}&query={{.Location.Latitude}},{{.Location.Longitude}}&output=json&temp_unit=c&ws_unit=kph`}, Keys: Keyring{Key: os.Getenv("MYWEATHER2_KEY"), Uref: os.Getenv("MYWEATHER2_UREF")}})

	sources = append(sources, SourceEntry{Name: "forecast.io", Urls: map[string]string{"current": `https://api.forecast.io/forecast/{{.Source.Keys.Key}}/{{.Location.Latitude}},{{.Location.Longitude}}?units=si`,
		"forecast": `https://api.forecast.io/forecast/{{.Source.Keys.Key}}/{{.Location.Latitude}},{{.Location.Longitude}}?units=si`}, Keys: Keyring{Key: os.Getenv("FORECASTIO_KEY")}})

	sources = append(sources, SourceEntry{Name: "worldweatheronline", Urls: map[string]string{"current": `http://api.worldweatheronline.com/free/v2/weather.ashx?key={{.Source.Keys.Key}}&q={{.Location.Latitude}},{{.Location.Longitude}}&format=json&fx=no&date_format=unix`,
		"forecast": `http://api.worldweatheronline.com/free/v2/weather.ashx?key={{.Source.Keys.Key}}&q={{.Location.Latitude}},{{.Location.Longitude}}&format=json&cc=no&fx=yes&num_of_days=5&tp=1&extra=utcDateTime`}, Keys: Keyring{Key: os.Getenv("WORLDWEATHERONLINE_KEY")}})

	return sources
}
