"use strict";

var helpers = require("./helpers.js");

exports.make_history_table_data = make_history_table_data;
exports.refresh_location_list = refresh_location_list;
exports.parseOWMhistory = parseOWMhistory;

function make_history_table_data(historyObject, history_entrypoint) {
    var content = historyObject.data;
    var values = [];

    for (var history_entry of content) {
        var history_entry_flat = {
            "json_link": "<a href='" + history_entrypoint + "?" + "entryid=" + history_entry.objectid + "'>Open</a>",
            "source_id": history_entry.source.name,
            "source_name": history_entry.source.prettyname,
            "raw_link": "N/A",
            "dt": moment.unix(history_entry.measurements[0].timestamp).format('YYYY-MM-DD HH:mm ZZ'),
            "request_dt": moment.unix(history_entry.request_time).format('YYYY-MM-DD HH:mm ZZ'),
            "temp": history_entry.measurements[0].data.temp.toFixed(1),
            "pressure": history_entry.measurements[0].data.pressure.toFixed(1),
            "humidity": history_entry.measurements[0].data.humidity.toFixed(1),
            "wind_speed": history_entry.measurements[0].data.wind.toFixed(1),
            "precipitation": history_entry.measurements[0].data.precipitation.toFixed(1)
        };
        if (history_entry.url !== undefined) {
            history_entry_flat.raw_link = "<a href='" + history_entry.url + "'>Open</a>";
        };

        values.push(history_entry_flat);
    };

    var columns = [{
        data: "source_name",
        title: "Provider",
        orderable: true
    }, {
        data: "dt",
        title: "Measurement date",
        orderable: true
    }, {
        data: "temp",
        title: "Temperature, C",
        orderable: false
    }, {
        data: "pressure",
        title: "Pressure, mbar",
        orderable: false
    }, {
        data: "humidity",
        title: "Humidity, percent",
        orderable: false
    }, {
        data: "wind_speed",
        title: "Wind speed, m/s",
        orderable: false
    }, {
        data: "precipitation",
        title: "Precipitation, mm",
        orderable: false
    }]

    var tableOpts = {
        order: [
            [1, "desc"],
            [0, "asc"]
        ]
    }

    return [values, columns, tableOpts]
}

function getlocations(dataObject) {
    var locations = [];
    var location_list = dataObject['content']['locations'];

    if (location_list != null) {
        for (var location_entry of location_list) {
            var entry = {};
            entry.id = location_entry['objectid'];
            entry.slug = location_entry['slug'];
            entry.name = location_entry['city_name'];
            entry.latitude = location_entry['latitude'];
            entry.longitude = location_entry['longitude'];
            locations.push(entry);
        };
    };

    return locations;
};

function makeCitySelectOptions(locationCollection) {
    var output = [];

    for (var entry of locationCollection) {
        var newOption = document.createElement('option')
        newOption.setAttribute('objectid', entry.id)
        newOption.setAttribute('slug', entry.slug)
        newOption.setAttribute('lat', entry.latitude)
        newOption.setAttribute('lon', entry.longitude)

        newOption.textContent = entry.name;
        output.push(newOption);
    };
    return output;
};

function refresh_location_list(location_list_model, entrypoints, spinnerContainer, callback) {
    helpers.clearChildren(location_list_model);

    var dataObject = {};
    var locationCollection = []

    helpers.get_with_spinner_and_callback(entrypoints.locations, spinnerContainer, function (data) {
        dataObject = JSON.parse(data);

        locationCollection = getlocations(dataObject);
        var options = makeCitySelectOptions(locationCollection);

        for (var option of options) {
            location_list_model.appendChild(option);
        };
        var locationMap = helpers.collectionToMap(locationCollection, 'id');

        if (callback != undefined) {
            callback(locationMap);
        };
    })
};


function parseOWMhistory(OWMHistoryObject) {
    var checkerHistory = []

    var historyList = OWMHistoryObject['list']
    for (var historyEntry of historyList) {
        var newEntry = {}
        newEntry['source'] = 'owm_history'
        newEntry['measurements'] = []

        var measurement = {}
        measurement['timestamp'] = historyEntry['dt']

        var measurementData = {}
        measurementData['temp'] = historyEntry['main']['temp'] - 273.15
        measurementData['wind_speed'] = historyEntry['wind']['speed']
        measurementData['humidity'] = historyEntry['main']['humidity']
        measurementData['pressure'] = historyEntry['main']['pressure']

        measurement['data'] = measurementData

        newEntry['measurements'].push(measurement)

        checkerHistory.push(newEntry)
    }
    return checkerHistory
}

function extractForecast(forecastArray, length) {
    var newForecasts = [];

    for (var forecastEntry of forecastArray) {
        var newEntry = {};

        var startDate = forecastEntry['request_dt'];

        newEntry['source'] = forecastEntry['source'];
        newEntry['request_dt'] = startDate;
        newEntry['measurements'] = [];

        var closestDate = helpers.find_closest(startDate + length, make_timestamplist(forecastEntry)).Closest;

        var matchedForecast = forecastEntry['measurements'].find(function (e, i, a) {
            if (e['timestamp'] == closestDate) {
                return true;
            };
        });

        newEntry['measurements'].push(matchedForecast);

        newForecasts.push(newEntry)
    };

    return newForecasts
};

function make_timestamplist(historyObject) {
    var timestamplist = []

    for (var measurement of historyObject['measurements']) {
        var timestamp = measurement['timestamp']
        timestamplist.push(timestamp)
    }

    return timestamplist
}
