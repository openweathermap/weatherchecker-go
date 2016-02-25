"use strict";

var STATUS = {
    OK: 0,
    LOADING: 1,
    ERROR: 2,
    HAND_LEFT: 3
};

exports.STATUS = STATUS;
exports.find_closest = find_closest;
exports.collectionToMap = collectionToMap;
exports.create_input_fields = create_input_fields;
exports.set_spinner_status = set_spinner_status;
exports.get_with_spinner_and_callback = get_with_spinner_and_callback;
exports.logger = logger;
exports.valueInSelect = valueInSelect;

function logger(data) {
    console.log(data);
};

function find_closest(x, range) {
    var closest = range[0];
    for (var value of range) {
        if (Math.abs(value - x) < Math.abs(closest - x)) {
            closest = value
        };
    };
    var closestIndex = range.indexOf(closest);
    var closestInfo = {
        Closest: closest,
        Index: closestIndex
    };
    return closestInfo;
};

function collectionToMap(objectCollection, keyName) {
    var newMap = new Map;

    for (var entry of objectCollection) {
        var key = entry[keyName];
        if (key != undefined) {
            newMap[key] = entry;
        };
    };

    return newMap;
}

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
    if (spinnerContainer == null) {
        return;
    };

    spinnerContainer.empty();
    var iconClass = "";
    switch (status) {
    case STATUS.OK: // OK
        iconClass = "fa fa-check";
        break;
    case STATUS.LOADING: // Loading
        iconClass = "fa fa-refresh fa-spin";
        break;
    case STATUS.ERROR: // Error
        iconClass = "fa fa-times";
        break;
    case STATUS.HAND_LEFT: // Select
        iconClass = "fa fa-hand-o-left";
        break;
    default:
        return;
    };
    spinnerContainer.append($('<span>', {
        class: iconClass
    }));
};

function get_with_spinner_and_callback(requestUrl, spinnerContainer, callbackFunc) {
    set_spinner_status(spinnerContainer, STATUS.LOADING);
    $.ajax({
        url: requestUrl,
        dataType: "text",
        success: function (data) {
            var jsonData = $.parseJSON(data);
            var status = jsonData.status;

            if (status == 200) {
                set_spinner_status(spinnerContainer, STATUS.OK);
            } else {
                set_spinner_status(spinnerContainer, STATUS.ERROR);
            }

            if (callbackFunc != undefined) {
                callbackFunc(data);
            }
        },
        error: function (data) {
            set_spinner_status(spinnerContainer, STATUS.ERROR);
        }
    });
};

function valueInSelect(selectModel, value) {
    var exists = false;
    console.log("Searching for " + value + " in the model")
    selectModel.contents().each(function () {
        if (this.value == value) {
            exists = true;
            console.log("Found!")
            return false;
        };
    });

    return exists;
};