package main

import (
        "encoding/json"
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
        err_msg["Status"] = "500"
        err_msg["Message"] = "The following parameters are missing: " + strings.Join(missing, ", ")
        MarshalPrintStuff(err_msg, w)
    } else {
        locationEntry := locations.AddLocation (city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name)
        PrintLocationEntry(locationEntry, w)
    }
}

func RefreshHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    wtypes := []string{"current"}
    locations_query := locations.RetrieveLocations()
    historyEntry := history.AddHistoryEntry(locations_query, sources, wtypes)
    PrintHistoryEntry(historyEntry, w)
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

    const ApiEntrypoint = "/api"

    const DataEntrypoint = ApiEntrypoint + "/data"
    const ActionEntrypoint = ApiEntrypoint + "/actions"

    goji.Use(Api)
    goji.Get(DataEntrypoint + "/locations", GetLocations)
    goji.Get(DataEntrypoint + "/history", GetHistory)
    goji.Get(ActionEntrypoint + "/add_location", AddLocation)
    goji.Get(ActionEntrypoint + "/refresh_history", RefreshHistory)
    goji.Serve()
}
