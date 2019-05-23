package requests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"net/url"
	"reflect"
	"io/ioutil"
	"bytes"
)

type testFile struct {
	name string
	data []byte
}

var (
	params = Params{
		"name": "John",
		"age": 18,
		"height": float64(55.6),
	}

	postParams = map[string][]interface{}{
		"width": {"20",40},
		"height": {30,float64(43),float32(12)},
	}

	filecontent = "this is a test file."
	filedata = []byte(filecontent)

	testFiles = map[string][]testFile {
		"firstFile": { {name: "firstFile.txt", data: filedata}, },
		"secondFile": { {name: "secondFile.txt", data: filedata}, },
		"thirdFile": { {name: "thirdFile.txt", data: filedata}, },
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
	_, err := Get(ts.URL, params)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Post(ts.URL, params, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostNilForm(t *testing.T) {
	PostHandler := func(w http.ResponseWriter, r *http.Request) {}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

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
		if !reflect.DeepEqual(want, got){
			t.Errorf("want: %v, got: %v", want, got)
		}

	}
	ts := httptest.NewServer(http.HandlerFunc(PostHandler))

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
					want[name] = append(want[name],fmt.Sprint(v))
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

					got[fieldname] = append(got[fieldname],testFile{
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

	_, err := PostBinary(ts.URL, params, bytes.NewReader(filedata))
	if err != nil {
		t.Fatal(err)
	}
}