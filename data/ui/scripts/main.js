"use strict"

requirejs(["charts"])

function getlocations(data) {
    let locations = new Array
    let location_list = data['content']['locations']

    for (let location_entry of location_list) {
        let entry = new Object
        entry.id = location_entry['objectid']
        entry.name = location_entry['city_name']
        locations.push(entry)
    }

    return locations
}

function create_input_fields(inputfields_desc) {
    let input_fields = new Array
    for (let inputfield of inputfields_desc) {
        let entry = $('<input>', {
            name: inputfield.Name,
            class: inputfield.Name + " entry",
            value: inputfield.Default,
            placeholder: inputfield.Placeholder
        })
        input_fields.push(entry)
    }

    return input_fields
}

$(document).ready(function() {
    let serveraddr = new String

    let entrypoints = {
        locations: new String,
        history: new String
    }

    let location_add_inputfields = [{
        Name: "city_name",
        Default: "",
        Placeholder: "Название города"
    }, {
        Name: "iso_country",
        Default: "",
        Placeholder: "Код страны"
    }, {
        Name: "country_name",
        Default: "",
        Placeholder: "Название страны"
    }, {
        Name: "latitude",
        Default: "",
        Placeholder: "Широта"
    }, {
        Name: "longitude",
        Default: "",
        Placeholder: "Долгота"
    }, {
        Name: "accuweather_id",
        Default: "",
        Placeholder: "ID AccuWeather"
    }, {
        Name: "accuweather_city_name",
        Default: "",
        Placeholder: "Название города AccuWeather"
    }, {
        Name: "gismeteo_id",
        Default: "",
        Placeholder: "ID Gismeteo"
    }, {
        Name: "gismeteo_city_name",
        Default: "",
        Placeholder: "Название города Gismeteo"
    }]

    let location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (для редактирования)"
    }].concat(location_add_inputfields)

    function logger(data) {
        $("pre.logger").text(data)
    }

    function reload_server_uri() {
        let APIEP = "api"
        let APIVER = "0.1"
        serveraddr = "" //$("input.server_uri")[0].value
        let serverEP = serveraddr + "/" + APIEP + "/" + APIVER
        entrypoints.locations = serverEP + "/" + "locations"
        entrypoints.history = serverEP + "/" + "history"
    }
    reload_server_uri()
        //$("input.server_uri").change(reload_server_uri)

    let model_object_id = "select.location_list"
    let model_object = $(model_object_id)

    $("form.serveraddr").submit(function() {
        event.preventDefault()
    })

    $("form.refresh > input.refresh_button").click(function() {
        $.get(entrypoints.history + "/refresh", function(data) {
            logger(data)
        })
    });

    function refresh_upsert_form(form, upsert_type) {
        form.empty()
        let inputfields = new Array
        if (upsert_type == 0) {
            inputfields = create_input_fields(location_add_inputfields)
        } else {
            inputfields = create_input_fields(location_update_inputfields)
        }
        for (let field of inputfields) {
            form.append(field)
        }
        let cancelButton = $("<input>", {
            type: "button",
            class: "location_upsert_cancel",
            value: "Отмена"
        })
        let sendButton = $("<input>", {
            type: "submit",
            class: "location_upsert_send",
            value: "Отправить"
        })
        cancelButton.click(function() {
            form.empty()
        })

        form.append(cancelButton)
        form.append(sendButton)
    }

    let location_upsert_form = $("form.location_upsert")
    location_upsert_form.submit(function() {
        event.preventDefault()
        let params = location_upsert_form.serialize()
        let url = entrypoints.locations + "/upsert"
        $.ajax({
            url: url + "?" + params,
            success: function(data) {
                logger(data)
            },
            error: function(jqXHR, textStatus, errorThrown) {
                logger("Ошибка запроса к " + url + ":   " + textStatus)
            }
        })
    })

    $("form.location > input.location_upsert").click(function() {
        event.preventDefault()
        refresh_upsert_form($("form.location_upsert"), 1)
    })

    $("form.location > input.location_data_download").click(function() {
        model_object.empty();

        $.get(entrypoints.locations, function(data) {
            logger(data)
            let data_object = $.parseJSON(data)
            let locations = getlocations(data_object)

            for (let entry of locations) {
                let entryOption = $("<option>", {
                    value: entry.id
                })
                entryOption.append(entry.name)
                model_object.append(entryOption)
            }

        })
    });

    $("form.weather").submit(function(event) {
        event.preventDefault();
        let locationid = $(model_object_id + " option:selected")[0].value
        let wtype = "current"
        $.get(entrypoints.history + "?" + "locationid=" + locationid + "&" + "wtype=" + wtype, function(data) {
            let jsonData = $.parseJSON(data)
            let status = jsonData['status']
            let message = jsonData['message']
            let content = jsonData['content']
            $(".weathertable").empty()
            $(".weatherchart").empty()
            logger(data)
            if (status != 200) {
                logger("Request failed with status " + String(status) + " and message: " + message)
            } else {
                $(".weathertable").append(build_weather_table(content['history']))
                build_weather_chart($('.weatherchart'), content['history'])
            }
        })
    });



    function build_weather_table(historyObject) {
        let table = $("<table>")

        let table_elements = [{
            id: "json_link",
            name: "ObjectId"
        }, {
            id: "source",
            name: "Источник"
        }, {
            id: "raw_link",
            name: "Ссылка на источник"
        }, {
            id: "dt",
            name: "Дата измерений"
        }, {
            id: "request_dt",
            name: "Дата запроса"
        }, {
            id: "temp",
            name: "Температура, C"
        }]
        let thead = $("<thead>")
        let theadtr = $("<tr>")
        for (let element of table_elements) {
            theadtr.append("<td>" + element.name + "</td>")
        }
        thead.append(theadtr)
        table.append(thead)

        let content = historyObject['data']
        let tbody = $("<tbody>")
        for (let history_entry of content) {
            if (history_entry['status'] != 200) {
                continue
            }
            let history_entry_row = $("<tr>")

            let history_entry_elements = {
                "json_link": "<a href='" + entrypoints.history + "?" + $.param({
                    entryid: history_entry['objectid']
                }) + "'>" + history_entry['objectid'] + "</a>",
                "source": history_entry['source']['name'],
                "raw_link": "<a href='" + history_entry['url'] + "'>Открыть</a>",
                "dt": history_entry['measurements'][0]['timestamp'],
                "request_dt": history_entry['request_time'],
                "temp": history_entry['measurements'][0]['data']['temp']
            }

            for (let row_cell of table_elements) {
                let text = history_entry_elements[row_cell.id]
                history_entry_row.append("<td>" + text + "</td>")
            }

            tbody.append(history_entry_row)
        }
        table.append(tbody)

        return table
    }
});
