package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	neturl "net/url"
)

// Params save key-value pairs which will be attached behind the url
type Params map[string]interface{}

// Response
type response struct {
	*http.Response
	data []byte
}

// return response string data
func (resp *response) String() string { return string(resp.data) }

// return an io.Reader which save the response bytes data
func (resp *response) Reader() io.Reader { return interface{}(bytes.NewReader(resp.data)).(io.Reader) }

// return response bytes data
func (resp *response) Bytes() []byte { return resp.data }

// unmarshal response data into v
func (resp *response) Unmarshal(v interface{}) error { return json.Unmarshal(resp.data, v) }

// attach params behind url
func genURL(url string, params map[string]interface{}) string {
	var ps = neturl.Values{}
	for k, v := range params {
		ps.Add(k, fmt.Sprint(v))
	}

	if len(params) != 0 {
		url = fmt.Sprintf("%s?%s", url, ps.Encode())
	}
	return url
}

// Send a Get Request
func Get(url string, params map[string]interface{}) (*response, error) {
	var (
		r    = &response{}
		_url = genURL(url, params)
		err  error
	)

	if r.Response, err = http.Get(_url); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil
}

// Send a Post Request with form
func Post(url string, params map[string]interface{}, f *form) (resp *response, err error) {
	if f == nil {
		resp, err = postNilForm(url, params)
		return
	}

	if f.hasFile {
		resp, err = postWithFile(url, params, f)
	} else {
		resp, err = postWithoutFile(url, params, f)
	}
	return
}

// file which will be sended
type postfile struct {
	name string
	data []byte
}

// web form
type form struct {
	fields  map[string][]interface{}
	hasFile bool
}

// add one field into form
func (f *form) AddField(name string, value interface{}) *form {
	f.fields[name] = append(f.fields[name], value)
	return f
}

// add one file into form
func (f *form) AddFile(name, filename string, data []byte) *form {
	f.hasFile = true
	f.fields[name] = append(f.fields[name], postfile{
		filename,
		data,
	})
	return f
}

// generate a new form
func NewForm() *form {
	return &form{
		make(map[string][]interface{}),
		false,
	}
}

// send post requests with file(s)
func postWithFile(url string, params map[string]interface{}, f *form) (*response, error) {
	var (
		r    = &response{}
		_url = genURL(url, params)
		err  error
	)

	// create a simulation form
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for fieldname, values := range f.fields {
		for _, value := range values {
			switch v := value.(type) {
			// add file
			case postfile:
				// create a file-upload item
				fileWriter, err := bodyWriter.CreateFormFile(fieldname, v.name)

				if err != nil {
					return nil, err
				}

				// put file data into file-upload item
				if _, err = io.Copy(fileWriter, bytes.NewReader(v.data)); err != nil {
					return nil, err
				}

			// add field
			default:
				bodyWriter.WriteField(fieldname, fmt.Sprint(v))
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// send post requests with data
	if r.Response, err = http.Post(_url, contentType, bodyBuf); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil
}

// send post requests with only fields, no files
func postWithoutFile(url string, params map[string]interface{}, f *form) (*response, error) {
	var (
		r    = &response{}
		_url = genURL(url, params)
		err  error
		data = neturl.Values{}
	)

	for field, values := range f.fields {
		for _, value := range values {
			if data.Get(field) == "" {
				data.Set(field, fmt.Sprint(value))
				continue
			}
			data.Add(field, fmt.Sprint(value))
		}
	}

	// send post request with data
	if r.Response, err = http.PostForm(_url, data); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil

}

// send post request with nothing
func postNilForm(url string, params map[string]interface{}) (*response, error) {
	var (
		r    = &response{}
		_url = genURL(url, params)
		err  error
	)

	// send post request
	if r.Response, err = http.PostForm(_url, nil); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}
	return r, nil
}

// send post request with binary data
func PostBinary(url string, params map[string]interface{}, body io.Reader) (*response, error) {
	var (
		r    = &response{}
		_url = genURL(url, params)
		err  error
	)

	// send post request
	if r.Response, err = http.Post(_url, "multipart/form-data", body); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil

}
