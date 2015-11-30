package main

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "text/template"
        "time"

        "github.com/BurntSushi/toml"

        "github.com/skybon/weatherchecker-go/adapters"
        )

type HistoryDataEntry struct {
    Location LocationEntry
    Source SourceEntry
    Measurements adapters.MeasurementArray
    Raw string
}

func (this *WeatherHistory) AddHistoryEntry (proxyTable []WeatherProxy) {
    var dataset HistoryDataArray

    for ip := 0 ; ip < len(proxyTable) ; ip++ {
        var proxy = proxyTable[ip]
        var raw = proxy.Data
        var measurements = adapters.AdaptWeather(proxy.Source.Name, "current", raw)
        var data = HistoryDataEntry {Source:proxy.Source, Location:proxy.Location, Measurements:measurements, Raw:raw}

        dataset = append(dataset, data)
    }

    this.Table = append(this.Table, HistoryEntry {Data:dataset, EntryTime: time.Now(), WType:"current"})
}

type HistoryDataArray []HistoryDataEntry

type HistoryEntry struct {
    EntryTime time.Time
    WType string
    Data HistoryDataArray
}

type HistoryArray []HistoryEntry

type WeatherHistory struct {
    Table HistoryArray
}

type Keyring struct {
    Key string
    Uref string
}

type SourceEntry struct {
    Name string
    Urls map[string]string
    Keys Keyring
}

type LocationEntry struct {
    City_name string `toml:"city_name"`
    Iso_country string `toml:"iso_country"`
    Country_name string `toml:"country_name"`
    Latitude string `toml:"latitude"`
    Longitude string `toml:"longitude"`
    Accuweather_id string `toml:"accuweather_id"`
    Accuweather_city_name string `toml:"accuweather_city_name"`
    Gismeteo_id string `toml:"gismeteo_id"`
    Gismeteo_city_name string `toml:"gismeteo_city_name"`
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

func (this *WeatherProxy) MakeUrl() {
    var urlBuf = new(bytes.Buffer)

    var t = template.New("URL template")
    t, _ = t.Parse(this.Source.Urls["current"])
    t.Execute(urlBuf, this)

    this.Url = urlBuf.String()
}

func (this *WeatherProxy) Refresh() {
    resp, err := http.Get(this.Url)
    if err != nil {
        fmt.Println(`Request finished with error`, err)
    } else {
        defer resp.Body.Close()
        readallContents, _ := ioutil.ReadAll(resp.Body)
        this.Data = string(readallContents)
    }
}

type WeatherProxyTable struct {
    Table []WeatherProxy
}

func (this *WeatherProxyTable) Refresh() {
    for it := 0 ; it < len(this.Table) ; it ++ {
        this.Table[it].Refresh()
    }
}

func (this *WeatherProxyTable) MakeTable(locations []LocationEntry, sources []SourceEntry) {
    for il := 0 ; il < len(locations) ; il++ {
        for is := 0 ; is < len(sources) ; is ++ {
            var proxy = makeProxy(sources[is], locations[il])
            this.Table = append(this.Table, proxy)
        }
    }
}

func makeProxy(source SourceEntry, location LocationEntry) WeatherProxy {
    var proxy = WeatherProxy{Source:source, Location:location}
    proxy.MakeUrl()

    return proxy
}

func loadLocations() []LocationEntry {
    var locationTable LocationTable

    var tomlString = `[[locations]]
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

    toml.Decode(tomlString, &locationTable)
    var locations = locationTable.Locations

    return locations
}

func createSources() []SourceEntry {
    var sources []SourceEntry
    var keys Keyring
    var urls map[string]string
    var entry SourceEntry

    keys = Keyring{Key:string(os.Getenv("OWM_KEY"))}
    urls = map[string]string {"current":`http://api.openweathermap.org/data/2.5/weather?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&units=metric`,
                              "forecast":``}
    entry = SourceEntry{Name:"OpenWeatherMap", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:string(os.Getenv("WUNDERGROUND_KEY"))}
    urls = map[string]string {"current": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/conditions/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`,
                              "forecast": ``}
    entry = SourceEntry{Name:"Weather Underground", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:string(os.Getenv("MYWEATHER2_KEY"))}
    urls = map[string]string {"current": `http://www.myweather2.com/developer/forecast.ashx?uac={{.Source.Keys.Key}}&query={{.Location.Latitude}},{{.Location.Longitude}}&temp_unit=c&output=json&ws_unit=kph`,
                              "forecast": ``}
    entry = SourceEntry{Name:"MyWeather2", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    return sources
}

func main() {
    var locations = loadLocations()
    var sources = createSources()
    var proxyTable WeatherProxyTable
    var history = WeatherHistory{}

    proxyTable.MakeTable(locations, sources)
    proxyTable.Refresh()

    history.AddHistoryEntry(proxyTable.Table)
}
