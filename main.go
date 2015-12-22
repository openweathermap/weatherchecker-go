package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/owm-inc/weatherchecker-go/db"
	"github.com/owm-inc/weatherchecker-go/structs"
)

type JsonResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Content interface {} `json:"content"`
}

var sources = structs.CreateSources()

var mongoDsn string

var db_instance = db.Db()
var locations = structs.NewLocationTable(db_instance)
var history = structs.NewWeatherHistory(db_instance)

func MarshalPrintStuff(stuff interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(stuff)
	jsonString := string(data)
	fmt.Fprintf(w, jsonString)
}

func MarshalPrintResponse(code int, message string, content interface {}, w http.ResponseWriter) {
	MarshalPrintStuff(JsonResponse{Code: code, Message: message, Content: content}, w)
}

func PrintHistory(w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string]interface{}{"history": history.ReadHistory()}, w)
}

func PrintLocationEntry(locationEntry structs.LocationEntry, w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]structs.LocationEntry{"location_entry": []structs.LocationEntry{locationEntry}}, w)
}

func PrintHistoryEntry(historyEntry []structs.HistoryDataEntry, w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]structs.HistoryDataEntry{"history_entry": historyEntry}, w)
}

func PrintLocations(w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]structs.LocationEntry{"locations": locations.ReadLocations()}, w)
}

func PrintStatus(err error, successMessage string, w http.ResponseWriter) {
	var status int
	var message string

	if err != nil {
		status = 500
		message = err.Error()
	} else {
		status = 200
		message = successMessage
	}

	MarshalPrintResponse(status, message, make(map[string]string), w)
}

func MakeMissingParamsList(query_holder url.Values, required_params []string) (missing []string) {
	for _, entry := range required_params {
		if query_holder.Get(entry) == "" {
			missing = append(missing, entry)
		}
	}

	return missing
}

func CreateLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	missing := MakeMissingParamsList(query_holder, []string{"city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})

	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: " + strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry := locations.CreateLocation(query_holder.Get("city_name"),
												  query_holder.Get("iso_country"),
												  query_holder.Get("country_name"),
												  query_holder.Get("latitude"),
												  query_holder.Get("longitude"),
												  query_holder.Get("accuweather_id"),
												  query_holder.Get("accuweather_city_name"),
												  query_holder.Get("gismeteo_id"),
												  query_holder.Get("gismeteo_city_name"))
		PrintLocationEntry(locationEntry, w)
	}
}

func ReadLocations(c web.C, w http.ResponseWriter, r *http.Request) {
	PrintLocations(w)
}

func UpdateLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	missing := MakeMissingParamsList(query_holder, []string{"location_id", "city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})
	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: " + strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry, _ := locations.UpdateLocation(query_holder.Get("location_id"),
													 query_holder.Get("city_name"),
													 query_holder.Get("iso_country"),
													 query_holder.Get("country_name"),
													 query_holder.Get("latitude"),
													 query_holder.Get("longitude"),
													 query_holder.Get("accuweather_id"),
													 query_holder.Get("accuweather_city_name"),
													 query_holder.Get("gismeteo_id"),
													 query_holder.Get("gismeteo_city_name"))
		PrintLocationEntry(locationEntry, w)
	}
}

func DeleteLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	missing := make([]string, 0)

	query_holder := r.URL.Query()

	MakeMissingParamsList(query_holder, []string{"location_id"})

	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: " + strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		err := locations.DeleteLocation(query_holder.Get("location_id"))
		PrintStatus(err, "Location removed successfully.", w)
	}
}

func ClearLocations(c web.C, w http.ResponseWriter, r *http.Request) {
	err := locations.Clear()

	PrintStatus(err, "Locations cleared successfully.", w)
}

func ReadHistory(c web.C, w http.ResponseWriter, r *http.Request) {
	PrintHistory(w)
}

func RefreshHistory(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	var wtypes []string
	wtype := query_holder.Get("wtype")

	if wtype == "" {
		wtypes = []string{"current", "forecast"}
	} else if wtype != "current" && wtype != "forecast" && wtype != "" {
		PrintStatus(errors.New(""), "Invalid wtype specified.", w)
		return
	} else {
		wtypes = []string{wtype}
	}

	locations_query := locations.ReadLocations()
	historyEntry := history.CreateHistoryEntry(locations_query, sources, wtypes)
	PrintHistoryEntry(historyEntry, w)
}

func ClearHistory(c web.C, w http.ResponseWriter, r *http.Request) {
	err := history.Clear()

	PrintStatus(err, "History cleared successfully.", w)
}

func Api(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Pass data through the environment
		c.Env["history"] = &history
		// Fully control how the next layer is called
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func init() {
	flag.StringVar(&mongoDsn, "mongo", "mongodb://127.0.0.1:27017/weatherchecker", "MongoDB DSN")
}

func main() {
	flag.Parse()

	if os.Getenv("MONGO") != "" {
		mongoDsn = os.Getenv("MONGO")
	}

	fmt.Println("Connecting to MongoDB at", mongoDsn)
	err := db_instance.Connect(mongoDsn)
	if err != nil {
		fmt.Println(fmt.Sprintf("Database error: %s", err))
		return
	}
	defer db_instance.Disconnect()

	const ApiVer = "0.1"

	const ApiEntrypoint = "/api" + "/" + ApiVer

	const LocationEntrypoint = ApiEntrypoint + "/locations"
	const HistoryEntrypoint = ApiEntrypoint + "/history"

	goji.Use(Api)
	goji.Get(LocationEntrypoint, ReadLocations)
	goji.Get(LocationEntrypoint+"/add", CreateLocation)
	goji.Get(LocationEntrypoint+"/edit", UpdateLocation)
	goji.Get(LocationEntrypoint+"/remove", DeleteLocation)
	goji.Get(LocationEntrypoint+"/clear", ClearLocations)
	goji.Get(HistoryEntrypoint, ReadHistory)
	goji.Get(HistoryEntrypoint+"/refresh", RefreshHistory)
	goji.Get(HistoryEntrypoint+"/clear", ClearHistory)
	goji.Serve()
}
