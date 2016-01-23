"use strict";

exports.build_weather_chart = build_weather_chart

function makeProviderObject(content, key) {
    var providerObject = {}

    for (var entry of content) {
        var provider = entry['source']

        var measurement = entry['measurements'][0]
        if (measurement == undefined) {
            continue
        }

        var dt = Number(measurement['timestamp']) * 1000
        var data = measurement['data'][key]

        if (providerObject[provider] == undefined) {
            providerObject[provider] = []
        }

        var dataEntry = [dt, data]

        providerObject[provider].push(dataEntry)
    }

    return providerObject
}

function makeSeriesObject(providerObject) {
    var seriesObject = []

    for (var providerName in providerObject) {
        var providerEntry = {
            name: providerName,
            data: providerObject[providerName]
        }

        seriesObject.push(providerEntry)
    }

    return seriesObject
}

function get_weatherchart_data(historyObject) {
    if (historyObject == undefined) {
        return []
    }
    var content = historyObject['data']
    var series = makeSeriesObject(makeProviderObject(content, 'temp'))

    return series
}

function build_weather_chart(containerObject, historyObject) {
    var chart_series = get_weatherchart_data(historyObject)
    containerObject.highcharts({
        chart: {
            type: 'spline'
        },
        title: {
            text: 'Weather chart'
        },
        xAxis: {
            type: 'datetime',
            title: {
                text: 'Date'
            }
        },
        yAxis: {
            title: {
                text: 'Temperature'
            }
        },
        series: chart_series
    })
}
