"use strict";

var charts, commonstuff, helpers;

charts = require("./charts.js");
commonstuff = require("./commonstuff.js");
helpers = require("./helpers.js");


function main() {
    var APIEP, APIVER, serveraddr, serverEP, entrypoints, adminKey, location_add_inputfields, location_update_inputfields, location_list_model, location_list_model_id, output;

    APIEP = "api";
    APIVER = "0.1";
    serveraddr = "";
    serverEP = serveraddr + "/" + APIEP + "/" + APIVER;
    entrypoints = {
        appid_check: serverEP + "/" + "check_appid",
        locations: serverEP + "/" + "locations",
        history: serverEP + "/" + "history"
    };

    adminKey = "";

    location_add_inputfields = [{
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
    }];

    location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (для редактирования)"
    }].concat(location_add_inputfields);

    var appid_check_spinner = $('#appid_check_spinner');
    var refresh_spinner = $('#refresh_spinner');
    var upsert_location_spinner = $('#upsert_location_spinner');
    var location_data_download_spinner = $('#location_data_download_spinner');
    var get_weatherdata_spinner = $('#get_weatherdata_spinner');

    location_list_model_id = "select#location_list";
    location_list_model = $(location_list_model_id);

    function refresh_location_list() {
        location_list_model.empty();

        output = "";

        helpers.get_with_spinner_and_callback(entrypoints.locations, location_data_download_spinner, function(data) {
            output = data;
            var data_object = $.parseJSON(data);
            var locations = helpers.getlocations(data_object);

            for (var entry of locations) {
                var entryOption = $("<option>", {
                    value: entry.id
                });
                entryOption.append(entry.name);
                location_list_model.append(entryOption);
            }

        })
        return output;
    }

    function refresh_location_list_log() {
        var data = refresh_location_list();
        helpers.logger(data);
    }

    var appid_check_form = $("#appid_check_form");

    var refresh_button = $('#refresh_button');
    var upsert_location_button = $('#upsert_location');

    var admin_buttons = [refresh_button, upsert_location_button];

    function disable_admin_buttons() {
        for (var button of admin_buttons) {
            button.attr("disabled", true);
        };
    };

    function enable_admin_buttons() {
        for (var button of admin_buttons) {
            button.attr("disabled", false);
        };
    };


    function check_appid(appid) {
        var url = entrypoints.appid_check;
        $.ajax({
            url: url + "?appid=" + appid,
            success: function(data) {
                helpers.logger(data)
                var content = $.parseJSON(data)
                if (content.status == 200) {
                    adminKey = appid_check_form.serializeArray()[0].value;
                    helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.OK);
                    enable_admin_buttons();
                } else {
                    helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.ERROR);
                    disable_admin_buttons();
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.ERROR);
                helpers.logger("Ошибка запроса к " + url + ":   " + textStatus);
                disable_admin_buttons();
            }
        });
    };

    function weather_refresh_url(entrypoints, status, locationid, wtype, adminKey) {
        return entrypoints.history + "?" + "status=" + status + "&" + "locationid=" + locationid + "&" + "wtype=" + wtype + "&" + "appid=" + adminKey
    };

    /* Actions on page load */
    refresh_location_list_log();
    disable_admin_buttons();
    check_appid("");
    for (var spinner of[location_data_download_spinner, get_weatherdata_spinner]) {
        helpers.set_spinner_status(spinner, helpers.STATUS.OK);
    };

    /* Events */
    appid_check_form.submit(function() {
        event.preventDefault();
        check_appid(appid_check_form.serializeArray()[0].value);
    });

    refresh_button.click(function() {
        helpers.get_with_spinner_and_callback(entrypoints.history + "/refresh" + "?appid=" + adminKey, refresh_spinner);
    });

    function refresh_upsert_form(form, upsert_type) {
        form.empty();

        var inputarea = $('<div>', {
            class: 'inputarea'
        });
        var inputfields = [];
        if (upsert_type == 0) {
            inputfields = helpers.create_input_fields(location_add_inputfields);
        } else {
            inputfields = helpers.create_input_fields(location_update_inputfields);
        }
        for (var field of inputfields) {
            var group = $('<div>', {
                class: 'form-group'
            });
            group.append(field);
            inputarea.append(group);
        };

        var buttonarea = $('<div>', {
            class: 'buttonarea'
        });
        var cancelButton = $("<input>", {
            type: "button",
            class: "location_upsert_cancel btn btn-danger",
            value: "Отмена"
        });
        var sendButton = $("<input>", {
            type: "submit",
            class: "location_upsert_send btn btn-default",
            value: "Отправить"
        });
        cancelButton.click(function() {
            form.empty();
        });

        buttonarea.append(cancelButton);
        buttonarea.append(sendButton);

        form.append(inputarea);
        form.append(buttonarea);
    };

    var location_upsert_form = $("#location_upsert_form");
    location_upsert_form.submit(function() {
        event.preventDefault()
        var params = location_upsert_form.serialize()
        var url = entrypoints.locations + "/upsert"
        $.ajax({
            url: url + "?" + params + "&appid=" + adminKey,
            success: function(data) {
                helpers.logger(data)
            },
            error: function(jqXHR, textStatus, errorThrown) {
                helpers.logger("Ошибка запроса к " + url + ":   " + textStatus)
            }
        })
        refresh_location_list()
    })

    upsert_location_button.click(function() {
        event.preventDefault()
        refresh_upsert_form(location_upsert_form, 1)
    })

    var location_data_download_button = $("#location_data_download")

    location_data_download_button.click(refresh_location_list_log)

    var weathertable_container = $("#weathertable")
    var weatherchart_container = $("#weatherchart")

    function show_data(data) {
        var jsonData = $.parseJSON(data);
        var status = jsonData.status;
        var message = jsonData.message;
        var content = jsonData.content;
        weathertable_container.empty();
        weathertable_container.empty();
        if (status != 200) {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.ERROR);
            helpers.logger("Request failed with status " + String(status) + " and message: " + message);
        } else {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.OK);
            var history_table_data = commonstuff.make_history_table_data(content.history, entrypoints.history);
            var history_table_values = history_table_data[0];
            var history_table_columns = history_table_data[1];
            var table = $("<table>", {
                class: "table table-striped"
            });
            weathertable_container.append(table);
            table.DataTable({
                data: history_table_values,
                columns: history_table_columns,
                paging: true,
                pagingType: "full_numbers"
            });
            charts.build_weather_chart(weatherchart_container, content.history);
        };
    };

    $("form#weather_request").submit(function(event) {
        event.preventDefault();
        var locationid = $(location_list_model_id + " option:selected").val();
        var wtype = "current";
        var download_url = weather_refresh_url(entrypoints, "200", locationid, wtype, adminKey);

        helpers.get_with_spinner_and_callback(download_url, get_weatherdata_spinner, show_data);
    });

};

$(document).ready(main);
