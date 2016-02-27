"use strict";
var charts = require("./charts.js");
var commonstuff = require("./commonstuff.js");
var helpers = require("./helpers.js");
var settings = require("./settings.js");

function makeLanding() {
    return helpers.quickParseHTML(`
    <div class="jumbotron">
      <h1>Weather info at your hand</h1>
      <p>Weather Checker lets you compare weather data from different providers.</p>
      <p>
        <a class="btn btn-primary btn-lg" role="button" disabled="true">Select the city in the box above to start</a>
      </p>
    </div>
    `)
};

function makeNoData() {
    return helpers.quickParseHTML(`
    <div class="jumbotron">
      <h1>No data available</h1>
      <p>Try expanding the time interval or choose another location.</p>
    </div>
    `)
};

function makeNotFound() {
    return helpers.quickParseHTML(`
    <div class="jumbotron">
      <h1>Location not found.</h1>
      <p>The specified location does not exist.</p>
    </div>
    `)
};

function makeShim() {
    var shim_spinner = document.createElement('span');
    shim_spinner.setAttribute('class', "fa fa-refresh fa-spin shim-spinner");

    var shim_spinner_p = document.createElement('p');
    shim_spinner_p.appendChild(shim_spinner);

    var shim_spinner_row = document.createElement('div');
    shim_spinner_row.setAttribute('class', "col-lg-offset-5");
    shim_spinner_row.appendChild(shim_spinner_p);

    var shim_spinner_container = document.createElement('div');
    shim_spinner_container.setAttribute('class', "container");
    shim_spinner_container.appendChild(shim_spinner_row);

    return shim_spinner_container;
}

function makeContactLink(emailaddr) {
    var contactLink = document.createElement('a')
    contactLink.setAttribute('href', "mailto:" + emailaddr)
    contactLink.innerText = 'Contact us.'

    return contactLink;
}

function makeLocationUrl(locationid) {
    return "/location" + "/" + locationid
}

function makeLocationSelectEntry() {
    var entry = document.createElement('option');
    entry.setAttribute('disabled', true);
    entry.setAttribute('objectid', "");
    entry.setAttribute('slug', "");
    entry.innerText = "(select your city)";

    return entry
}

function main() {
    var entrypoints = settings.entrypoints;

    var adminKey = "";

    var datepickers = document.getElementsByClassName('daterange');
    var activeZone = document.getElementsByClassName('activezone');

    var landing_container = document.getElementById('landing');

    var nodata_container = document.getElementById('nodata');

    var loading_shim_container = document.getElementById('loading_shim');

    var get_weatherdata_spinner = document.getElementById('get_weatherdata_spinner');

    var location_list_model = document.getElementById('location_list');

    var weathertable_container = document.getElementById("weathertable");
    var weatherchart_container = document.getElementById("weatherchart");
    var weather_request_form = document.getElementById("weather_request_form");

    var request_range_picker = document.getElementById('request_range_picker');
    var request_range_span = document.getElementById('request_range_span');
    var request_location_select = location_list_model;

    var contactInfoContainer = document.getElementById('contact-info');

    function refresh_location_list_simple(cb) {
        commonstuff.refresh_location_list(location_list_model, entrypoints, null, cb);
    };

    function weather_refresh_url(entrypoints, status, locationid, wtype,
        request_start, request_end, adminKey) {
        var url = entrypoints.history + "?" + "status=" + status + "&" +
            "locationid=" + locationid + "&" + "wtype=" + wtype + "&" +
            "appid=" + adminKey
        if (request_start != null) {
            url = url + '&requeststart=' + String(request_start)
        }
        if (request_end != null) {
            url = url + '&requestend=' + String(request_end)
        }
        return url
    };

    /* Model init */
    jQuery(datepickers).daterangepicker({
        timePicker: true,
        timePicker24Hour: true,
        locale: {
            format: "YYYY-MM-DD HH:mm"
        }
    });

    /* Events */
    function empty_body() {
        for (var node of Array.from(activeZone)) {
            helpers.clearChildren(node);
        };
    };

    function show_landing() {
        empty_body();
        var landingBody = makeLanding();
        landing_container.appendChild(landingBody);
    };

    function show_nodata() {
        empty_body();
        var noDataBody = makeNoData();
        nodata_container.appendChild(noDataBody);
    };

    function showNotFound() {
        empty_body();
        var notFoundBody = makeNotFound();
        nodata_container.appendChild(notFoundBody);
    };

    function show_shim() {
        empty_body();

        loading_shim_container.appendChild(makeShim());
    };

    function show_data(data) {
        var jsonData = JSON.parse(data);
        var status = jsonData.status;
        var message = jsonData.message;
        var content = jsonData['content'];
        empty_body();
        if (status != 200) {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.ERROR);
            helpers.logger("Request failed with status " + String(status) +
                " and message: " + message);
        } else {
            helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.OK);
            if (content['history']['data'].length > 0) {
                var history_table_data = commonstuff.make_history_table_data(
                    content['history'], entrypoints['history']);
                var history_table_values = history_table_data[0];
                var history_table_columns = history_table_data[1];
                var history_table_opts = history_table_data[2];

                var table = document.createElement("table");
                table.setAttribute("class", "table table-striped");

                weathertable_container.appendChild(table);
                var tableInitOpts = {
                    data: history_table_values,
                    columns: history_table_columns,
                    paging: true,
                    pagingType: "full_numbers"
                };
                Object.assign(tableInitOpts, history_table_opts);

                jQuery(table).DataTable(tableInitOpts);
                charts.build_weather_chart(weatherchart_container, content['history']['data']);
            } else {
                show_nodata();
            }
        };
    };

    function download_weather_data() {
        var locoption = location_list_model.options[location_list_model.selectedIndex];
        var locationid = locoption.getAttribute("objectid");
        var locationslug = locoption.getAttribute("slug");
        var locationName = locoption.innerText;
        var wtype = "current";

        var request_start_momentObject = jQuery(request_range_picker).data("daterangepicker").startDate;
        var request_end_momentObject = jQuery(request_range_picker).data("daterangepicker").endDate;

        history.replaceState({}, "", makeLocationUrl(locationslug));
        document.title = "Weather in " + locationName;

        request_range_span.innerHTML = request_start_momentObject.format('D MMMM YYYY HH:mm') + ' - ' + request_end_momentObject.format('D MMMM YYYY HH:mm');

        var request_start = request_start_momentObject.unix();
        var request_end = request_end_momentObject.unix();

        if (locationid != "") {
            var download_url = weather_refresh_url(entrypoints, "200", locationid, wtype, request_start, request_end, adminKey);
            show_shim();
            helpers.get_with_spinner_and_callback(download_url, get_weatherdata_spinner, show_data);
        };
    };

    function weather_request_form_submit() {
        download_weather_data();
    }

    weather_request_form.onsubmit = function (event) {
        event.preventDefault();
        weather_request_form_submit();
    };

    function make_request_range_picker_span() {
        var request_start_momentObject = jQuery(request_range_picker).data("daterangepicker").startDate;
        var request_end_momentObject = jQuery(request_range_picker).data("daterangepicker").endDate;
        request_range_span.innerHTML = request_start_momentObject.format('D MMMM YYYY HH:mm') + ' - ' + request_end_momentObject.format('D MMMM YYYY HH:mm');
    };

    request_location_select.onchange = function () {
        weather_request_form_submit();
    };

    jQuery(request_range_picker).on("apply.daterangepicker", function () {
        make_request_range_picker_span();

        weather_request_form_submit();
    });


    /* Actions on page load */
    var start_time = moment().subtract(3, 'days');
    var end_time = moment();
    jQuery.ajax({
        url: entrypoints.settingsData,
        contentType: "text/plain",
        success: function (settingsObject) {
            var settingsMap = settingsObject["content"]["settings"];
            var minstart = moment().subtract(settingsMap["max-depth"], 'hours');
            jQuery(request_range_picker).data('daterangepicker').minDate = minstart;
            contactInfoContainer.appendChild(makeContactLink(settingsMap["email"]));
        }
    });

    jQuery(request_range_picker).data('daterangepicker').setStartDate(start_time);
    jQuery(request_range_picker).data('daterangepicker').setEndDate(end_time);
    make_request_range_picker_span();

    helpers.set_spinner_status(get_weatherdata_spinner, helpers.STATUS.HAND_LEFT);

    show_landing();

    var pathName = decodeURI(window.location.pathname);

    var pathNameSplit = pathName.split('/').slice(1);

    var prefix = pathNameSplit[0];
    var value = pathNameSplit[1];

    var refresh_cb_base = function () {
        var selectCityEntry = makeLocationSelectEntry();
        location_list_model.insertBefore(selectCityEntry, null);
        location_list_model.selectedIndex = Array.from(location_list_model.childNodes).indexOf(selectCityEntry);
    };

    var refresh_cb = refresh_cb_base;

    if (prefix == "location") {
        var preselectedSlug = value;

        refresh_cb = function () {
            refresh_cb_base();
            var foundIndex = helpers.fieldValueInSelect(location_list_model, "slug", preselectedSlug);

            if (foundIndex != -1) {
                location_list_model.selectedIndex = foundIndex;
                request_location_select.onchange();
            } else {
                showNotFound();
            };
        };
    };

    refresh_location_list_simple(refresh_cb);
};

document.addEventListener('DOMContentLoaded', main);
