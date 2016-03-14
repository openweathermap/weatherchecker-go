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

type jsonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}

var mongoDsn string
var refreshInterval int
var maxDepth int
var email string

func renderJSON(stuff interface{}, w http.ResponseWriter) {
	data, _ := json.Marshal(stuff)
	jsonString := string(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, jsonString)
}

func renderResponse(status int, message string, content interface{}, w http.ResponseWriter) {
	renderJSON(jsonResponse{Status: status, Message: message, Content: content}, w)
}

func renderLocationEntry(locationEntry LocationEntry, w http.ResponseWriter) {
	renderResponse(200, "OK", map[string][]LocationEntry{"location_entry": []LocationEntry{locationEntry}}, w)
}

func renderStatus(err error, successMessage string, w http.ResponseWriter) {
	var status int
	var message string

	if err != nil {
		status = 500
		message = err.Error()
	} else {
		status = 200
		message = successMessage
	}

	renderResponse(status, message, make(map[string]string), w)
}

func notAllowedForPublic(w http.ResponseWriter, r *http.Request) {
	renderStatus(errors.New("This entrypoint has been disabled for public installation of Weather Checker."), "", w)
}

func invalidAPIKey(w http.ResponseWriter, r *http.Request) {
	renderStatus(errors.New("Invalid API key."), "", w)
}

func validAPIKey(w http.ResponseWriter, r *http.Request) {
	renderStatus(nil, "API key valid.", w)
}

type serverActions struct {
	sources   []SourceEntry
	locations LocationTable
	history   WeatherHistory
}

func (a *serverActions) ReadLocations(w http.ResponseWriter, r *http.Request) {
	renderResponse(200, "OK", map[string][]LocationEntry{"locations": a.locations.ReadLocations()}, w)
}

func (a *serverActions) ReadSources(w http.ResponseWriter, r *http.Request) {
	renderResponse(200, "OK", map[string][]SourceEntry{"sources": a.sources}, w)
}

func (a *serverActions) ReadSanitizedSources(w http.ResponseWriter, r *http.Request) {
	output := make([]map[string]interface{}, len(a.sources))
	for i, source := range a.sources {
		output[i] = source.GetSanitizedInfo()
	}
	renderResponse(200, "OK", map[string]interface{}{"sources": output}, w)
}

func makeMissingParamsList(queryHolder url.Values, requiredParams []string) (missing []string) {
	for _, entry := range requiredParams {
		if queryHolder.Get(entry) == "" {
			missing = append(missing, entry)
		}
	}

	return missing
}

func (a *serverActions) CreateLocation(w http.ResponseWriter, r *http.Request) {
	queryHolder := r.URL.Query()

	missing := makeMissingParamsList(queryHolder, []string{"city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})

	if len(missing) > 0 {
		renderResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry := a.locations.CreateLocation(
			queryHolder.Get("city_name"),
			queryHolder.Get("iso_country"),
			queryHolder.Get("country_name"),
			queryHolder.Get("latitude"),
			queryHolder.Get("longitude"))
		renderLocationEntry(locationEntry, w)
	}
}

func (a *serverActions) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	queryHolder := r.URL.Query()

	missing := makeMissingParamsList(queryHolder, []string{"location_id", "city_name", "iso_country", "iso_country", "country_name", "latitude", "longitude"})
	if len(missing) > 0 {
		renderResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		locationEntry, _ := a.locations.UpdateLocation(
			queryHolder.Get("location_id"),
			queryHolder.Get("city_name"),
			queryHolder.Get("iso_country"),
			queryHolder.Get("country_name"),
			queryHolder.Get("latitude"),
			queryHolder.Get("longitude"))
		renderLocationEntry(locationEntry, w)
	}
}

func (a *serverActions) UpsertLocation(w http.ResponseWriter, r *http.Request) {
	queryHolder := r.URL.Query()

	if len(queryHolder.Get("location_id")) == 0 {
		a.CreateLocation(w, r)
	} else {
		a.UpdateLocation(w, r)
	}
}

func (a *serverActions) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	missing := []string{}

	queryHolder := r.URL.Query()

	makeMissingParamsList(queryHolder, []string{"location_id"})

	if len(missing) > 0 {
		renderResponse(500, "The following parameters are missing: "+strings.Join(missing, ", "), make(map[string]string), w)
	} else {
		err := a.locations.DeleteLocation(queryHolder.Get("location_id"))
		renderStatus(err, "Location removed successfully.", w)
	}
}

func (a *serverActions) ClearLocations(w http.ResponseWriter, r *http.Request) {
	err := a.locations.Clear()

	renderStatus(err, "Locations cleared successfully.", w)
}

func (a *serverActions) ReadHistory(w http.ResponseWriter, r *http.Request, sanitize bool) {
	queryHolder := r.URL.Query()
	entryid := queryHolder.Get("entryid")
	status, _ := strconv.ParseInt(queryHolder.Get("status"), 10, 64)
	source := queryHolder.Get("source")
	wtype := queryHolder.Get("wtype")
	country := queryHolder.Get("country")
	locationid := queryHolder.Get("locationid")
	requestend := queryHolder.Get("requestend")

	var requeststart string
	requeststartRaw := queryHolder.Get("requeststart")
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

	historyData := a.history.ReadHistory(entryid, status, source, wtype, country, locationid, requeststart, requestend)
	historyFilters := map[string]string{"entryid": entryid, "status": strconv.FormatInt(status, 10), "source": source, "wtype": wtype, "country": country, "locationid": locationid, "requeststart": requeststart, "requestend": requestend}

	output := make([]interface{}, len(historyData))
	for i, historyEntry := range historyData {
		entry := make(map[string]interface{})
		entry["objectid"] = historyEntry.Id
		entry["status"] = historyEntry.Status
		entry["location"] = historyEntry.Location
		entry["source"] = historyEntry.Source.GetSanitizedInfo()
		entry["measurements"] = historyEntry.Measurements
		entry["request_time"] = historyEntry.RequestTime
		entry["wtype"] = historyEntry.WType

		if !sanitize {
			entry["url"] = historyEntry.Url
		}
		output[i] = entry
	}

	renderResponse(200, "OK", map[string]interface{}{"history": map[string]interface{}{"data": output, "filters": historyFilters}}, w)
}

func (a *serverActions) ReadFullHistory(w http.ResponseWriter, r *http.Request) {
	a.ReadHistory(w, r, false)
}

func (a *serverActions) ReadSanitizedHistory(w http.ResponseWriter, r *http.Request) {
	a.ReadHistory(w, r, true)
}

func (a *serverActions) RefreshHistory(w http.ResponseWriter, r *http.Request) {
	queryHolder := r.URL.Query()

	var wtypes []string
	wtype := queryHolder.Get("wtype")

	if wtype == "" {
		wtypes = []string{"current", "forecast"}
	} else if wtype != "current" && wtype != "forecast" && wtype != "" {
		renderStatus(errors.New(""), "Invalid wtype specified.", w)
		return
	} else {
		wtypes = []string{wtype}
	}

	historyEntry := PollAll(&a.history, a.locations.ReadLocations(), a.sources, wtypes)
	renderResponse(200, "OK", map[string][]HistoryDataEntry{"historyEntry": historyEntry}, w)
}

func (a *serverActions) ClearHistory(w http.ResponseWriter, r *http.Request) {
	renderStatus(a.history.Clear(), "History cleared successfully.", w)
}

func (a *serverActions) CheckAPIKey(w http.ResponseWriter, r *http.Request, adminKey string, cbSuccess, cbFail func(http.ResponseWriter, *http.Request)) {
	key := r.URL.Query().Get("appid")

	if key == adminKey {
		cbSuccess(w, r)
	} else {
		cbFail(w, r)
	}
}

func (a *serverActions) GetSettings() map[string]interface{} {
	settingsMap := map[string]interface{}{}

	settingsMap["mongo"] = mongoDsn
	settingsMap["refresh-interval"] = refreshInterval
	settingsMap["max-depth"] = maxDepth
	settingsMap["email"] = email

	return settingsMap
}

func (a *serverActions) ReadSettings(w http.ResponseWriter, r *http.Request) {
	renderResponse(200, "OK", map[string]interface{}{"settings": a.GetSettings()}, w)
}

func startActionInstance(dbConn *MongoDb) *serverActions {
	return &serverActions{sources: CreateSources(), locations: NewLocationTable(dbConn), history: NewWeatherHistory(dbConn)}
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

	var dbConn = Db()

	fmt.Println("Connecting to MongoDB at", mongoDsn)
	err := dbConn.Connect(mongoDsn)
	if err != nil {
		fmt.Println(fmt.Sprintf("Database error: %s", err))
		return
	}
	defer dbConn.Disconnect()

	var actions = startActionInstance(dbConn)

	var sMux = http.NewServeMux()

	const APIVer = "0.1"

	const APIEntrypoint = "/" + APIVer

	const KeyCheckEntrypoint = APIEntrypoint + "/check_appid"
	const SettingsEntrypoint = APIEntrypoint + "/settings"
	const SourcesEntrypoint = APIEntrypoint + "/sources"
	const LocationEntrypoint = APIEntrypoint + "/locations"
	const HistoryEntrypoint = APIEntrypoint + "/history"

	sMux.HandleFunc(SettingsEntrypoint, actions.ReadSettings)
	sMux.HandleFunc(LocationEntrypoint, actions.ReadLocations)
	if !closedForPublic {
		sMux.HandleFunc(HistoryEntrypoint, actions.ReadFullHistory)
		sMux.HandleFunc(SourcesEntrypoint, actions.ReadSources)
		sMux.HandleFunc(KeyCheckEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, r.URL.Query().Get("appid"), validAPIKey, invalidAPIKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/refresh", actions.RefreshHistory)
		sMux.HandleFunc(LocationEntrypoint+"/add", actions.CreateLocation)
		sMux.HandleFunc(LocationEntrypoint+"/edit", actions.UpdateLocation)
		sMux.HandleFunc(LocationEntrypoint+"/upsert", actions.UpsertLocation)
		sMux.HandleFunc(LocationEntrypoint+"/remove", actions.DeleteLocation)
		sMux.HandleFunc(LocationEntrypoint+"/clear", actions.ClearLocations)
		sMux.HandleFunc(HistoryEntrypoint+"/clear", actions.ClearHistory)
	} else {
		sMux.HandleFunc(HistoryEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.ReadFullHistory, actions.ReadSanitizedHistory)
		})
		sMux.HandleFunc(SourcesEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.ReadSources, actions.ReadSanitizedSources)
		})
		sMux.HandleFunc(KeyCheckEntrypoint, func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, validAPIKey, invalidAPIKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/refresh", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.RefreshHistory, invalidAPIKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/add", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.CreateLocation, invalidAPIKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/edit", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.UpdateLocation, invalidAPIKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/upsert", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.UpsertLocation, invalidAPIKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/remove", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.DeleteLocation, invalidAPIKey)
		})
		sMux.HandleFunc(LocationEntrypoint+"/clear", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.ClearLocations, invalidAPIKey)
		})
		sMux.HandleFunc(HistoryEntrypoint+"/clear", func(w http.ResponseWriter, r *http.Request) {
			actions.CheckAPIKey(w, r, adminKey, actions.ClearHistory, invalidAPIKey)
		})
	}

	if refreshInterval > 0 {
		go func() {
			for {
				PollAll(&actions.history, actions.locations.ReadLocations(), actions.sources, []string{"current", "forecast"})
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
