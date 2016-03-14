package main

import "time"

func CreateHistoryDataEntry(location LocationEntry, source SourceEntry, measurements MeasurementArray, wtype string, url string, err error) (entry HistoryDataEntry) {
	var status int64
	var message string
	if err != nil {
		status = 500
		message = err.Error()
	} else {
		status = 200
		message = "OK"
	}
	entry = MakeHistoryDataEntry()
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
func StatSource(location LocationEntry, source SourceEntry, wtype string) (entry HistoryDataEntry) {
	var err error
	var url string
	var raw string
	measurements := make(MeasurementArray, 0)

	adaptFunc, adaptFuncLookupErr := GetAdaptFunc(source.Name, wtype)

	if adaptFuncLookupErr == nil {
		url = MakeURL(source.Urls[wtype], UrlData{Source: source, Location: location})

		var downloadErr error
		raw, downloadErr = Download(url)

		if downloadErr != nil {
			measurements = AdaptStub(raw)
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

func PollAll(h *WeatherHistory, locations []LocationEntry, sources []SourceEntry, wtypes []string) (dataset []HistoryDataEntry) {
	dt := time.Now().Unix()

	dataChan := make(chan HistoryDataEntry, 9999)
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
