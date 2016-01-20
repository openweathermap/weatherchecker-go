"use strict"

requirejs(["charts"])

var STATUS = {
    OK: 0,
    LOADING: 1,
    ERROR: 2
}

function getlocations(data) {
    var locations = new Array
    var location_list = data['content']['locations']

    for (var location_entry of location_list) {
        var entry = new Object
        entry.id = location_entry['objectid']
        entry.name = location_entry['city_name']
        locations.push(entry)
    }

    return locations
}

function create_input_fields(inputfields_desc) {
    var input_fields = new Array
    for (var inputfield of inputfields_desc) {
        var entry = $('<input>', {
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
    var iconClass = ""
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
    var serveraddr = new String

    var entrypoints = {
        locations: new String,
        history: new String
    }

    var adminKey = new String

    var location_add_inputfields = [{
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
    }, {
        Name: "yandex_id",
        Default: "",
        Placeholder: "ID Яндекс"
    }]

    var location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (для редактирования)"
    }].concat(location_add_inputfields)

    function logger(data) {
        console.log(data)
    }

    function reload_server_uri() {
        var APIEP = "api"
        var APIVER = "0.1"
        serveraddr = ""
        var serverEP = serveraddr + "/" + APIEP + "/" + APIVER
        entrypoints.appid_check = serverEP + "/" + "check_appid"
        entrypoints.locations = serverEP + "/" + "locations"
        entrypoints.history = serverEP + "/" + "history"
    }

    var location_list_model_id = "select.location_list"
    var location_list_model = $(location_list_model_id)

    function refresh_location_list() {
        location_list_model.empty()

        var output = new String

        get_with_spinner_and_callback(entrypoints.locations, location_data_download_spinner, function(data) {
            output = data
            var data_object = $.parseJSON(data)
            var locations = getlocations(data_object)

            for (var entry of locations) {
                var entryOption = $("<option>", {
                    value: entry['id']
                })
                entryOption.append(entry['name'])
                location_list_model.append(entryOption)
            }

        })
        return output
    }

    function refresh_location_list_log() {
        var data = refresh_location_list()
        logger(data)
    }

    function get_with_spinner_and_callback(requestUrl, spinnerObject, callbackFunc) {
        set_spinner_status(spinnerObject, STATUS.LOADING)
        $.ajax({
            url: requestUrl,
            success: function(data) {
                var jsonData = $.parseJSON(data)
                var status = jsonData['status']

                if (status == 200) {
                    set_spinner_status(spinnerObject, STATUS.OK)
                } else {
                    set_spinner_status(spinnerObject, STATUS.ERROR)
                }

                if (callbackFunc != undefined) {
                    callbackFunc(data)
                }

                logger(data)
            },
            error: function(data) {
                set_spinner_status(spinnerObject, STATUS.ERROR)
            }
        })
    }

    var appid_check_form = $(".appid_check_form")

    var refresh_button = $('.refresh_button')
    var upsert_location_button = $('.upsert_location')

    var admin_buttons = [refresh_button, upsert_location_button]

    function disable_admin_buttons() {
        for (var button of admin_buttons) {
            button.attr("disabled", true)
        }
    }

    function enable_admin_buttons() {
        for (var button of admin_buttons) {
            button.attr("disabled", false)
        }
    }

    var appid_check_spinner = $('.appid_check_spinner')
    var refresh_spinner = $('.refresh_spinner')
    var upsert_location_spinner = $('.upsert_location_spinner')
    var location_data_download_spinner = $('.location_data_download_spinner')
    var get_weatherdata_spinner = $('.get_weatherdata_spinner')

    function check_appid(appid) {
        var url = entrypoints.appid_check
        $.ajax({
            url: url + "?appid=" + appid,
            success: function(data) {
                logger(data)
                var content = $.parseJSON(data)
                if (content['status'] == 200) {
                    adminKey = appid_check_form.serializeArray()[0].value
                    set_spinner_status(appid_check_spinner, STATUS.OK)
                    enable_admin_buttons()
                } else {
                    set_spinner_status(appid_check_spinner, STATUS.ERROR)
                    disable_admin_buttons()
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                set_spinner_status(appid_check_spinner, STATUS.ERROR)
                logger("Ошибка запроса к " + url + ":   " + textStatus)
                disable_admin_buttons()
            }
        })
    }

    // Actions on page load
    reload_server_uri()
    refresh_location_list_log()
    disable_admin_buttons()
    check_appid(new String)
    for (var spinner of[location_data_download_spinner, get_weatherdata_spinner]) {
        set_spinner_status(spinner, STATUS.OK)
    }

    // Events
    appid_check_form.submit(function() {
        event.preventDefault()
        check_appid(appid_check_form.serializeArray()[0].value)
    })

    refresh_button.click(function() {
        get_with_spinner_and_callback(entrypoints.history + "/refresh" + "?appid=" + adminKey, refresh_spinner)
    })

    function refresh_upsert_form(form, upsert_type) {
        form.empty()

        var inputarea = $('<div>', {
            class: 'inputarea'
        })
        var inputfields = new Array
        if (upsert_type == 0) {
            inputfields = create_input_fields(location_add_inputfields)
        } else {
            inputfields = create_input_fields(location_update_inputfields)
        }
        for (var field of inputfields) {
            var group = $('<div>', {
                class: 'form-group'
            })
            group.append(field)
            inputarea.append(group)
        }

        var buttonarea = $('<div>', {
            class: 'buttonarea'
        })
        var cancelButton = $("<input>", {
            type: "button",
            class: "location_upsert_cancel btn btn-danger",
            value: "Отмена"
        })
        var sendButton = $("<input>", {
            type: "submit",
            class: "location_upsert_send btn btn-default",
            value: "Отправить"
        })
        cancelButton.click(function() {
            form.empty()
        })

        buttonarea.append(cancelButton)
        buttonarea.append(sendButton)

        form.append(inputarea)
        form.append(buttonarea)
    }

    var location_upsert_form = $(".location_upsert_form")
    location_upsert_form.submit(function() {
        event.preventDefault()
        var params = location_upsert_form.serialize()
        var url = entrypoints.locations + "/upsert"
        $.ajax({
            url: url + "?" + params + "&appid=" + adminKey,
            success: function(data) {
                logger(data)
            },
            error: function(jqXHR, textStatus, errorThrown) {
                logger("Ошибка запроса к " + url + ":   " + textStatus)
            }
        })
        refresh_location_list()
    })

    upsert_location_button.click(function() {
        event.preventDefault()
        refresh_upsert_form($(".location_upsert_form"), 1)
    })

    $(".location_data_download").click(refresh_location_list_log)

    $("form.weather").submit(function(event) {
        event.preventDefault();
        var locationid = $(location_list_model_id + " option:selected").val()
        var wtype = "current"
        set_spinner_status(get_weatherdata_spinner, STATUS.LOADING)
        $.ajax({
            url: entrypoints.history + "?" + "locationid=" + locationid + "&" + "wtype=" + wtype + "&" + "appid=" + adminKey,
            success: function(data) {
                var jsonData = $.parseJSON(data)
                var status = jsonData['status']
                var message = jsonData['message']
                var content = jsonData['content']
                $(".weathertable").empty()
                $(".weatherchart").empty()
                logger(data)
                if (status != 200) {
                    set_spinner_status(get_weatherdata_spinner, STATUS.ERROR)
                    logger("Request failed with status " + String(status) + " and message: " + message)
                } else {
                    set_spinner_status(get_weatherdata_spinner, STATUS.OK)
                    $(".weathertable").append(build_weather_table(content['history']))
                    $(".weathertable > table").DataTable({
                        "paging": true,
                        "pagingType": "full_numbers"
                    })
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
        var table = $("<table>", {
            class: "table table-striped table-bordered"
        })

        var table_elements = [{
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
        var thead = $("<thead>")
        var theadtr = $("<tr>")
        for (var element of table_elements) {
            theadtr.append("<td>" + element.name + "</td>")
        }
        thead.append(theadtr)
        table.append(thead)

        var content = historyObject['data']
        var tbody = $("<tbody>")
        for (var history_entry of content) {
            if (history_entry['status'] != 200) {
                continue
            }
            var history_entry_row = $("<tr>")

            var history_entry_elements = {
                "json_link": "<a href='" + entrypoints.history + "?" + $.param({
                    entryid: history_entry['objectid']
                }) + "'>" + "Открыть" + "</a>",
                "source": history_entry['source']['name'],
                "dt": new Date(history_entry['measurements'][0]['timestamp'] * 1000).toISOString(),
                "request_dt": new Date(history_entry['request_time'] * 1000).toISOString(),
                "temp": history_entry['measurements'][0]['data']['temp'].toFixed(1),
                "pressure": history_entry['measurements'][0]['data']['pressure'].toFixed(1),
                "humidity": history_entry['measurements'][0]['data']['humidity'].toFixed(1),
                "wind_speed": history_entry['measurements'][0]['data']['wind'].toFixed(1),
                "precipitation": history_entry['measurements'][0]['data']['precipitation'].toFixed(1)
            }
            if (history_entry['url'] != undefined) {
                history_entry_elements["raw_link"] = "<a href='" + history_entry['url'] + "'>Открыть</a>"
            } else {
                history_entry_elements["raw_link"] = "Недоступен"
            }

            for (var row_cell of table_elements) {
                var text = history_entry_elements[row_cell.id]
                history_entry_row.append("<td>" + text + "</td>")
            }

            tbody.append(history_entry_row)
        }
        table.append(tbody)

        return table
    }
});
