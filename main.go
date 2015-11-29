package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"

        "github.com/skybon/weatherchecker-go/adapters"
        )

func download_data(url string) string {
    resp, err := http.Get(url)
    contents := ""
    if err != nil {
        fmt.Println(`Request finished with error`, err)
    } else {
        defer resp.Body.Close()
        readall_contents, _ := ioutil.ReadAll(resp.Body)
        contents = string(readall_contents)
    }
    return contents
}

func main() {
    var url_base = "http://pro.openweathermap.org/data/2.5/weather?q=Moscow,RU&appid="
    var keys = make(map[string]string)
    keys["owm"] = os.Getenv("OWM_KEY")
    var appid = keys["owm"]
    var data = download_data(url_base + appid)

    var dataset = adapters.Owm_adapt_weather(data)
    fmt.Println(dataset)
}
