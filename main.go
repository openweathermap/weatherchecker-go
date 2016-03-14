package main

//go:generate browserify ./ui/scripts/main_index.js -o ./ui/bundle/index.js
//go:generate browserify ./ui/scripts/main_admin.js -o ./ui/bundle/admin.js
//go:generate browserify ./ui/scripts/main_analytics.js -o ./ui/bundle/analytics.js
//go:generate browserify ./ui/scripts/ga.js -o ./ui/bundle/ga.js
import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type JsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

var mongoDsn string
var refreshInterval int
var maxDepth int
var email string

var db_instance *MongoDb
var sources []SourceEntry
var locations LocationTable
var history WeatherHistory

func MarshalPrintStuff(stuff interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(stuff)
	jsonString := string(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, jsonString)
}

func MarshalPrintResponse(status int, message string, content interface{}, w http.ResponseWriter) {
	MarshalPrintStuff(JsonResponse{Status: status, Message: message, Content: content}, w)
}

func PrintLocationEntry(locationEntry LocationEntry, w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]LocationEntry{"location_entry": []LocationEntry{locationEntry}}, w)
}

func PrintHistoryEntry(historyEntry []HistoryDataEntry, w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]HistoryDataEntry{"history_entry": historyEntry}, w)
}

func PrintLocations(w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]LocationEntry{"locations": locations.ReadLocations()}, w)
}

func PrintSources(w http.ResponseWriter) {
	MarshalPrintResponse(200, "OK", map[string][]SourceEntry{"sources": sources}, w)
}

func PrintSanitizedSources(w http.ResponseWriter) {
	output := make([]map[string]interface{}, len(sources))
	for i, source := range sources {
		entry := SanitizeSource(source)
		output[i] = entry
	}
	MarshalPrintResponse(200, "OK", map[string]interface{}{"sources": output}, w)
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

func NotAllowedForPublic(w http.ResponseWriter, r *http.Request) {
	PrintStatus(errors.New("This entrypoint has been disabled for public installation of Weather Checker."), "", w)
}

func InvalidApiKey(w http.ResponseWriter, r *http.Request) {
	PrintStatus(errors.New("Invalid API key."), "", w)
}

func ValidApiKey(w http.ResponseWriter, r *http.Request) {
	PrintStatus(nil, "API key valid.", w)
}

func SanitizeSource(source SourceEntry) map[string]interface{} {
	entry := make(map[string]interface{})
	entry["name"] = source.Name
	entry["prettyname"] = source.PrettyName
	entry["urls"] = source.Urls

	return entry
}

func MakeMissingParamsList(query_holder url.Values, required_params []string) (missing []string) {
	for _, entry := range required_params {
		if query_holder.Get(entry) == "" {
			missing = append(missing, entry)
		}
	}

	return missing
}

func ReadSources(w http.ResponseWriter, r *http.Request) {
	PrintSources(w)
}

func ReadSanitizedSources(w http.ResponseWriter, r *http.Request) {
	PrintSanitizedSources(w)
}

func CreateLocation(w http.ResponseWriter, r *http.Request) {
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
			query_holder.Get("longitude"))
		PrintLocationEntry(locationEntry, w)
	}
}

func ReadLocations(w http.ResponseWriter, r *http.Request) {
	PrintLocations(w)
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
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
			query_holder.Get("longitude"))
		PrintLocationEntry(locationEntry, w)
	}
}

func UpsertLocation(w http.ResponseWriter, r *http.Request) {
	query_holder := r.URL.Query()

	if len(query_holder.Get("location_id")) == 0 {
		CreateLocation(w, r)
	} else {
		UpdateLocation(w, r)
	}
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	missing := []string{}

	query_holder := r.URL.Query()

	MakeMissingParamsList(query_holder, []string{"location_id"})

	if len(missing) > 0 {
		MarshalPrintResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		err := locations.DeleteLocation(query_holder.Get("location_id"))
		PrintStatus(err, "Location removed successfully.", w)
	}
}

func ClearLocations(w http.ResponseWriter, r *http.Request) {
	err := locations.Clear()

	PrintStatus(err, "Locations cleared successfully.", w)
}

func ReadHistory(w http.ResponseWriter, r *http.Request, sanitize bool) {
	query_holder := r.URL.Query()
	entryid := query_holder.Get("entryid")
	status, _ := strconv.ParseInt(query_holder.Get("status"), 10, 64)
	source := query_holder.Get("source")
	wtype := query_holder.Get("wtype")
	country := query_holder.Get("country")
	locationid := query_holder.Get("locationid")
	requestend := query_holder.Get("requestend")

	var requeststart string
	requeststartRaw := query_holder.Get("requeststart")
	if sanitize && maxDepth > 0 {
		requeststartRawInt, _ := strconv.ParseInt(requeststartRaw, 10, 64)

		currentTime := time.Now().Unix()

		a := requeststartRawInt
		b := currentTime - int64(3600*maxDepth)
		requestStartInt := map[bool]int64{true: a, false: b}[a > b]
		requeststart = strconv.FormatInt(requestStartInt, 10)
	} else {
		requeststart = requeststartRaw
	}

	history_data := history.ReadHistory(entryid, status, source, wtype, country, locationid, requeststart, requestend)
	history_filters := map[string]string{"entryid": entryid, "status": strconv.FormatInt(status, 10), "source": source, "wtype": wtype, "country": country, "locationid": locationid, "requeststart": requeststart, "requestend": requestend}

	output := make([]interface{}, len(history_data))
	for i, history_entry := range history_data {
		entry := make(map[string]interface{})
		entry["objectid"] = history_entry.Id
		entry["status"] = history_entry.Status
		entry["location"] = history_entry.Location
		entry["source"] = SanitizeSource(history_entry.Source)
		entry["measurements"] = history_entry.Measurements
		entry["request_time"] = history_entry.RequestTime
		entry["wtype"] = history_entry.WType

		if !sanitize {
			entry["url"] = history_entry.Url
		}
		output[i] = entry
	}

	MarshalPrintResponse(200, "OK", map[string]interface{}{"history": map[string]interface{}{"data": output, "filters": history_filters}}, w)
}

func ReadFullHistory(w http.ResponseWriter, r *http.Request) {
	ReadHistory(w, r, false)
}

func ReadSanitizedHistory(w http.ResponseWriter, r *http.Request) {
	ReadHistory(w, r, true)
}

func RefreshHistory(w http.ResponseWriter, r *http.Request) {
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

	historyEntry := RefreshHistoryCore(sources, wtypes)
	PrintHistoryEntry(historyEntry, w)
}

func RefreshHistoryCore(sources []SourceEntry, wtypes []string) []HistoryDataEntry {
	locations_query := locations.ReadLocations()
	historyEntry := PollAll(&history, locations_query, sources, wtypes)

	return historyEntry
}

func ClearHistory(w http.ResponseWriter, r *http.Request) {
	err := history.Clear()

	PrintStatus(err, "History cleared successfully.", w)
}

func CheckApiKey(w http.ResponseWriter, r *http.Request, adminKey string, cbSuccess, cbFail func(http.ResponseWriter, *http.Request)) {
	key := r.URL.Query().Get("appid")

	if key == adminKey {
		cbSuccess(w, r)
	} else {
		cbFail(w, r)
	}
}

func GetSettings() map[string]interface{} {
	settingsMap := map[string]interface{}{}

	settingsMap["mongo"] = mongoDsn
	settingsMap["refresh-interval"] = refreshInterval
	settingsMap["max-depth"] = maxDepth
	settingsMap["email"] = email

	return settingsMap
}

func ReadSettings(w http.ResponseWriter, r *http.Request) {
	MarshalPrintResponse(200, "OK", map[string]interface{}{"settings": GetSettings()}, w)
}

func init() {
	flag.StringVar(&mongoDsn, "mongo", "mongodb://127.0.0.1:27017/weatherchecker", "MongoDB DSN")
	flag.IntVar(&refreshInterval, "refresh-interval", 0, "Refresh interval")
	flag.IntVar(&maxDepth, "max-depth", 0, "Maximum depth (h) for unpriveleged requests")
}

func main() {
	flag.Parse()

	if os.Getenv("MONGO") != "" {
		mongoDsn = os.Getenv("MONGO")
	}

	email = os.Getenv("WC_EMAIL")

	var adminKey = os.Getenv("ADMIN_PASS")

	var closedForPublic bool
	if os.Getenv("CLOSED_FOR_PUBLIC") == "1" {
		closedForPublic = true
	}

	sources = CreateSources()
	db_instance = Db()
	locations = NewLocationTable(db_instance)
	history = NewWeatherHistory(db_instance)

	fmt.Println("Connecting to MongoDB at", mongoDsn)
	err := db_instance.Connect(mongoDsn)
	if err != nil {
		fmt.Println(fmt.Sprintf("Database error: %s", err))
		return
	}
	defer db_instance.Disconnect()

	var sMux = http.NewServeMux()

	const APIVer = "0.1"

	const APIEntrypoint = "/" + APIVer

	const KeyCheckEntrypoint = APIEntrypoint + "/check_appid"
	const SettingsEntrypoint = APIEntrypoint + "/settings"
	const SourcesEntrypoint = APIEntrypoint + "/sources"
	const LocationEntrypoint = APIEntrypoint + "/locations"
	const HistoryEntrypoint = APIEntrypoint + "/history"

	sMux.HandleFunc(SettingsEntrypoint, ReadSettings)
	sMux.HandleFunc(LocationEntrypoint, ReadLocations)
	if !closedForPublic {
		sMux.HandleFunc(HistoryEntrypoint, ReadFullHistory)
		sMux.HandleFunc(SourcesEntrypoint, ReadSources)
		sMux.HandleFunc(KeyCheckEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, r.URL.Query().Get("appid"), ValidApiKey, InvalidApiKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/refresh", RefreshHistory)
		sMux.HandleFunc(LocationEntrypoint+"/add", CreateLocation)
		sMux.HandleFunc(LocationEntrypoint+"/edit", UpdateLocation)
		sMux.HandleFunc(LocationEntrypoint+"/upsert", UpsertLocation)
		sMux.HandleFunc(LocationEntrypoint+"/remove", DeleteLocation)
		sMux.HandleFunc(LocationEntrypoint+"/clear", ClearLocations)
		sMux.HandleFunc(HistoryEntrypoint+"/clear", ClearHistory)
	} else {
		sMux.HandleFunc(HistoryEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, ReadFullHistory, ReadSanitizedHistory)
		})
		sMux.HandleFunc(SourcesEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, ReadSources, ReadSanitizedSources)
		})
		sMux.HandleFunc(KeyCheckEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, ValidApiKey, InvalidApiKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/refresh", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, RefreshHistory, InvalidApiKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/add", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, CreateLocation, InvalidApiKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/edit", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, UpdateLocation, InvalidApiKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/upsert", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, UpsertLocation, InvalidApiKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/remove", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, DeleteLocation, InvalidApiKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/clear", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, ClearLocations, InvalidApiKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/clear", func(w http.ResponseWriter, r *http.Request) {
			CheckApiKey(w, r, adminKey, ClearHistory, InvalidApiKey)
		})
	}

	if refreshInterval > 0 {
		go func() {
			for {
				RefreshHistoryCore(sources, []string{"current", "forecast"})
				time.Sleep(time.Duration(refreshInterval) * time.Minute)
			}
		}()
	}

	var server = &http.Server{
		Addr:    ":8000",
		Handler: sMux,
	}
	server.ListenAndServe()
}
