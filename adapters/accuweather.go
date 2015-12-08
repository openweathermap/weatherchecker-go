package adapters

import (
        "errors"
        "strconv"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

func AccuweatherAdaptCurrentHumidity (htmlString string) (float64, error) {
    nodeErr := errors.New(`Node not found`)

    htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
    if htmlErr != nil {return float64(0), htmlErr}

    nodeGroup1 := htmlDoc.Find(`tr`).Nodes
    nodeIndex1 := 5
    if len(nodeGroup1) - 1 < nodeIndex1 {return float64(0), nodeErr}
    node1 := nodeGroup1[nodeIndex1]

    nodeGroup2 := goquery.NewDocumentFromNode(node1).Contents().Nodes
    nodeIndex2 := 3
    if len(nodeGroup2) - 1 < nodeIndex2 {return float64(0), nodeErr}
    node2 := nodeGroup2[nodeIndex2]

    finalData := node2.FirstChild.Data

    humidityRaw := finalData
    humidityString := strings.TrimRight(strings.TrimSpace(humidityRaw), "%")

    humidity, convErr := strconv.ParseFloat(humidityString, 64)
    if convErr != nil {return float64(0), convErr}

    return humidity, nil
}

func AccuweatherAdaptCurrentTemp (htmlString string) (float64, error) {
    nodeErr := errors.New(`Node not found`)

    htmlDoc, htmlErr := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
    if htmlErr != nil {return float64(0), htmlErr}

    nodeGroup1 := htmlDoc.Find(`tr.temp`).Find(`td.first-col`).Contents().Nodes
    nodeIndex1 := 0
    if len(nodeGroup1) - 1 < nodeIndex1 {return float64(0), nodeErr}
    node1 := nodeGroup1[nodeIndex1]

    finalData := node1.Data

    tempRaw := finalData
    tempString := strings.TrimRight(strings.TrimSpace(tempRaw), "°")
    temp, convErr := strconv.ParseFloat(tempString, 64)

    if convErr != nil {return float64(0), convErr}

    return temp, nil
}

func AccuweatherAdaptCurrentWeather (htmlString string) MeasurementArray {
    var measurements MeasurementArray

    humidity, _ := AccuweatherAdaptCurrentHumidity(htmlString)
    precipitation := float64(0)
    pressure := float64(0)
    temp, _ := AccuweatherAdaptCurrentTemp(htmlString)
    wind := float64(0)

    measurements = append(measurements, MeasurementSchema{Humidity:humidity, Precipitation:precipitation, Pressure:pressure, Temp:temp, Wind:wind})

    return measurements
}
