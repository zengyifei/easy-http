package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type testFile struct {
	name string
	data []byte
}

type Response struct {
	Result string `json:"result"`
}

var (
	// get method params
	params = Params{
		"name":   "John",
		"age":    18,
		"height": float64(55.6),
	}

	// post method params
	postParams = map[string][]interface{}{
		"width":  {"20", 40},
		"height": {30, float64(43), float32(12)},
	}

	// simulate real file content
	filecontent = "this is a test file."
	filedata    = []byte(filecontent)

	// simulate multi files
	testFiles = map[string][]testFile{
		"firstFile":  {{name: "firstFile.txt", data: filedata}},
		"secondFile": {{name: "secondFile.txt", data: filedata}},
		"thirdFile":  {{name: "thirdFile.txt", data: filedata}},
	}

	testResponse = Response{
		Result: "success",
	}
)

func TestURLParams(t *testing.T) {

	GETHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		for key, value := range params {
			if v := query.Get(key); fmt.Sprint(value) != v {
				t.Errorf("query param %s = %s; want = %s", key, v, value)
			}
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(GETHandler))

	// test get url params
	_, err := Get(ts.URL, params)
	if err != nil {
		t.Fatal(err)
	}

	// test post url params
	_, err = Post(ts.URL, params, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostNilForm(t *testing.T) {
	PostHandler := func(w http.ResponseWriter, r *http.Request) {}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

	// post nil form
	_, err := Post(ts.URL, params, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostWithoutFile(t *testing.T) {
	form := NewForm()
	for fieldname, values := range postParams {
		for _, v := range values {
			form.AddField(fieldname, v)
		}
	}

	want := url.Values{}

	for name, values := range postParams {
		for _, v := range values {
			want.Add(name, fmt.Sprint(v))
		}
	}

	PostHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		got := r.PostForm
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v, got: %v", want, got)
		}

	}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

	// post only fields, no file
	_, err := Post(ts.URL, params, form)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostWithMultiFiles(t *testing.T) {
	form := NewForm()
	for fieldname, values := range postParams {
		for _, v := range values {
			form.AddField(fieldname, v)
		}
	}

	for fieldname, files := range testFiles {
		for _, file := range files {
			form.AddFile(fieldname, file.name, file.data)
		}
	}

	PostHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 5); err != nil {
			t.Fatal(err)
		}

		// check post fields
		{
			want := map[string][]string{}
			got := r.MultipartForm.Value
			for name, values := range postParams {
				for _, v := range values {
					want[name] = append(want[name], fmt.Sprint(v))
				}
			}

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want: %v[%T], got: %v[%T]", want, want, got, got)
			}
		}
		// check post files
		{
			want := testFiles
			got := map[string][]testFile{}
			for fieldname, files := range r.MultipartForm.File {
				for _, file := range files {
					f, err := file.Open()
					if err != nil {
						t.Fatal(err)
					}
					defer f.Close()

					data, err := ioutil.ReadAll(f)
					if err != nil {
						t.Fatal(err)
					}

					got[fieldname] = append(got[fieldname], testFile{
						name: file.Filename,
						data: data,
					})
				}
			}

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want: %v, got: %v", want, got)
			}
		}

	}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

	// post fields and multi files
	_, err := Post(ts.URL, params, form)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostBinary(t *testing.T) {
	PostHandler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		want := filedata
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

	// post binary data
	_, err := PostBinary(ts.URL, params, bytes.NewReader(filedata))
	if err != nil {
		t.Fatal(err)
	}
}

func TestResponse(t *testing.T) {
	GETHandler := func(w http.ResponseWriter, r *http.Request) {
		respdata, err := json.Marshal(Response{
			Result: "success",
		})
		if err != nil {
			t.Fatal(err)
		}
		w.Write(respdata)
	}
	ts := httptest.NewServer(http.HandlerFunc(GETHandler))

	resp, err := Get(ts.URL, params)
	if err != nil {
		t.Fatal(err)
	}

	// test response Bytes() method
	t.Run("TestResponse_Bytes", func(t *testing.T) {
		got := resp.Bytes()
		want, err := json.Marshal(testResponse)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v[%T], got: %v[%T]", want, want, got, got)
		}
	})

	// test response String() method
	t.Run("TestResponse_String", func(t *testing.T) {
		got := resp.String()
		testdata, err := json.Marshal(testResponse)
		if err != nil {
			t.Fatal(err)
		}

		want := string(testdata)

		if want != got {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	// test response Reader() method
	t.Run("TestResponse_Reader", func(t *testing.T) {
		respRd := resp.Reader()

		got, err := ioutil.ReadAll(respRd)
		if err != nil {
			t.Fatal(err)
		}

		want, err := json.Marshal(testResponse)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	// test response Unmarshal() method
	t.Run("TestResponse_Unmarshal", func(t *testing.T) {
		want := testResponse
		got := Response{}
		if err := resp.Unmarshal(&got); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

}
