"use strict"
var charts = require("./charts.js");
var commonstuff = require("./commonstuff.js");
var helpers = require("./helpers.js");
var settings = require("./settings.js");

var testingData = settings.testing

function main() {
    var entrypoints = settings.entrypoints;

    var dt_current_ms = Date.now();
    var dt_3days = 3 /*d*/ * 24 /*h*/ * 60 /*m*/ * 60 /*s*/ ;


    var requestDeviation = 2 /*h*/ * 60 /*m*/ * 60 /*s*/ ;


    var get_weatherdata = document.getElementById("get_weatherdata");
    var location_list = document.getElementById("location_list");
    var forecastchart_container = document.getElementById("forecast_chart");

    var location_download_button = document.getElementById('location_download_button')
    var location_download_spinner = document.getElementById('location_download_spinner')

    var location_list_model = document.getElementById("location_list");

    var get_weather_form = document.getElementById('get_weather_form');

    var locationMap = new Map;

    var owm_data = [];

    // Funcs
    function refresh_location_list_log() {
        commonstuff.refresh_location_list(location_list_model, entrypoints, location_download_spinner, function (locMap) {
            locationMap = locMap;
        });
    };

    function get_weather() {
        var locationOption = location_list_model.options[location_list_model.selectedIndex];
        var locationId = locationOption.getAttribute("value");
        var locationObject = locationMap[locationId];

        var historyRequestParams = {
            lat: locationObject.latitude,
            lon: locationObject.longitude,
            start: requestStart,
            end: requestEnd,
            appid: testingData.appid
        };

        var requestEnd = Math.floor(dt_current_ms / 1000);
        var requestStart = requestEnd - dt_3days;

        var historyRequestUrl = "http://history.openweathermap.org/data/2.5/history/city" + "?" + jQuery.param(historyRequestParams);

        helpers.get_with_spinner_and_callback(historyRequestUrl, null, function (data) {
            var jsonData = JSON.parse(data)
            var checkerHistory = commonstuff.parseOWMhistory(jsonData)
            charts.build_weather_chart(forecastchart_container, checkerHistory)
        });
    };

    // Events
    location_download_button.onclick = refresh_location_list_log;
    get_weather_form.onsubmit = function (event) {
        event.preventDefault();
        get_weather();
    };

};

document.addEventListener('DOMContentLoaded', main);
