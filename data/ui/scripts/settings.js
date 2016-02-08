"use strict";

var APIEP = "api";
var APIVER = "0.1";
var serveraddr = "";
var serverEP = serveraddr + "/" + APIEP + "/" + APIVER;

exports.entrypoints = {
    appid_check: serverEP + "/" + "check_appid",
    settingsData: serverEP + "/" + "settings",
    locations: serverEP + "/" + "locations",
    history: serverEP + "/" + "history"
};

exports.testing = {
    appid: '',
    location: {
        latitude: '',
        longitude: ''
    }
};
