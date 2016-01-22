"use strict";

var STATUS = {
    OK: 0,
    LOADING: 1,
    ERROR: 2
};

exports.STATUS = STATUS;
exports.getlocations = getlocations;
exports.create_input_fields = create_input_fields;
exports.set_spinner_status = set_spinner_status;
exports.get_with_spinner_and_callback = get_with_spinner_and_callback;
exports.logger = logger;

function logger(data) {
    console.log(data);
};

function getlocations(data) {
    var i, locations, location_list, location_entry, entry;

    locations = [];
    location_list = data.content.locations;

    for (i in location_list) {
        if (location_list.hasOwnProperty(i)) {
            location_entry = location_list[i];
            entry = {};
            entry.id = location_entry.objectid;
            entry.name = location_entry.city_name;
            locations.push(entry);
        };
    };

    return locations
};

function create_input_fields(inputfields_desc) {
    var input_fields = new Array;
    for (var inputfield of inputfields_desc) {
        var entry = $('<input>', {
            name: inputfield.Name,
            type: "text",
            class: inputfield.Name + " form-control",
            value: inputfield.Default,
            placeholder: inputfield.Placeholder
        });
        input_fields.push(entry);
    };

    return input_fields
};

function set_spinner_status(spinnerContainer, status) {
    spinnerContainer.empty();
    var iconClass = "";
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
    };
    spinnerContainer.append($('<span>', {
        class: iconClass
    }));
};

function get_with_spinner_and_callback(requestUrl, spinnerObject, callbackFunc) {
    set_spinner_status(spinnerObject, STATUS.LOADING);
    $.ajax({
        url: requestUrl,
        success: function(data) {
            var jsonData = $.parseJSON(data);
            var status = jsonData.status;

            if (status == 200) {
                set_spinner_status(spinnerObject, STATUS.OK);
            } else {
                set_spinner_status(spinnerObject, STATUS.ERROR);
            }

            if (callbackFunc != undefined) {
                callbackFunc(data);
            }

            logger(data);
        },
        error: function(data) {
            set_spinner_status(spinnerObject, STATUS.ERROR);
        }
    });
};
