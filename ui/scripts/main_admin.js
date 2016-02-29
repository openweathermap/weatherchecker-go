"use strict";
var helpers = require("./helpers.js");
var settings = require("./settings.js");

function main() {
    var entrypoints = settings.entrypoints;
    var adminKey = "";

    var location_add_inputfields = [{
        Name: "city_name",
        Default: "",
        Placeholder: "Location name"
    }, {
        Name: "iso_country",
        Default: "",
        Placeholder: "Country ISO code"
    }, {
        Name: "country_name",
        Default: "",
        Placeholder: "Country name"
    }, {
        Name: "latitude",
        Default: "",
        Placeholder: "Latitude"
    }, {
        Name: "longitude",
        Default: "",
        Placeholder: "Longitude"
    }];

    var location_update_inputfields = [{
        Name: "entryid",
        Default: "",
        Placeholder: "ObjectID (for editing existing location)"
    }].concat(location_add_inputfields);

    var appid_check_spinner = document.getElementById('appid_check_spinner');
    var refresh_spinner = document.getElementById('refresh_spinner');
    var upsert_location_spinner = document.getElementById('upsert_location_spinner');

    var location_upsert_form = document.getElementById("location_upsert_form");

    var appid_check_form = document.getElementById("appid_check_form");

    var refresh_button = document.getElementById('refresh_button');
    var upsert_location_button = document.getElementById('upsert_location');

    var admin_buttons = [refresh_button, upsert_location_button];

    function disable_admin_buttons() {
        admin_buttons.forEach(function (button) {
            button.setAttribute("disabled", "");
        });
    };

    function enable_admin_buttons() {
        admin_buttons.forEach(function (button) {
            button.removeAttribute("disabled");
        });
    };

    function setAdminKey(value) {
        adminKey = value;
    };

    function getAdminKey() {
        return adminKey;
    };

    function check_appid(appid) {
        var url = entrypoints.appid_check;
        jQuery.ajax({
            url: url + "?appid=" + appid,
            success: function (data) {
                if (data.status == 200) {
                    setAdminKey(appid);
                    helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.OK);
                    enable_admin_buttons();
                } else {
                    helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.ERROR);
                    disable_admin_buttons();
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                helpers.set_spinner_status(appid_check_spinner, helpers.STATUS.ERROR);
                disable_admin_buttons();
            }
        });
    };

    /* Events */
    appid_check_form.onsubmit = function (event) {
        event.preventDefault();
        check_appid($(appid_check_form).serializeArray()[0].value);
    };

    refresh_button.onclick = function (event) {
        helpers.get_with_spinner_and_callback(entrypoints.history + "/refresh" + "?appid=" + getAdminKey(), refresh_spinner);
    };

    function refresh_upsert_form(form, upsert_type) {
        helpers.clearChildren(form);

        var inputarea = document.createElement("div");
        inputarea.setAttribute("class", "inputarea");

        var inputfields = [];
        if (upsert_type == 0) {
            inputfields = helpers.create_input_fields(location_add_inputfields);
        } else {
            inputfields = helpers.create_input_fields(location_update_inputfields);
        }
        inputfields.forEach(function (field) {
            var group = document.createElement("div");
            group.setAttribute("class", "form-group");

            group.appendChild(field);
            this.appendChild(group);
        }, inputarea);

        var buttonarea = document.createElement('div')
        buttonarea.setAttribute("class", "buttonarea")

        var cancelButton = document.createElement("input")
        cancelButton.setAttribute("type", "button");
        cancelButton.setAttribute("class", "location_upsert_cancel btn btn-danger");
        cancelButton.setAttribute("value", "Cancel");

        var sendButton = document.createElement("input")
        sendButton.setAttribute("type", "submit");
        sendButton.setAttribute("class", "location_upsert_send btn btn-default");
        sendButton.setAttribute("value", "Send");

        cancelButton.onclick = function () {
            helpers.clearChildren(form);
        };

        buttonarea.appendChild(cancelButton);
        buttonarea.appendChild(sendButton);

        form.appendChild(inputarea);
        form.appendChild(buttonarea);
    };

    function closeUpsertForm() {
        helpers.clearChildren(location_upsert_form);
    };

    function upsert_location() {
        var params = jQuery(location_upsert_form).serialize();
        var url = entrypoints.locations + "/upsert" + "?" + params + "&appid=" + getAdminKey();
        helpers.get_with_spinner_and_callback(url, upsert_location_spinner, function (data) {
            if (JSON.parse(data).status == 200) {
                closeUpsertForm();
            };
        });
    };

    location_upsert_form.onsubmit = function (event) {
        event.preventDefault();
        upsert_location();
    };

    upsert_location_button.onclick = function (event) {
        event.preventDefault();
        refresh_upsert_form(location_upsert_form, 1);
    };

    /* Actions on page load */
    disable_admin_buttons();
};

document.addEventListener('DOMContentLoaded', main);
