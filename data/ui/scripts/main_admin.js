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
    }];

    var location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (для редактирования)"
    }].concat(location_add_inputfields);

    var appid_check_spinner = $('#appid_check_spinner');
    var refresh_spinner = $('#refresh_spinner');
    var upsert_location_spinner = $('#upsert_location_spinner');
    var get_weatherdata_spinner = $('#get_weatherdata_spinner');

    var location_list_model_id = "select#location_list";
    var location_list_model = $(location_list_model_id);

    var location_upsert_form = $("#location_upsert_form");
    var weathertable_container = $("#weathertable");
    var weatherchart_container = $("#weatherchart");
    var weather_request_form = $("form#weather_request");

    var appid_check_form = $("#appid_check_form");

    var datepickers = $(".datepicker");

    var refresh_button = $('#refresh_button');
    var request_start_picker = $('#request_start_picker');
    var request_end_picker = $('#request_end_picker');
    var request_location_select = location_list_model;
    var upsert_location_button = $('#upsert_location');

    var admin_buttons = [refresh_button, upsert_location_button];


    function refresh_location_list_log() {
        commonstuff.refresh_location_list(location_list_model, entrypoints, null, helpers.logger);
    }



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
                helpers.logger(data);
                if (data.status == 200) {
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

    function weather_refresh_url(entrypoints, status, locationid, wtype, request_start, request_end, adminKey) {
        var url = entrypoints.history + "?" + "status=" + status + "&" + "locationid=" + locationid + "&" + "wtype=" + wtype + "&" + "appid=" + adminKey
        if (request_start != null) {
            url = url + '&requeststart=' + String(request_start)
        }
        if (request_end != null) {
            url = url + '&requestend=' + String(request_end)
        }
        return url
    };

    /* Model init */
    datepickers.datetimepicker();

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

    function upsert_location() {
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
        commonstuff.refresh_location_list(location_list_model, entrypoints, null, helpers.logger)
    }

    location_upsert_form.submit(function(event) {
        event.preventDefault()
        upsert_location()

    })

    upsert_location_button.click(function(event) {
        event.preventDefault()
        refresh_upsert_form(location_upsert_form, 1)
    })

    /* Actions on page load */
    refresh_location_list_log();
    disable_admin_buttons();
    check_appid("");
};

$(document).ready(main);
