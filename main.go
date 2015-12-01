package main

import (
        "encoding/json"
        "fmt"
        "net/http"

        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"

        "github.com/owm-inc/weatherchecker-go/structs"
        )

var proxyTable = structs.NewWeatherProxyTable(locations, sources)
var locations = structs.LoadLocations()
var sources = structs.CreateSources()
var history = structs.NewWeatherHistory()


func MarshalPrintStuff(stuff interface{}, w http.ResponseWriter) {
    data, _ := json.Marshal(stuff)
    jsonString := string(data)
    fmt.Fprintf(w, jsonString)
}

func PrintProxies(w http.ResponseWriter) {
    MarshalPrintStuff(proxyTable, w)
}

func PrintHistory(w http.ResponseWriter) {
    MarshalPrintStuff(history, w)
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

func main() {
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
