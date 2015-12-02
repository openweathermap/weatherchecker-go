package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "net/http"
        "os"

        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"

        "github.com/owm-inc/weatherchecker-go/db"
        "github.com/owm-inc/weatherchecker-go/structs"
        )

var mongoDsn string

var db_instance = db.Db()

var proxyTable = structs.NewWeatherProxyTable(locations, sources)
var locations = structs.LoadLocations()
var sources = structs.CreateSources()
var history = structs.NewWeatherHistory(db_instance)


func MarshalPrintStuff(stuff interface{}, w http.ResponseWriter) {
    data, _ := json.Marshal(stuff)
    jsonString := string(data)
    fmt.Fprintf(w, jsonString)
}

func PrintProxies(w http.ResponseWriter) {
    MarshalPrintStuff(proxyTable, w)
}

func PrintHistory(w http.ResponseWriter) {
    MarshalPrintStuff(history.ShowFullHistory(), w)
}

func PrintHistoryEntry(historyEntry structs.HistoryEntry, w http.ResponseWriter) {
    MarshalPrintStuff(historyEntry, w)
}

func GetHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    PrintHistory(w)
}

func RefreshHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    proxyTable.Refresh()
    historyEntry := history.AddHistoryEntry(proxyTable.Table)
    PrintHistoryEntry(historyEntry, w)
}

func GetProxies(c web.C, w http.ResponseWriter, r *http.Request) {
    PrintProxies(w)
}

func RefreshProxies(c web.C, w http.ResponseWriter, r *http.Request) {
    proxyTable.Refresh()
    PrintProxies(w)
}

func Api(c *web.C, h http.Handler) http.Handler {
    fn := func (w http.ResponseWriter, r *http.Request) {
        // Pass data through the environment
        c.Env["proxyTable"] = &proxyTable
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
        fmt.Println(fmt.Sprintf("db error: %s", err))
        return
	}
	defer db_instance.Disconnect()

    const ApiEntrypoint = "/api"

    const DataEntrypoint = ApiEntrypoint + "/data"
    const ActionEntrypoint = ApiEntrypoint + "/actions"

    goji.Use(Api)
    goji.Get(DataEntrypoint + "/proxies", GetProxies)
    goji.Get(DataEntrypoint + "/history", GetHistory)
    goji.Get(ActionEntrypoint + "/refresh_history", RefreshHistory)
    goji.Get(ActionEntrypoint + "/refresh_proxies", RefreshProxies)
    goji.Serve()
}
