package main

import (
        "encoding/json"
        "fmt"
        "net/http"

        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"

        "github.com/skybon/weatherchecker-go/structs"
        )

var proxyTable = structs.NewWeatherProxyTable(locations, sources)
var locations = structs.LoadLocations()
var sources = structs.CreateSources()
var history = structs.WeatherHistory{}


func GetHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    data, _ := json.Marshal(history)
    jsonString := string(data)
    fmt.Fprintf(w, jsonString)
}

func RefreshHistory(c web.C, w http.ResponseWriter, r *http.Request) {
    RefreshProxies(c, w, r)

    historyEntry := history.AddHistoryEntry(proxyTable.Table)
    data, _ := json.Marshal(historyEntry)
    jsonString := string(data)
    fmt.Fprintf(w, jsonString)
}


func GetProxies(c web.C, w http.ResponseWriter, r *http.Request) {
    data, _ := json.Marshal(proxyTable)
    jsonString := string(data)
    fmt.Fprintf(w, jsonString)
}

func RefreshProxies(c web.C, w http.ResponseWriter, r *http.Request) {
    proxyTable.Refresh()
    jsonString := `{"cod": "200", "message": "Proxies refreshed successfully"}`
    fmt.Fprintf(w, jsonString)
}

func Api(c *web.C, h http.Handler) http.Handler {
    fn := func (w http.ResponseWriter, r *http.Request) {
        //proxyTable.Refresh()

        //history.AddHistoryEntry(proxyTable.Table)
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
    goji.Get(ActionEntrypoint + "/refresh", RefreshHistory)
    goji.Get(ActionEntrypoint + "/refresh_proxies", RefreshProxies)
    goji.Serve()
}
