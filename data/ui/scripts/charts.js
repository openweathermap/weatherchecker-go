"use strict"

define(function() {
    return {
        make_providers_list: make_providers_list,
        make_timestamps_object:  make_timestamps_object,
        get_weatherchart_data: get_weatherchart_data,
        build_weather_chart: build_weather_chart
    }
})

let make_providers_list = function(content) {
    let providersSet = new Set
    for (let entry of content) {
        let provider = entry["source"]["name"]
        providersSet.add(provider)
    }
    return Array.from(providersSet)
}

let make_timestamps_object = function(content) {
    let timestampsObject = new Object
    for (let entry of content) {
        let measurement = entry['measurements'][0]

        if (measurement == undefined) {
            continue
        }

        let dt = String(measurement['timestamp'])

        if (timestampsObject[dt] == undefined) {
            timestampsObject[dt] = new Object
        }

        let provider = entry["source"]["name"]
        timestampsObject[dt][provider] = measurement["data"]["temp"]
    }

    return timestampsObject
}

let make_series_object = function(providers) {
    let seriesObject = new Object

    for (let provider of providers) {
        seriesObject[provider] = {
            name: provider,
            data: new Array
        }
    }
    return seriesObject
}

let get_weatherchart_data = function(historyObject) {
    let series = new Array
    let timestamps = new Array

    if (historyObject == undefined) {
        return series, timestamps
    }

    let content = historyObject['data']

    let providers = make_providers_list(content)

    let timestampsObject = make_timestamps_object(content)

    let seriesObject = make_series_object(providers)

    timestamps = Object.keys(timestampsObject)

    for (let key of timestamps) {
        let dt = timestampsObject[key]
        for (let provider of providers) {
            if (seriesObject[provider].data == undefined) {
                seriesObject[provider].data = new Array
            }
            let data = dt[provider]
            let temp = data
            if (data == undefined) {
                continue
            }
            let unixDate = Number(key) * 1000
            seriesObject[provider].data.push([unixDate, temp])
        }
    }

    for (let k in seriesObject) {
        let v = seriesObject[k]
        series.push(v)
    }

    return series
}

let build_weather_chart = function(containerObject, historyObject) {
    let chart_series = get_weatherchart_data(historyObject)

    containerObject.highcharts({
        chart: {
            type: 'spline'
        },
        title: {
            text: 'Таблица погоды'
        },
        xAxis: {
            type: 'datetime',
            title: {
                text: 'Дата'
            }
        },
        yAxis: {
            title: {
                text: 'Температура'
            }
        },
        series: chart_series
    })
}
