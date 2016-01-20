"use strict"

define(function() {
    return {
        make_providers_list: make_providers_list,
        make_timestamps_object:  make_timestamps_object,
        get_weatherchart_data: get_weatherchart_data,
        build_weather_chart: build_weather_chart
    }
})

var make_providers_list = function(content) {
    var providersSet = new Set
    for (var entry of content) {
        var provider = entry["source"]["name"]
        providersSet.add(provider)
    }
    return Array.from(providersSet)
}

var make_timestamps_object = function(content) {
    var timestampsObject = new Object
    for (var entry of content) {
        var measurement = entry['measurements'][0]

        if (measurement == undefined) {
            continue
        }

        var dt = String(measurement['timestamp'])

        if (timestampsObject[dt] == undefined) {
            timestampsObject[dt] = new Object
        }

        var provider = entry["source"]["name"]
        timestampsObject[dt][provider] = measurement["data"]["temp"]
    }

    return timestampsObject
}

var make_series_object = function(providers) {
    var seriesObject = new Object

    for (var provider of providers) {
        seriesObject[provider] = {
            name: provider,
            data: new Array
        }
    }
    return seriesObject
}

var get_weatherchart_data = function(historyObject) {
    var series = new Array
    var timestamps = new Array

    if (historyObject == undefined) {
        return series, timestamps
    }

    var content = historyObject['data']

    var providers = make_providers_list(content)

    var timestampsObject = make_timestamps_object(content)

    var seriesObject = make_series_object(providers)

    timestamps = Object.keys(timestampsObject)

    for (var key of timestamps) {
        var dt = timestampsObject[key]
        for (var provider of providers) {
            if (seriesObject[provider].data == undefined) {
                seriesObject[provider].data = new Array
            }
            var data = dt[provider]
            var temp = data
            if (data == undefined) {
                continue
            }
            var unixDate = Number(key) * 1000
            seriesObject[provider].data.push([unixDate, temp])
        }
    }

    for (var k in seriesObject) {
        var v = seriesObject[k]
        series.push(v)
    }

    return series
}

var build_weather_chart = function(containerObject, historyObject) {
    var chart_series = get_weatherchart_data(historyObject)

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
