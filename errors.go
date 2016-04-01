package main

import "errors"

var AdapterPanicErr = errors.New("Adapter panic.")
var NoAdaptFunc = errors.New("No adapt function.")
var NodeErr = errors.New("Node error.")
var MalformedEntry = errors.New("Malformed entry.")
var MalformedResponse = errors.New("Malformed response.")
var InvalidTimeOffsetString = errors.New("Invalid time offset string.")
var InvalidTimeString = errors.New("Invalid time string.")

var CompErr = errors.New("Mismatch.")
