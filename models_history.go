package main

import (
	"strconv"

	"github.com/skybon/mgoHelpers"

	"gopkg.in/mgo.v2/bson"
)

type HistoryDataEntryBase struct {
	Status       int64
	Message      string
	Location     LocationEntry
	Source       SourceEntry
	Measurements MeasurementArray
	RequestTime  int64
	WType        string
	Url          string
}

type HistoryDataEntry struct {
	mgoHelpers.DbEntryBase `bson:",inline"`
	HistoryDataEntryBase   `bson:",inline"`
}

func MakeHistoryDataEntry() HistoryDataEntry {
	var entry HistoryDataEntry
	entry.SetBsonID(bson.NewObjectId())

	return entry
}

type WeatherHistory struct {
	Database   *mgoHelpers.MongoDb
	Collection string
}

func (h *WeatherHistory) Add(entry HistoryDataEntry) {
	h.Database.Insert(h.Collection, entry)
}

func (h *WeatherHistory) ReadHistory(entryid string, status int64, source string, wtype string, country string, locationid string, requeststart string, requestend string) (result []HistoryDataEntry) {
	result = []HistoryDataEntry{}
	query := make(map[string]interface{})
	if entryid != "" {
		query["_id"], _ = mgoHelpers.GetObjectIDFromString(entryid)
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
			query["location._id"], _ = mgoHelpers.GetObjectIDFromString(locationid)
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

func NewWeatherHistory(db_instance *mgoHelpers.MongoDb) (history WeatherHistory) {
	history = WeatherHistory{Database: db_instance, Collection: "WeatherHistory"}

	return history
}
