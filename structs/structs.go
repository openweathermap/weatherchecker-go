package structs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/owm-inc/weatherchecker-go/adapters"
	"github.com/owm-inc/weatherchecker-go/db"
	"gopkg.in/mgo.v2/bson"
)

type DbEntryBase struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"objectid"`
}

type HistoryDataEntryBase struct {
	Status       int64
	Message      string
	Location     LocationEntry
	Source       SourceEntry
	Measurements adapters.MeasurementArray
	RequestTime  int64
	WType        string
	Url          string
}

type HistoryDataEntry struct {
	DbEntryBase          `bson:",inline"`
	HistoryDataEntryBase `bson:",inline"`
}

func NewHistoryDataEntry(location LocationEntry, source SourceEntry, measurements adapters.MeasurementArray, wtype string, url string, err error) (entry HistoryDataEntry) {
	var status int64
	var message string
	if err != nil {
		status = 500
		message = err.Error()
	} else {
		status = 200
		message = "OK"
	}
	entry = HistoryDataEntry{DbEntryBase{Id: bson.NewObjectId()}, HistoryDataEntryBase{Status: status, Message: message, Location: location, Source: source, Measurements: measurements, WType: wtype, Url: url}}

	return entry
}

func MakeDataEntry(location LocationEntry, source SourceEntry, wtype string) (entry HistoryDataEntry) {
	var err error
	var url string
	var raw string
	measurements := make(adapters.MeasurementArray, 0)

	adaptFunc, adaptFuncLookupErr := adapters.GetAdaptFunc(source.Name, wtype)

	if adaptFuncLookupErr == nil {
		url = makeUrl(source.Urls[wtype], UrlData{Source: source, Location: location})

		var downloadErr error
		raw, downloadErr = download(url)

		if downloadErr != nil {
			measurements = adapters.AdaptStub(raw)
			err = downloadErr
		} else {
			var adaptErr error
			measurements, adaptErr = adaptFunc(raw)

			err = adaptErr
		}

	} else {
		err = adaptFuncLookupErr
	}

	entry = NewHistoryDataEntry(location, source, measurements, wtype, url, err)

	return entry
}

type WeatherHistory struct {
	Database   *db.MongoDb
	Collection string
}

func (h *WeatherHistory) CreateHistoryEntry(locations []LocationEntry, sources []SourceEntry, wtypes []string) (dataset []HistoryDataEntry) {
	dt := time.Now().Unix()

	dataChan := make(chan HistoryDataEntry, 9999)
	doneChan := make(chan struct{})

	go func() {
		for entry := range dataChan {
			h.Database.Insert(h.Collection, entry)
			dataset = append(dataset, entry)
		}
		doneChan <- struct{}{}
	}()

	for _, location := range locations {
		for _, source := range sources {
			for _, wtype := range wtypes {
				data := MakeDataEntry(location, source, wtype)
				data.RequestTime = dt

				dataChan <- data
			}
		}
	}
	close(dataChan)
	<-doneChan

	return dataset
}

func (h *WeatherHistory) ReadHistory(entryid string, status int64, source string, wtype string, country string, locationid string, requeststart string, requestend string) (result []HistoryDataEntry) {
	result = []HistoryDataEntry{}
	query := make(map[string]interface{})
	if entryid != "" {
		query["_id"], _ = db.GetObjectIDFromString(entryid)
	} else {
		if status != 0 {
			query["status"] = status
		}
		if source != "" {
			query["source.name"] = source
		}
		if wtype != "" {
			query["wtype"] = wtype
		}
		if country != "" {
			query["location.iso_country"] = country
		}
		if locationid != "" {
			query["location._id"], _ = db.GetObjectIDFromString(locationid)
		}
		if requeststart != "" || requestend != "" {
			requestquery := make(map[string]int64)
			if requeststart != "" {
				requestquery[`$gte`], _ = strconv.ParseInt(requeststart, 10, 64)
			}
			if requestend != "" {
				requestquery[`$lte`], _ = strconv.ParseInt(requestend, 10, 64)
			}
			query["requesttime"] = requestquery
		}
	}

	h.Database.Find(h.Collection, query, &result)
	return result
}

func (this *WeatherHistory) Clear() (err error) {
	err = this.Database.RemoveAll(this.Collection)

	return err
}

func NewWeatherHistory(db_instance *db.MongoDb) (history WeatherHistory) {
	history = WeatherHistory{Database: db_instance, Collection: "WeatherHistory"}

	return history
}

type Keyring struct {
	Key  string `json:"key"`
	Uref string `json:"uref"`
}

type UrlData struct {
	Source   SourceEntry
	Location LocationEntry
}

func makeUrl(url_template string, data UrlData) (urlString string) {
	var urlBuf = new(bytes.Buffer)

	var t = template.New("URL template")
	t, _ = t.Parse(url_template)
	t.Execute(urlBuf, data)

	urlString = urlBuf.String()
	return urlString
}

func download(url string) (data string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		readallContents, _ := ioutil.ReadAll(resp.Body)
		data = string(readallContents)
	}
	return data, err
}
