"use strict"

requirejs(["charts"])

let STATUS = {
    OK: 0,
    LOADING: 1,
    ERROR: 2
}

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
            type: "text",
            class: inputfield.Name + " form-control",
            value: inputfield.Default,
            placeholder: inputfield.Placeholder
        })
        input_fields.push(entry)
    }

    return input_fields
}

function set_spinner_status(spinnerContainer, status) {
    spinnerContainer.empty()
    let iconClass = ""
    switch (status) {
        case STATUS.OK: // OK
            iconClass = "fa fa-check"
            break
        case STATUS.LOADING: // Loading
            iconClass = "fa fa-spin fa-refresh"
            break
        case STATUS.ERROR: // Error
            iconClass = "fa fa-minus-circle"
            break
        default:
            return
    }
    spinnerContainer.append($('<span>', {
        class: iconClass
    }))
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
        console.log(data)
    }

    function reload_server_uri() {
        let APIEP = "api"
        let APIVER = "0.1"
        serveraddr = ""
        let serverEP = serveraddr + "/" + APIEP + "/" + APIVER
        entrypoints.locations = serverEP + "/" + "locations"
        entrypoints.history = serverEP + "/" + "history"
    }

    let location_list_model_id = "select.location_list"
    let location_list_model = $(location_list_model_id)

    function refresh_location_list() {
        location_list_model.empty()

        let output = new String

        $.get(entrypoints.locations, function(data) {
            output = data
            let data_object = $.parseJSON(data)
            let locations = getlocations(data_object)

            for (let entry of locations) {
                let entryOption = $("<option>", {
                    value: entry['id']
                })
                entryOption.append(entry['name'])
                location_list_model.append(entryOption)
            }

        })
        return output
    }

    function refresh_location_list_log() {
        let data = refresh_location_list()
        logger(data)
    }


    let get_weatherdata_spinner = $('.get_weatherdata_spinner')

    // Actions on page load
    reload_server_uri()
    refresh_location_list_log()
    set_spinner_status(get_weatherdata_spinner, STATUS.OK)

    // Events

    $(".refresh_button").click(function() {
        $.get(entrypoints.history + "/refresh", function(data) {
            logger(data)
        })
    });

    function refresh_upsert_form(form, upsert_type) {
        form.empty()

        let inputarea = $('<div>', {
            class: 'inputarea'
        })
        let inputfields = new Array
        if (upsert_type == 0) {
            inputfields = create_input_fields(location_add_inputfields)
        } else {
            inputfields = create_input_fields(location_update_inputfields)
        }
        for (let field of inputfields) {
            let group = $('<div>', {
                class: 'form-group'
            })
            group.append(field)
            inputarea.append(group)
        }

        let buttonarea = $('<div>', {
            class: 'buttonarea'
        })
        let cancelButton = $("<input>", {
            type: "button",
            class: "location_upsert_cancel btn btn-danger",
            value: "Отмена"
        })
        let sendButton = $("<input>", {
            type: "submit",
            class: "location_upsert_send btn btn-default",
            value: "Отправить"
        })
        cancelButton.click(function() {
            form.empty()
        })

        buttonarea.append(cancelButton)
        buttonarea.append(sendButton)

        form.append($('<br>'))
        form.append(inputarea)
        form.append(buttonarea)
        form.append($('<hr>'))
    }

    let location_upsert_form = $(".location_upsert_form")
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
        refresh_location_list()
    })

    $(".upsert_location").click(function() {
        event.preventDefault()
        refresh_upsert_form($(".location_upsert_form"), 1)
    })

    $(".location_data_download").click(refresh_location_list_log)

    $("form.weather").submit(function(event) {
        event.preventDefault();
        let locationid = $(location_list_model_id + " option:selected").val()
        let wtype = "current"
        set_spinner_status(get_weatherdata_spinner, STATUS.LOADING)
        $.ajax({
            url: entrypoints.history + "?" + "locationid=" + locationid + "&" + "wtype=" + wtype,
            success: function(data) {
                let jsonData = $.parseJSON(data)
                let status = jsonData['status']
                let message = jsonData['message']
                let content = jsonData['content']
                $(".weathertable").empty()
                $(".weatherchart").empty()
                logger(data)
                if (status != 200) {
                    set_spinner_status(get_weatherdata_spinner, STATUS.ERROR)
                    logger("Request failed with status " + String(status) + " and message: " + message)
                } else {
                    set_spinner_status(get_weatherdata_spinner, STATUS.OK)
                    $(".weathertable").append(build_weather_table(content['history']))
                    build_weather_chart($('.weatherchart'), content['history'])
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                set_spinner_status(get_weatherdata_spinner, STATUS.ERROR)
                logger("Ошибка запроса к " + url + ":   " + textStatus)
            }
        })
    });

    function build_weather_table(historyObject) {
        let table = $("<table>")

        let table_elements = [{
            id: "json_link",
            name: "Запись в БД"
        }, {
            id: "source",
            name: "Погодный сервис"
        }, {
            id: "raw_link",
            name: "Источник"
        }, {
            id: "dt",
            name: "Дата измерений"
        }, {
            id: "request_dt",
            name: "Дата запроса"
        }, {
            id: "temp",
            name: "Температура, C"
        }, {
            id: "pressure",
            name: "Давление, бар"
        }, {
            id: "humidity",
            name: "Влажность, процентов"
        }, {
            id: "wind_speed",
            name: "Скорость ветра, м/с"
        }, {
            id: "precipitation",
            name: "Осадки, мм"
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
                }) + "'>" + "Открыть" + "</a>",
                "source": history_entry['source']['name'],
                "raw_link": "<a href='" + history_entry['url'] + "'>Открыть</a>",
                "dt": new Date(history_entry['measurements'][0]['timestamp'] * 1000).toISOString(),
                "request_dt": new Date(history_entry['request_time'] * 1000).toISOString(),
                "temp": history_entry['measurements'][0]['data']['temp'].toFixed(1),
                "pressure": history_entry['measurements'][0]['data']['pressure'].toFixed(1),
                "humidity": history_entry['measurements'][0]['data']['humidity'].toFixed(1),
                "wind_speed": history_entry['measurements'][0]['data']['wind'].toFixed(1),
                "precipitation": history_entry['measurements'][0]['data']['precipitation'].toFixed(1)
            }

            for (let row_cell of table_elements) {
                let text = history_entry_elements[row_cell.id]
                history_entry_row.append("<td>" + text + "</td>")
            }

            tbody.append(history_entry_row)
        }
        table.append(tbody)
        table.DataTable()

        return table
    }
});
