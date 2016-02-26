package core

import (
	"time"

	"github.com/owm-inc/weatherchecker-go/adapters"
	"github.com/owm-inc/weatherchecker-go/models"
	"github.com/owm-inc/weatherchecker-go/util"
)

func CreateHistoryDataEntry(location models.LocationEntry, source models.SourceEntry, measurements adapters.MeasurementArray, wtype string, url string, err error) (entry models.HistoryDataEntry) {
	var status int64
	var message string
	if err != nil {
		status = 500
		message = err.Error()
	} else {
		status = 200
		message = "OK"
	}
	entry = models.MakeHistoryDataEntry()
	entry.Status = status
	entry.Message = message
	entry.Location = location
	entry.Source = source
	entry.Measurements = measurements
	entry.WType = wtype
	entry.Url = url

	return entry
}

// StatSource polls single provider for data on specified location and wtype.
func StatSource(location models.LocationEntry, source models.SourceEntry, wtype string) (entry models.HistoryDataEntry) {
	var err error
	var url string
	var raw string
	measurements := make(adapters.MeasurementArray, 0)

	adaptFunc, adaptFuncLookupErr := adapters.GetAdaptFunc(source.Name, wtype)

	if adaptFuncLookupErr == nil {
		url = util.MakeURL(source.Urls[wtype], models.UrlData{Source: source, Location: location})

		var downloadErr error
		raw, downloadErr = util.Download(url)

		if downloadErr != nil {
			measurements = adapters.AdaptStub(raw)
			err = downloadErr
		} else {
			var adaptErr error
			measurements, adaptErr = adaptFunc(raw)

			err = adaptErr
		}

	} else {
		err = adaptFuncLookupErr
	}

	entry = CreateHistoryDataEntry(location, source, measurements, wtype, url, err)

	return entry
}

func PollAll(h *models.WeatherHistory, locations []models.LocationEntry, sources []models.SourceEntry, wtypes []string) (dataset []models.HistoryDataEntry) {
	dt := time.Now().Unix()

	dataChan := make(chan models.HistoryDataEntry, 9999)
	doneChan := make(chan struct{})

	go func() {
		for entry := range dataChan {
			h.Add(entry)
			dataset = append(dataset, entry)
		}
		doneChan <- struct{}{}
	}()

	for _, location := range locations {
		for _, source := range sources {
			for _, wtype := range wtypes {
				data := StatSource(location, source, wtype)
				data.RequestTime = dt

				dataChan <- data
			}
		}
		time.Sleep(6 * time.Second)
	}
	close(dataChan)
	<-doneChan

	return dataset
}
