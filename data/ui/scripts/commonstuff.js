"use strict";

exports.make_history_table_data = make_history_table_data;

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
