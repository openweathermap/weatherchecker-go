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
exports.polyfill = polyfill;

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

function polyfillObjectAssign() {
    if (!Object.assign) {
        Object.assign = function (target) {
            'use strict';
            if (target === undefined || target === null) {
                throw new TypeError('Cannot convert undefined or null to object');
            }

            var output = Object(target);
            for (var index = 1; index < arguments.length; index++) {
                var source = arguments[index];
                if (source !== undefined && source !== null) {
                    for (var nextKey in source) {
                        if (source.hasOwnProperty(nextKey)) {
                            output[nextKey] = source[nextKey];
                        }
                    }
                }
            }
            return output;
        };
    };
};

function polyfillArrayFrom() {
    // Production steps of ECMA-262, Edition 6, 22.1.2.1
    // Reference: https://people.mozilla.org/~jorendorff/es6-draft.html#sec-array.from
    if (!Array.from) {
        Array.from = (function () {
            var toStr = Object.prototype.toString;
            var isCallable = function (fn) {
                return typeof fn === 'function' || toStr.call(fn) === '[object Function]';
            };
            var toInteger = function (value) {
                var number = Number(value);
                if (isNaN(number)) {
                    return 0;
                }
                if (number === 0 || !isFinite(number)) {
                    return number;
                }
                return (number > 0 ? 1 : -1) * Math.floor(Math.abs(number));
            };
            var maxSafeInteger = Math.pow(2, 53) - 1;
            var toLength = function (value) {
                var len = toInteger(value);
                return Math.min(Math.max(len, 0), maxSafeInteger);
            };

            // The length property of the from method is 1.
            return function from(arrayLike /*, mapFn, thisArg */ ) {
                // 1. Let C be the this value.
                var C = this;

                // 2. Let items be ToObject(arrayLike).
                var items = Object(arrayLike);

                // 3. ReturnIfAbrupt(items).
                if (arrayLike == null) {
                    throw new TypeError("Array.from requires an array-like object - not null or undefined");
                }

                // 4. If mapfn is undefined, then let mapping be false.
                var mapFn = arguments.length > 1 ? arguments[1] : void undefined;
                var T;
                if (typeof mapFn !== 'undefined') {
                    // 5. else
                    // 5. a If IsCallable(mapfn) is false, throw a TypeError exception.
                    if (!isCallable(mapFn)) {
                        throw new TypeError('Array.from: when provided, the second argument must be a function');
                    }

                    // 5. b. If thisArg was supplied, let T be thisArg; else let T be undefined.
                    if (arguments.length > 2) {
                        T = arguments[2];
                    }
                }

                // 10. Let lenValue be Get(items, "length").
                // 11. Let len be ToLength(lenValue).
                var len = toLength(items.length);

                // 13. If IsConstructor(C) is true, then
                // 13. a. Let A be the result of calling the [[Construct]] internal method of C with an argument list containing the single item len.
                // 14. a. Else, Let A be ArrayCreate(len).
                var A = isCallable(C) ? Object(new C(len)) : new Array(len);

                // 16. Let k be 0.
                var k = 0;
                // 17. Repeat, while k < len… (also steps a - h)
                var kValue;
                while (k < len) {
                    kValue = items[k];
                    if (mapFn) {
                        A[k] = typeof T === 'undefined' ? mapFn(kValue, k) : mapFn.call(T, kValue, k);
                    } else {
                        A[k] = kValue;
                    }
                    k += 1;
                }
                // 18. Let putStatus be Put(A, "length", len, true).
                A.length = len;
                // 20. Return A.
                return A;
            };
        }());
    };
};

function polyfill() {
    polyfillArrayFrom();
    polyfillObjectAssign();
};
