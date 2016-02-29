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
exports.selectOption = selectOption;
exports.clearChildren = clearChildren;
exports.quickParseHTML = quickParseHTML;

function logger(data) {
    console.log(data);
};

function quickParseHTML(sourceText) {
    return (new DOMParser()).parseFromString(sourceText, "text/html").firstChild;
}

function clearChildren(node) {
    var nodeChildren = Array.from(node.children);

    nodeChildren.forEach(function (child) {
        node.removeChild(child);
    });
};

function find_closest(x, range) {
    var closestIndex = 0;
    var closest = range[closestIndex];

    range.forEach(function (value, i) {
        if (Math.abs(value - x) < Math.abs(closest - x)) {
            closestIndex = i;
            closest = value;
        };
    });

    var closestInfo = {
        Closest: closest,
        Index: closestIndex
    };
    return closestInfo;
};

function collectionToMap(objectCollection, keyName) {
    var newMap = new Map;

    objectCollection.forEach(function (entry) {
        var key = entry[keyName];
        if (key != undefined) {
            newMap[key] = entry;
        };
    });

    return newMap;
}

function create_input_fields(inputfields_desc) {
    var input_fields = new Array;
    inputfields_desc.forEach(function (inputfield) {
        var entry = document.createElement('input')
        entry.setAttribute("name", inputfield.Name)
        entry.setAttribute("type", "text")
        entry.setAttribute("class", inputfield.Name + " form-control")
        entry.setAttribute("value", inputfield.Default)
        entry.setAttribute("placeholder", inputfield.Placeholder)

        input_fields.push(entry);
    });

    return input_fields
};

function makeSpinner(status) {
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
    var newSpinner = document.createElement('span');
    newSpinner.setAttribute('class', iconClass);

    return newSpinner;
}

function set_spinner_status(spinnerContainer, status) {
    if (spinnerContainer == null) {
        return;
    };

    clearChildren(spinnerContainer);

    spinnerContainer.appendChild(makeSpinner(status));
};

function get_with_spinner_and_callback(requestUrl, spinnerContainer, callbackFunc) {
    set_spinner_status(spinnerContainer, STATUS.LOADING);
    jQuery.ajax({
        url: requestUrl,
        dataType: "text",
        success: function (data) {
            var jsonData = JSON.parse(data);
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

function valueInSelect(selectModel, field, value) {
    var existIndex = -1;
    var optionList = Array.from(selectModel.childNodes);

    optionList.forEach(function (option, optionIndex) {
        var optionValue = option.getAttribute(field);
        if (optionValue == value && existIndex < 0) {
            existIndex = optionIndex;
        };
    });

    return existIndex;
};

function selectOption(selectModel, field, value) {
    var foundIndex = valueInSelect(selectModel, field, value);

    var exists = (foundIndex != -1);
    if (exists) {
        selectModel.selectedIndex = foundIndex;
    };

    return exists;
};
