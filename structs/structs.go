package structs

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "time"
        "text/template"

        "gopkg.in/mgo.v2/bson"
        "github.com/owm-inc/weatherchecker-go/db"
        "github.com/owm-inc/weatherchecker-go/adapters"
        )

type HistoryDataEntry struct {
    Id bson.ObjectId `bson:"_id,omitempty" json:"objectid"`
    Location LocationEntry
    Source SourceEntry
    Measurements adapters.MeasurementArray
    WType string
    Url string
}

func NewHistoryDataEntry (location LocationEntry, source SourceEntry, measurements adapters.MeasurementArray, wtype string, url string) HistoryDataEntry {
    entry := HistoryDataEntry {Location:location, Source:source, Measurements:measurements, WType:wtype, Url:url}
    entry.Id = bson.NewObjectId()

    return entry
}

func GetDataEntry (location LocationEntry, source SourceEntry, wtype string) HistoryDataEntry {
    url := makeUrl (source.Urls[wtype], UrlData {Source:source, Location:location})
    raw := download (url)
    measurements := adapters.AdaptWeather(source.Name, wtype, raw)
    data := NewHistoryDataEntry(location, source, measurements, wtype, url)

    return data
}

type HistoryEntry struct {
    Id bson.ObjectId `bson:"_id,omitempty"`
    EntryTime time.Time
    WType string
    Data []HistoryDataEntry
}

func NewHistoryEntry (dataset []HistoryDataEntry, entryTime time.Time, wType string) HistoryEntry {
    var historyEntry = HistoryEntry {Data:dataset, EntryTime: entryTime, WType:wType}

    return historyEntry
}

type WeatherHistory struct {
    Database *db.MongoDb
    Collection string
}

func (this *WeatherHistory) AddHistoryEntry (locations []LocationEntry, sources []SourceEntry, wtypes []string) HistoryEntry {
    var dataset []HistoryDataEntry

    for _, location := range locations {
        for _, source := range sources {
            for _, wtype := range wtypes {
                data := GetDataEntry(location, source, wtype)

                dataset = append(dataset, data)
            }
        }
    }

    newHistoryEntry := NewHistoryEntry(dataset, time.Now(), "current")
    this.Database.Insert(this.Collection, newHistoryEntry)

    return newHistoryEntry
}

func (this *WeatherHistory) ShowFullHistory () []HistoryEntry {
    var result []HistoryEntry
    this.Database.FindAll(this.Collection, &result)
    return result
}

func NewWeatherHistory (db_instance *db.MongoDb) WeatherHistory {
    var history = WeatherHistory {Database:db_instance, Collection:"WeatherHistory"}

    return history
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

    keys = Keyring{Key:os.Getenv("OWM_KEY")}
    urls = map[string]string {"current":`http://api.openweathermap.org/data/2.5/weather?appid={{.Source.Keys.Key}}&lat={{.Location.Latitude}}&lon={{.Location.Longitude}}&units=metric`,
                              "forecast":``}
    entry = SourceEntry{Name:"OpenWeatherMap", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("WUNDERGROUND_KEY")}
    urls = map[string]string {"current": `http://api.wunderground.com/api/{{.Source.Keys.Key}}/conditions/q/{{.Location.Latitude}},{{.Location.Longitude}}.json`,
                              "forecast": ``}
    entry = SourceEntry{Name:"Weather Underground", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{Key:os.Getenv("MYWEATHER2_KEY"), Uref:os.Getenv("MYWEATHER2_UREF")}
    urls = map[string]string {"current": `http://www.myweather2.com/developer/forecast.ashx?uac={{.Source.Keys.Key}}&query={{.Location.Latitude}},{{.Location.Longitude}}&temp_unit=c&output=json&ws_unit=kph`,
                              "forecast": ``}
    entry = SourceEntry{Name:"MyWeather2", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    keys = Keyring{}
    urls = map[string]string {"current": `http://www.accuweather.com/ru/{{.Location.Iso_country}}/{{.Location.Accuweather_city_name}}/{{.Location.Accuweather_id}}/hourly-weather-forecast/{{.Location.Accuweather_id}}`,
                              "forecast": ``}
    entry = SourceEntry{Name:"AccuWeather", Urls:urls, Keys:keys}
    sources = append(sources, entry)

    return sources
}


type LocationEntry struct {
    Id bson.ObjectId `bson:"_id,omitempty" json:"objectid"`
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

func NewLocationEntry (city_name string, iso_country string, country_name string, latitude string, longitude string, accuweather_id string, accuweather_city_name string, gismeteo_id string, gismeteo_city_name string) LocationEntry {
    model := LocationEntry {City_name:city_name, Iso_country:iso_country, Country_name:country_name, Latitude:latitude, Longitude:longitude, Accuweather_id:accuweather_id, Accuweather_city_name:accuweather_city_name,Gismeteo_id:gismeteo_id, Gismeteo_city_name:gismeteo_city_name}
    model.Id = bson.NewObjectId()

    return model
}

type LocationTable struct {
    Database *db.MongoDb
    Collection string
}

func (this *LocationTable) AddLocation (city_name string, iso_country string, country_name string, latitude string, longitude string, accuweather_id string, accuweather_city_name string, gismeteo_id string, gismeteo_city_name string) LocationEntry {
    newLocationEntry := NewLocationEntry (city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name)
    this.Database.Insert(this.Collection, newLocationEntry)

    return newLocationEntry
}

func (this *LocationTable) RetrieveLocations () []LocationEntry {
    var result []LocationEntry
    this.Database.FindAll(this.Collection, &result)
    return result
}

func (this *LocationTable) RemoveLocation (location_id bson.ObjectId) error {
    err := this.Database.Remove(this.Collection, location_id)

    return err
}

func NewLocationTable (db_instance *db.MongoDb) LocationTable {
    var locations = LocationTable {Database:db_instance, Collection:"Locations"}

    return locations
}

type UrlData struct {
    Source SourceEntry
    Location LocationEntry
}

func makeUrl(url_template string, data UrlData) string {
    var urlBuf = new(bytes.Buffer)

    var t = template.New("URL template")
    t, _ = t.Parse(url_template)
    t.Execute(urlBuf, data)

    urlString := urlBuf.String()
    return urlString
}

func download(url string) string {
    resp, err := http.Get(url)
    var data string
    if err != nil {
        fmt.Println(`Request finished with error`, err)
        data = ``
    } else {
        defer resp.Body.Close()
        readallContents, _ := ioutil.ReadAll(resp.Body)
        data = string(readallContents)
    }
    return data
}
