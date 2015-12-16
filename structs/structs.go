package structs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"github.com/owm-inc/weatherchecker-go/adapters"
	"github.com/owm-inc/weatherchecker-go/db"
	"gopkg.in/mgo.v2/bson"
)

type HistoryDataEntry struct {
	Id           bson.ObjectId             `bson:"_id,omitempty" json:"objectid"`
	Location     LocationEntry             `json:"location"`
	Source       SourceEntry               `json:"source"`
	Measurements adapters.MeasurementArray `json:"measurements"`
	RequestTime  time.Time                 `json:"request_time"`
	WType        string                    `json:"wtype"`
	Url          string                    `json:"url"`
}

func NewHistoryDataEntry(location LocationEntry, source SourceEntry, measurements adapters.MeasurementArray, wtype string, url string) HistoryDataEntry {
	entry := HistoryDataEntry{Location: location, Source: source, Measurements: measurements, WType: wtype, Url: url}
	entry.Id = bson.NewObjectId()

	return entry
}

func GetDataEntry(location LocationEntry, source SourceEntry, wtype string) HistoryDataEntry {
	url := makeUrl(source.Urls[wtype], UrlData{Source: source, Location: location})
	raw := download(url)
	measurements := adapters.AdaptWeather(source.Name, wtype, raw)
	data := NewHistoryDataEntry(location, source, measurements, wtype, url)

	return data
}

type WeatherHistory struct {
	Database   *db.MongoDb
	Collection string
}

func (this *WeatherHistory) AddHistoryEntry(locations []LocationEntry, sources []SourceEntry, wtypes []string) (dataset []HistoryDataEntry) {
	dt := time.Now()
	for _, location := range locations {
		for _, source := range sources {
			for _, wtype := range wtypes {
				data := GetDataEntry(location, source, wtype)
				data.RequestTime = dt

				dataset = append(dataset, data)
			}
		}
	}

	for _, entry := range dataset {
		this.Database.Insert(this.Collection, entry)
	}

	return dataset
}

func (this *WeatherHistory) ShowFullHistory() []HistoryDataEntry {
	var result []HistoryDataEntry
	this.Database.FindAll(this.Collection, &result)
	return result
}

func (this *WeatherHistory) Clear() error {
	err := this.Database.RemoveAll(this.Collection)

	return err
}

func NewWeatherHistory(db_instance *db.MongoDb) WeatherHistory {
	var history = WeatherHistory{Database: db_instance, Collection: "WeatherHistory"}

	return history
}

type Keyring struct {
	Key  string
	Uref string
}

type UrlData struct {
	Source   SourceEntry
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
