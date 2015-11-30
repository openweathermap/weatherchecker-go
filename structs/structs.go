package structs

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "time"
        "text/template"

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
    Name string `json:"name"`
    Urls map[string]string `json:"urls"`
    Keys Keyring
}

func CreateSources() []SourceEntry {
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

type LocationEntry struct {
    City_name string `json:"city_name"`
    Iso_country string `json:"iso_country"`
    Country_name string `json:"country_name"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Accuweather_id string `json:"accuweather_id"`
    Accuweather_city_name string `json:"accuweather_city_name"`
    Gismeteo_id string `json:"gismeteo_id"`
    Gismeteo_city_name string `json:"gismeteo_city_name"`
}

type LocationTable struct {
    Locations []LocationEntry `json:"locations"`
}

func LoadLocations() []LocationEntry {
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

type WeatherProxy struct {
    Source SourceEntry `json:"source"`
    Location LocationEntry `json:"location"`
    Data string `json:"data"`
}

func (this *WeatherProxy) MakeUrl() string {
    var urlBuf = new(bytes.Buffer)

    var t = template.New("URL template")
    t, _ = t.Parse(this.Source.Urls["current"])
    t.Execute(urlBuf, this)

    urlString := urlBuf.String()
    return urlString
}

func NewProxy(source SourceEntry, location LocationEntry) WeatherProxy {
    proxy := WeatherProxy{Source:source, Location:location}

    return proxy
}

func (this *WeatherProxy) Refresh() {
    url := this.MakeUrl()
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(`Request finished with error`, err)
    } else {
        defer resp.Body.Close()
        readallContents, _ := ioutil.ReadAll(resp.Body)
        this.Data = string(readallContents)
    }
}

type WeatherProxyTable struct {
    Table []WeatherProxy `json:"proxies"`
}


func (this *WeatherProxyTable) Refresh() {
    for it := 0 ; it < len(this.Table) ; it ++ {
        this.Table[it].Refresh()
    }
}

func NewWeatherProxyTable(locations []LocationEntry, sources []SourceEntry) WeatherProxyTable {
    var newProxyTable WeatherProxyTable
    for il := 0 ; il < len(locations) ; il++ {
        for is := 0 ; is < len(sources) ; is ++ {
            var proxy = NewProxy(sources[is], locations[il])
            newProxyTable.Table = append(newProxyTable.Table, proxy)
        }
    }

    return newProxyTable
}
