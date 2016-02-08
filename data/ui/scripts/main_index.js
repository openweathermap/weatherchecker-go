"use strict";
var charts = require("./charts.js");
var commonstuff = require("./commonstuff.js");
var helpers = require("./helpers.js");
var settings = require("./settings.js");

function makeLanding() {
    var landingBody = '<div class="jumbotron"><h1>Weather info at your hand</h1><p>OWM Weather Checker lets you compare weather data from different providers.</p><p><a class="btn btn-primary btn-lg" role="button" disabled=true>Select the city in the box above to start</a></p></div>'

    return $(landingBody)
}

function main() {
    var entrypoints = settings.entrypoints;

    var adminKey = "";

    var datepickers = $('.daterange');
    var activeZone = $('.activezone');

    var landing_container = $('#landing');

    var loading_shim_container = $('#loading_shim');

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

    var refresh_button = $('#refresh_button');
    var request_range_picker = $('#request_range_picker');
    var request_range_span = $('#request_range_span')
    var request_location_select = location_list_model;
    var upsert_location_button = $('#upsert_location');

    var admin_buttons = [refresh_button, upsert_location_button];

    function refresh_location_list_log() {
        commonstuff.refresh_location_list(location_list_model, entrypoints, null, helpers.logger);
    }

    function refresh_location_list_nolog() {
        commonstuff.refresh_location_list(location_list_model, entrypoints, null);
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
    datepickers.daterangepicker({
        timePicker: true,
        timePicker24Hour: true,
        locale: {
            format: "YYYY-MM-DD HH:mm"
        }
    });

    /* Events */
    function empty_body() {
        activeZone.empty();
    };

    function show_landing() {
        empty_body();
        var landingBody = makeLanding();
        landing_container.append(landingBody);
    };

    function show_shim() {
        empty_body();

        var shim_spinner = $('<span>', {
            class: "fa fa-refresh fa-spin shim-spinner"
        })

        var shim_spinner_p = $('<p>')
        shim_spinner_p.append(shim_spinner)

        var shim_spinner_row = $('<div>', {
            class: "col-lg-offset-5"
        })
        shim_spinner_row.append(shim_spinner_p)

        var shim_spinner_container = $('<div>', {
            class: "container"
        })
        shim_spinner_container.append(shim_spinner_row)

        loading_shim_container.append(shim_spinner_container);
    };

    function show_data(data) {
        var jsonData = $.parseJSON(data);
        var status = jsonData.status;
        var message = jsonData.message;
        var content = jsonData['content'];
        empty_body();
        if (status != 200) {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.ERROR);
            helpers.logger("Request failed with status " + String(status) + " and message: " + message);
        } else {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.OK);
            var history_table_data = commonstuff.make_history_table_data(content['history'], entrypoints['history']);
            var history_table_values = history_table_data[0];
            var history_table_columns = history_table_data[1];
            var history_table_opts = history_table_data[2];
            var table = $("<table>", {
                class: "table table-striped"
            });
            weathertable_container.append(table);
            var tableInitOpts = {
                data: history_table_values,
                columns: history_table_columns,
                paging: true,
                pagingType: "full_numbers"
            };
            Object.assign(tableInitOpts, history_table_opts);

            table.DataTable(tableInitOpts);
            charts.build_weather_chart(weatherchart_container, content['history']['data']);
        };
    };

    function download_weather_data() {
        var locationid = $(location_list_model_id + " option:selected").val();
        var wtype = "current";

        var request_start_momentObject = request_range_picker.data("daterangepicker").startDate;
        var request_end_momentObject = request_range_picker.data("daterangepicker").endDate;
        request_range_span.html(request_start_momentObject.format('D MMMM YYYY HH:mm') + ' - ' + request_end_momentObject.format('D MMMM YYYY HH:mm'));

        var request_start = request_start_momentObject.unix();
        var request_end = request_end_momentObject.unix();

        if (locationid != "") {
            var download_url = weather_refresh_url(entrypoints, "200", locationid, wtype, request_start, request_end, adminKey);
            show_shim();
            helpers.get_with_spinner_and_callback(download_url, get_weatherdata_spinner, show_data);
        };
    };

    weather_request_form.submit(function(event) {
        event.preventDefault();
        download_weather_data();
    });

    function make_request_range_picker_span() {
        var request_start_momentObject = request_range_picker.data("daterangepicker").startDate;
        var request_end_momentObject = request_range_picker.data("daterangepicker").endDate;
        request_range_span.html(request_start_momentObject.format('D MMMM YYYY HH:mm') + ' - ' + request_end_momentObject.format('D MMMM YYYY HH:mm'));
    };

    for (var formObject of[request_location_select]) {
        formObject.on("change", function(event) {
            weather_request_form.submit();
        });
    };


    request_range_picker.on("apply.daterangepicker", function(event) {
        make_request_range_picker_span();

        weather_request_form.submit();
    });


    /* Actions on page load */
    var start_time = moment().subtract(3, 'days')
    var end_time = moment();
    $.ajax({
        url: entrypoints.settingsData,
        contentType: "text/plain",
        success: function(settingsData) {
            var minstart = moment().subtract(settingsData["content"]["settings"]["max-depth"], 'hours');
            request_range_picker.data('daterangepicker').minDate = minstart;
        }
    });

    request_range_picker.data('daterangepicker').setStartDate(start_time);
    request_range_picker.data('daterangepicker').setEndDate(end_time);
    make_request_range_picker_span();

    refresh_location_list_nolog();
    helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.HAND_LEFT)

    var selectCityEntry = $('<option>', {
        disabled: true,
        selected: true,
        value: ""
    });
    selectCityEntry.append("(select your city)");

    location_list_model.prepend(selectCityEntry);

    show_landing();
};

$(document).ready(main);
