package main

import (
        "encoding/json"
        "errors"
        "flag"
        "fmt"
        "net/http"
        "os"
        "strings"

        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"

        "github.com/owm-inc/weatherchecker-go/db"
        "github.com/owm-inc/weatherchecker-go/structs"
        )

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

func PrintHistory(w http.ResponseWriter) {
    MarshalPrintStuff(history.ShowFullHistory(), w)
}

func PrintLocationEntry(locationEntry structs.LocationEntry, w http.ResponseWriter) {
    MarshalPrintStuff(locationEntry, w)
}

func PrintHistoryEntry(historyEntry structs.HistoryEntry, w http.ResponseWriter) {
    MarshalPrintStuff(historyEntry, w)
}

func PrintLocations(w http.ResponseWriter) {
    MarshalPrintStuff(locations.RetrieveLocations(), w)
}

func PrintStatus(err error, successMessage string, w http.ResponseWriter) {
    err_msg := make(map[string]string)
    if err != nil {
        err_msg["status"] = "500"
        err_msg["message"] = err.Error()
    } else {
        err_msg["status"] = "200"
        err_msg["message"] = successMessage
    }

    MarshalPrintStuff(err_msg, w)
}

func GetHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    PrintHistory(w)
}

func GetLocations(c web.C, w http.ResponseWriter, r *http.Request) {
    PrintLocations(w)
}

func AddLocation(c web.C, w http.ResponseWriter, r *http.Request) {
    missing := make([]string, 0)

    query_holder := r.URL.Query()

    city_name := query_holder.Get("city_name") ; if city_name == "" {missing = append(missing, "city name")}
    iso_country := query_holder.Get("iso_country") ; if iso_country == "" {missing = append(missing, "country code")}
    country_name := query_holder.Get("country_name") ; if country_name == "" {missing = append(missing, "country name")}
    latitude := query_holder.Get("latitude") ; if latitude == "" {missing = append(missing, "latitude")}
    longitude := query_holder.Get("longitude") ; if longitude == "" {missing = append(missing, "longitude")}
    accuweather_id := query_holder.Get("accuweather_id")
    accuweather_city_name := query_holder.Get("accuweather_city_name")
    gismeteo_id := query_holder.Get("gismeteo_id")
    gismeteo_city_name := query_holder.Get("gismeteo_city_name")

    if len(missing) > 0 {
        err_msg := make(map[string]string)
        err_msg["status"] = "500"
        err_msg["message"] = "The following parameters are missing: " + strings.Join(missing, ", ")
        MarshalPrintStuff(err_msg, w)
    } else {
        locationEntry := locations.AddLocation (city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name)
        PrintLocationEntry(locationEntry, w)
    }
}

func RemoveLocation(c web.C, w http.ResponseWriter, r *http.Request) {
    missing := make([]string, 0)

    query_holder := r.URL.Query()

    location_id := query_holder.Get("location_id") ; if location_id == "" {missing = append(missing, "location_id")}

    if len(missing) > 0 {
        err_msg := make(map[string]string)
        err_msg["status"] = "500"
        err_msg["message"] = "The following parameters are missing: " + strings.Join(missing, ", ")
        MarshalPrintStuff(err_msg, w)
    } else {
        err := locations.RemoveLocation (location_id)
        PrintStatus(err, "Location removed successfully.", w)
    }
}

func ClearLocations(c web.C, w http.ResponseWriter, r *http.Request) {
    err := locations.Clear()

    PrintStatus(err, "Locations cleared successfully.", w)
}

func RefreshHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    query_holder := r.URL.Query()

    var wtypes []string
    wtype := query_holder.Get("wtype")

    if wtype == "" {
        wtypes = []string{"current", "forecast"}
    } else if (wtype != "current" && wtype != "forecast" && wtype != "") {
        PrintStatus(errors.New(""), "Invalid wtype specified.", w); return
    } else {
        wtypes = []string{wtype}
    }

    locations_query := locations.RetrieveLocations()
    historyEntry := history.AddHistoryEntry(locations_query, sources, wtypes)
    PrintHistoryEntry(historyEntry, w)
}

func ClearHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    err := history.Clear()

    PrintStatus(err, "History cleared successfully.", w)
}

func Api(c *web.C, h http.Handler) http.Handler {
    fn := func (w http.ResponseWriter, r *http.Request) {
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
    goji.Get(LocationEntrypoint, GetLocations)
    goji.Get(LocationEntrypoint + "/add", AddLocation)
    goji.Get(LocationEntrypoint + "/remove", RemoveLocation)
    goji.Get(LocationEntrypoint + "/clear", ClearLocations)
    goji.Get(HistoryEntrypoint, GetHistory)
    goji.Get(HistoryEntrypoint + "/refresh", RefreshHistory)
    goji.Get(HistoryEntrypoint + "/clear", ClearHistory)
    goji.Serve()
}
