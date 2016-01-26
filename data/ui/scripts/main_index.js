"use strict";
var charts = require("./charts.js");
var commonstuff = require("./commonstuff.js");
var helpers = require("./helpers.js");
var settings = require("./settings.js");

function main() {
    var entrypoints = settings.entrypoints;

    var adminKey = "";

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
    }];

    var location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (для редактирования)"
    }].concat(location_add_inputfields);

    var appid_check_spinner = $('#appid_check_spinner');
    var refresh_spinner = $('#refresh_spinner');
    var upsert_location_spinner = $('#upsert_location_spinner');
    var location_data_download_spinner = $('#location_data_download_spinner');
    var get_weatherdata_spinner = $('#get_weatherdata_spinner');

    var location_list_model_id = "select#location_list";
    var location_list_model = $(location_list_model_id);

    var location_upsert_form = $("#location_upsert_form");
    var location_data_download_button = $("#location_data_download");
    var weathertable_container = $("#weathertable");
    var weatherchart_container = $("#weatherchart");
    var weather_request_form = $("form#weather_request");

    var appid_check_form = $("#appid_check_form");

    var refresh_button = $('#refresh_button');
    var upsert_location_button = $('#upsert_location');

    function refresh_location_list_log() {
        commonstuff.refresh_location_list(location_list_model, entrypoints, location_data_download_spinner, helpers.logger);
    }


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
            success: function (data) {
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
            error: function (jqXHR, textStatus, errorThrown) {
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
    appid_check_form.submit(function () {
        event.preventDefault();
        check_appid(appid_check_form.serializeArray()[0].value);
    });

    refresh_button.click(function () {
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
        cancelButton.click(function () {
            form.empty();
        });

        buttonarea.append(cancelButton);
        buttonarea.append(sendButton);

        form.append(inputarea);
        form.append(buttonarea);
    };

    location_upsert_form.submit(function () {
        event.preventDefault()
        var params = location_upsert_form.serialize()
        var url = entrypoints.locations + "/upsert"
        $.ajax({
            url: url + "?" + params + "&appid=" + adminKey,
            success: function (data) {
                helpers.logger(data)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                helpers.logger("Ошибка запроса к " + url + ":   " + textStatus)
            }
        })
        commonstuff.refresh_location_list(location_list_model, entrypoints, location_data_download_spinner, helpers.logger)
    })

    upsert_location_button.click(function () {
        event.preventDefault()
        refresh_upsert_form(location_upsert_form, 1)
    })


    location_data_download_button.click(refresh_location_list_log)

    function show_data(data) {
        var jsonData = $.parseJSON(data);
        var status = jsonData.status;
        var message = jsonData.message;
        var content = jsonData['content'];
        weathertable_container.empty();
        weathertable_container.empty();
        if (status != 200) {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.ERROR);
            helpers.logger("Request failed with status " + String(status) + " and message: " + message);
        } else {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.OK);
            var history_table_data = commonstuff.make_history_table_data(content['history'], entrypoints['history']);
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
            charts.build_weather_chart(weatherchart_container, content['history']['data']);
        };
    };

    function download_weather_data() {
        var locationid = $(location_list_model_id + " option:selected").val();
        var wtype = "current";
        var download_url = weather_refresh_url(entrypoints, "200", locationid, wtype, adminKey);

        helpers.get_with_spinner_and_callback(download_url, get_weatherdata_spinner, show_data);
    };

    weather_request_form.submit(function (event) {
        event.preventDefault();
        download_weather_data();
    });

};

$(document).ready(main);
