"use strict"

requirejs.config({
    "baseUrl": 'scripts',
    "paths": {
    	"jquery": "//code.jquery.com/jquery-2.1.4.min"
    }
});

requirejs(['main']);
