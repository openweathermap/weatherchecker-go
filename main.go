package main

//go:generate go-bindata -o "bindata/bindata.go" -pkg "bindata" "data/..."
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

	"github.com/owm-inc/weatherchecker-go/bindata"
	"github.com/owm-inc/weatherchecker-go/db"
	"github.com/owm-inc/weatherchecker-go/structs"
)

type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
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

func MarshalPrintResponse(status int, message string, content interface{}, w http.ResponseWriter) {
	MarshalPrintStuff(JsonResponse{Status: status, Message: message, Content: content}, w)
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

func PrintSources(w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]structs.SourceEntry{"sources": sources}, w)
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

func ReadSources(c web.C, w http.ResponseWriter, r *http.Request) {
	PrintSources(w)
}

func CreateLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	missing := MakeMissingParamsList(query_holder, []string{"city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})

	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry := locations.CreateLocation(
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

func ReadLocations(c web.C, w http.ResponseWriter, r *http.Request) {
	PrintLocations(w)
}

func UpdateLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	missing := MakeMissingParamsList(query_holder, []string{"location_id", "city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})
	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry, _ := locations.UpdateLocation(
			query_holder.Get("location_id"),
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

func UpsertLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	if len(query_holder.Get("location_id")) == 0 {
		CreateLocation(c, w, r)
	} else {
		UpdateLocation(c, w, r)
	}
}

func DeleteLocation(c web.C, w http.ResponseWriter, r *http.Request) {
	missing := make([]string, 0)

	query_holder := r.URL.Query()

	MakeMissingParamsList(query_holder, []string{"location_id"})

	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
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
	query_holder := r.URL.Query()
	entryid := query_holder.Get("entryid")
	source := query_holder.Get("source")
	wtype := query_holder.Get("wtype")
	country := query_holder.Get("country")
	locationid := query_holder.Get("locationid")
	requeststart := query_holder.Get("requeststart")
	requestend := query_holder.Get("requestend")

	history_data := history.ReadHistory(entryid, source, wtype, country, locationid, requeststart, requestend)
	history_filters := map[string]string{"entryid": entryid, "source": source, "wtype": wtype, "country": country, "locationid": locationid, "requeststart": requeststart, "requestend": requestend}

	MarshalPrintResponse(200, "OK", map[string]interface{}{"history": map[string]interface{}{"data": history_data, "filters": history_filters}}, w)
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetPath(c web.C, w http.ResponseWriter, r *http.Request) {
	assetPath := "data" + r.URL.Path
	asset, err := bindata.Asset(assetPath)
	if err == nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, string(asset))
	} else {
		fmt.Fprintf(w, err.Error()+"\n")
	}
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

	const SourcesEntrypoint = ApiEntrypoint + "/sources"
	const LocationEntrypoint = ApiEntrypoint + "/locations"
	const HistoryEntrypoint = ApiEntrypoint + "/history"

	const UIEntrypoint = "/ui"

	goji.Use(Api)
	goji.Get(SourcesEntrypoint, ReadSources)
	goji.Get(LocationEntrypoint, ReadLocations)
	goji.Get(LocationEntrypoint+"/add", CreateLocation)
	goji.Get(LocationEntrypoint+"/edit", UpdateLocation)
	goji.Get(LocationEntrypoint+"/upsert", UpsertLocation)
	goji.Get(LocationEntrypoint+"/remove", DeleteLocation)
	goji.Get(LocationEntrypoint+"/clear", ClearLocations)
	goji.Get(HistoryEntrypoint, ReadHistory)
	goji.Get(HistoryEntrypoint+"/refresh", RefreshHistory)
	goji.Get(HistoryEntrypoint+"/clear", ClearHistory)

	goji.Get(UIEntrypoint+"/*", GetPath)
	goji.Get(UIEntrypoint, http.RedirectHandler(UIEntrypoint+"/index.html", 301))
	goji.Serve()
}
