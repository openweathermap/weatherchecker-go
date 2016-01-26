"use strict";

var STATUS = {
    OK: 0,
    LOADING: 1,
    ERROR: 2
};

exports.STATUS = STATUS;
exports.create_input_fields = create_input_fields;
exports.set_spinner_status = set_spinner_status;
exports.get_with_spinner_and_callback = get_with_spinner_and_callback;
exports.logger = logger;

function logger(data) {
    console.log(data);
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
        iconClass = "fa fa-spin fa-refresh";
        break;
    case STATUS.ERROR: // Error
        iconClass = "fa fa-minus-circle";
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

            logger(data);
        },
        error: function (data) {
            set_spinner_status(spinnerContainer, STATUS.ERROR);
        }
    });
};
