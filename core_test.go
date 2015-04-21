package sloth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"testing"
)

type Item struct{}

func (item Item) Get(request *http.Request) (int, interface{}, http.Header) {
	items := []string{"item1", "item2"}
	data := map[string][]string{"items": items}
	return http.StatusOK, data, nil
}

func (item Item) Post(request *http.Request) (int, interface{}, http.Header) {
	data := fmt.Sprintf("You sent: %s", request.Form.Get("hello"))
	return http.StatusOK, data, http.Header{"Content-Type": {"text/plain"}}
}

type Upload struct{}

func (upload Upload) Post(request *http.Request) (int, interface{}, http.Header) {
	if err := request.ParseMultipartForm(1024); err != nil {
		return http.StatusBadRequest, err.Error(), http.Header{"Content-Type": {"text/plain"}}
	}
	data := request.MultipartForm.Value["title"]
	return http.StatusOK, data, nil
}

func TestMain(m *testing.M) {

	item := new(Item)
	upload := new(Upload)

	var api = NewAPI()

	api.AddResource(item, "/items", "/bar", "/baz")
	api.AddResource(upload, "/upload")

	go api.Start(3000)

	os.Exit(m.Run())
}

func TestBasicGet(t *testing.T) {

	resp, err := http.Get("http://localhost:3000/items")
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "{\n  \"items\": [\n    \"item1\",\n    \"item2\"\n  ]\n}" {
		t.Error("Not equal.")
	}

}

func TestBasicPostWithTextPlainResponse(t *testing.T) {

	resp, err := http.PostForm("http://localhost:3000/items", url.Values{"hello": {"sloth"}})
	if err != nil {
		t.Error(err)
	}
	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Error("Content-Type wrong.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "You sent: sloth" {
		t.Error("Not equal.")
	}

}

func TestMultipartFormDataPost(t *testing.T) {

	params := map[string]string{
		"title":       "My Document",
		"author":      "Simon Eisenmann",
		"description": "A document with all the my secrets",
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	request, _ := http.NewRequest("POST", "http://localhost:3000/upload", data)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type wrong.")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) != "[\n  \"My Document\"\n]" {
		t.Error("Not equal.")
	}

}
