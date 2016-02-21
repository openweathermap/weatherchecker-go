package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"
)

func MakeUrl(url_template string, data interface{}) (urlString string) {
	var urlBuf = new(bytes.Buffer)

	var t = template.New("URL template")
	t, _ = t.Parse(url_template)
	t.Execute(urlBuf, data)

	urlString = urlBuf.String()
	return urlString
}

func Download(url string) (data string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		readallContents, _ := ioutil.ReadAll(resp.Body)
		data = string(readallContents)
	}
	return data, err
}
