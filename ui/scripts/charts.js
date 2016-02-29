"use strict";

exports.buildWeatherChart = buildWeatherChart

function makeProviderMap(content, key) {
    var providerMap = new Map;

    content.forEach(function (entry) {
        var provider = entry['source']['prettyname'];

        var measurement = entry['measurements'][0];
        if (measurement == undefined) {
            return;
        };

        var dt = Number(measurement['timestamp']) * 1000;
        var data = measurement['data'][key];

        if (this.get(provider) == undefined) {
            this.set(provider, []);
        };

        var dataEntry = [dt, data];

        this.get(provider).push(dataEntry);
    }, providerMap);

    return providerMap;
}

function makeSeriesArray(historyDataObject, key) {
    var seriesArray = [];

    makeProviderMap(historyDataObject, key).forEach(function (providerData, providerName) {
        var providerEntry = {
            data: providerData,
            name: providerName
        };

        this.push(providerEntry);
    }, seriesArray);

    return seriesArray;
};

function makeWeatherChartData(historyDataObject) {
    if (historyDataObject == undefined) {
        return [];
    } else {
        return makeSeriesArray(historyDataObject, 'temp');
    };
};

function buildWeatherChart(containerObject, historyDataObject) {
    jQuery(containerObject).highcharts({
        chart: {
            type: 'spline'
        },
        title: {
            text: ''
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
        series: makeWeatherChartData(historyDataObject)
    });
};
