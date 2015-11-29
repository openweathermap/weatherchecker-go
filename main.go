package main

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "text/template"

        "github.com/BurntSushi/toml"

        "github.com/skybon/weatherchecker-go/adapters"
        )

type HistoryDataEntry struct {
    Location LocationEntry
    Source SourceEntry
    Measurements adapters.MeasurementSchema
}

type HistoryEntry struct {
    Time string
    WType string
    Data []HistoryDataEntry
}

type WeatherHistory struct {
    Table []HistoryEntry
}

type Keyring struct {
    Key string
    Uref string
}

type SourceEntry struct {
    Name string
    Url string
    Keys Keyring
}

type LocationEntry struct {
    City_name string `toml:"city_name"`
    Iso_country string
    Country_name string
    Latitude string
    Longitude string
    Accuweather_id string
    Accuweather_city_name string
    Gismeteo_id string
    Gismeteo_city_name string
}

type LocationTable struct {
    Locations []LocationEntry `toml:"locations"`
}

type WeatherProxy struct {
    Url string
    Source SourceEntry
    Location LocationEntry
    Data string
}

func (p *WeatherProxy) MakeUrl() {
    var url_buf = new(bytes.Buffer)

    var t = template.New("URL template")
    t, _ = t.Parse(p.Source.Url)
    t.Execute(url_buf, p)

    p.Url = url_buf.String()
}

func (p *WeatherProxy) Refresh() {
    resp, err := http.Get(p.Url)
    if err != nil {
        fmt.Println(`Request finished with error`, err)
    } else {
        defer resp.Body.Close()
        readall_contents, _ := ioutil.ReadAll(resp.Body)
        p.Data = string(readall_contents)
    }
}

type WeatherProxyTable struct {
    Table []WeatherProxy
}

func (t *WeatherProxyTable) Refresh() {
    for it := 0 ; it < len(t.Table) ; it ++ {
        t.Table[it].Refresh()
    }
}

func make_proxy(source SourceEntry, location LocationEntry) WeatherProxy {
    var proxy = WeatherProxy{Source:source, Location:location}
    proxy.MakeUrl()

    return proxy
}

func load_locations() []LocationEntry {
    var location_table LocationTable

    var toml_string = `[[locations]]
                       city_name = "Москва"
                       iso_country = "RU"
                       country_name = "Россия"
                       latitude = "55.75"
                       longitude = "37.62"
                       accuweather_id = "294021"
                       accuweather_city_name = "moscow"
                       gismeteo_id = "4368"
                       gismeteo_city_name = "moscow"

                       [[locations]]
                       city_name = "Санкт-Петербург"
                       iso_country = "RU"
                       country_name = "Россия"
                       latitude = "59.95"
                       longitude = "30.3"
                       accuweather_id = "295212"
                       accuweather_city_name = "saint-petersburg"
                       gismeteo_id = "4079"
                       gismeteo_city_name = "sankt-peterburg"`

    toml.Decode(toml_string, &location_table)
    var locations = location_table.Locations

    return locations
}

func create_sources() []SourceEntry {
    var sources []SourceEntry
    var keys Keyring
    var entry SourceEntry

    keys = Keyring{Key:string(os.Getenv("OWM_KEY"))}
    entry = SourceEntry{Name:"OpenWeatherMap", Url:`http://api.openweathermap.org/data/2.5/weather?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&units=metric`, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:string(os.Getenv("WUNDERGROUND_KEY"))}
    entry = SourceEntry{Name:"Weather Underground", Url:`http://api.wunderground.com/api/{{.Source.Keys.Key}}/conditions/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`, Keys:keys}
    sources = append(sources, entry)

    return sources
}

func main() {
    var locations = load_locations()
    var sources = create_sources()
    var proxy_table WeatherProxyTable
    var history = WeatherHistory{}

    for il := 0 ; il < len(locations) ; il++ {
        for is := 0 ; is < len(sources) ; is ++ {
            var proxy = make_proxy(sources[is], locations[il])
            proxy_table.Table = append(proxy_table.Table, proxy)
        }
    }

    proxy_table.Refresh()

    var dataset []HistoryDataEntry

    for ip := 0 ; ip < len(proxy_table.Table) ; ip++ {
        var proxy = proxy_table.Table[ip]
        var measurement = adapters.Owm_adapt_weather(proxy.Data)
        var data = HistoryDataEntry {Source:proxy.Source, Location:proxy.Location, Measurements:measurement}

        dataset = append(dataset, data)
        fmt.Println(proxy.Source.Name)
        fmt.Println(proxy.Location.City_name)
        fmt.Println(dataset[ip].Measurements)
        fmt.Println("---------------------")
    }

    var history_entry = HistoryEntry {Data:dataset, Time:"0", WType:"current"}

    history.Table = append(history.Table, history_entry)
}
