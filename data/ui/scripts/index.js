"use strict"

requirejs.config({
    "baseUrl": 'scripts',
    "paths": {
        "jquery": "//code.jquery.com/jquery-2.1.4.min",
        "datatables": "//cdn.datatables.net/1.10.10/js/jquery.dataTables.min"
    }
});

requirejs(['main']);
