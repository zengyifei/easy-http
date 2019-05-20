package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
)

type response struct {
	*http.Response
	data []byte
}

type Params map[string]interface{}

func (resp *response) String() string { return string(resp.data) }

func (resp *response) Reader() *bytes.Reader { return bytes.NewReader(resp.data) }

func (resp *response) Bytes() []byte { return resp.data }

func (resp *response) Unmarshal(v interface{}) error { return json.Unmarshal(resp.data, v) }

func Get(url string, ps map[string]interface{}) (*response, error) {
	var (
		params = neturl.Values{}
		r      = &response{}
		err    error
	)

	for k, v := range ps {
		params.Add(k, fmt.Sprint(v))
	}

	if len(params) != 0 {
		url = fmt.Sprintf("%s?%s", url, params.Encode())
	}

	if r.Response, err = http.Get(url); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil
}
