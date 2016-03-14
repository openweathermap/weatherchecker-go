package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

// MakeURL creates a URL out of a template.
func MakeURL(urlTemplate string, data interface{}) (urlString string) {
	var urlBuf = new(bytes.Buffer)

	var t = template.New("URL template")
	t, _ = t.Parse(urlTemplate)
	t.Execute(urlBuf, data)

	urlString = urlBuf.String()
	return urlString
}

// Download retrieves data from the specified HTTP address.
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

// MakeSlug creates a slug out of specified title.
func MakeSlug(title string) string {
	return strings.ToLower(strings.Replace(title, " ", "_", -1))
}
