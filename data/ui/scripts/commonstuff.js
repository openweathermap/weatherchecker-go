"use strict";

var helpers = require("./helpers.js");

exports.make_history_table_data = make_history_table_data;
exports.make_city_select_options = make_city_select_options;
exports.refresh_location_list = refresh_location_list;
exports.parseOWMhistory = parseOWMhistory;

function make_history_table_data(historyObject, history_entrypoint) {
    var content = historyObject.data;
    var values = [];

    for (var history_entry of content) {

        var history_entry_flat = {
            "json_link": "<a href='" + history_entrypoint + "?" + $.param({
                entryid: history_entry.objectid
            }) + "'>Open</a>",
            "source": history_entry.source,
            "raw_link": "N/A",
            "dt": new Date(history_entry.measurements[0].timestamp * 1000).toISOString(),
            "request_dt": new Date(history_entry.request_time * 1000).toISOString(),
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
        data: "json_link",
        title: "DB entry"
    }, {
        data: "source",
        title: "Provider"
    }, {
        data: "raw_link",
        title: "Source"
    }, {
        data: "dt",
        title: "Measurement date"
    }, {
        data: "request_dt",
        title: "Request date"
    }, {
        data: "temp",
        title: "Temperature, C"
    }, {
        data: "pressure",
        title: "Pressure, mbar"
    }, {
        data: "humidity",
        title: "Humidity, percent"
    }, {
        data: "wind_speed",
        title: "Wind speed, m/s"
    }, {
        data: "precipitation",
        title: "Precipitation, mm"
    }]

    return [values, columns]
}

function getlocations(dataObject) {
    var locations = [];
    var location_list = dataObject['content']['locations'];

    if (location_list != null) {
        for (var location_entry of location_list) {
            var entry = {};
            entry.id = location_entry['objectid'];
            entry.name = location_entry['city_name'];
            entry.latitude = location_entry['latitude'];
            entry.longitude = location_entry['longitude'];
            locations.push(entry);
        };
    };

    return locations;
};

function make_city_select_options(locationCollection) {
    var output = [];

    for (var entry of locationCollection) {
        var newOption = $("<option>", {
            value: entry.id,
            lat: entry.latitude,
            lon: entry.longitude
        });
        newOption.append(entry.name);
        output.push(newOption);
    };
    return output;
};

function refresh_location_list(location_list_model, entrypoints, spinnerContainer, callback) {
    location_list_model.empty();

    var dataObject = {};
    var locationCollection = []

    helpers.get_with_spinner_and_callback(entrypoints.locations, spinnerContainer, function (data) {
        dataObject = $.parseJSON(data);

        locationCollection = getlocations(dataObject);
        var options = make_city_select_options(locationCollection);

        for (var option of options) {
            location_list_model.append(option);
        };
        var locationMap = helpers.collectionToMap(locationCollection, 'id');
        callback(locationMap);
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

function make_timestamplist(historyObject) {
    var timestamplist = []

    for (var measurement of historyObject['measurements']) {
        var timestamp = measurement['timestamp']
        timestamplist.push(timestamp)
    }

    return timestamplist
}
