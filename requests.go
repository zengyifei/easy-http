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

type response struct {
	*http.Response
	data []byte
}

type Params map[string]interface{}

func (resp *response) String() string { return string(resp.data) }

func (resp *response) Reader() *bytes.Reader { return bytes.NewReader(resp.data) }

func (resp *response) Bytes() []byte { return resp.data }

func (resp *response) Unmarshal(v interface{}) error { return json.Unmarshal(resp.data, v) }

func Get(url string, params map[string]interface{}) (*response, error) {
	var (
		ps  = neturl.Values{}
		r   = &response{}
		err error
	)

	for k, v := range params {
		ps.Add(k, fmt.Sprint(v))
	}

	if len(params) != 0 {
		url = fmt.Sprintf("%s?%s", url, ps.Encode())
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

func Post(url string, params map[string]interface{}, f *form) (*response, error) {
	var (
		ps  = neturl.Values{}
		r   = &response{}
		err error
	)

	for k, v := range params {
		ps.Add(k, fmt.Sprint(v))
	}

	if len(params) != 0 {
		url = fmt.Sprintf("%s?%s", url, ps.Encode())
	}

	//创建一个模拟的form中的一个选项,这个form项现在是空的
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for fieldname, values := range f.fields {
		for _, value := range values {
			switch v := value.(type) {
			case postfile:
				//关键的一步操作, 设置文件的上传参数叫image, 文件名是tmp.png,
				//相当于现在还没选择文件, form项里选择文件的选项
				fileWriter, err := bodyWriter.CreateFormFile(fieldname, v.name)

				//bodyWriter.WriteField()
				if err != nil {
					return nil, err
				}
				//iocopy 这里相当于选择了文件,将文件放到form中
				_, err = io.Copy(fileWriter, bytes.NewReader(v.data))
				if err != nil {
					return nil, err
				}
			default:
				bodyWriter.WriteField(fieldname, fmt.Sprint(v))
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//发送post请求到服务端
	if r.Response, err = http.Post(url, contentType, bodyBuf); err != nil {
		return nil, err
	}
	defer r.Response.Body.Close()

	if r.data, err = ioutil.ReadAll(r.Response.Body); err != nil {
		return nil, err
	}

	return r, nil
}

type form struct {
	fields map[string][]interface{}
}

type postfile struct {
	name string
	data []byte
}

func (f *form) AddField(name string, value interface{}) *form {
	f.fields[name] = append(f.fields[name], value)
	return f
}

func (f *form) AddFile(name, filename string, data []byte) *form {
	f.fields[name] = append(f.fields[name], postfile{
		filename,
		data,
	})
	return f
}

func NewForm() *form {
	return &form{map[string][]interface{}{}}
}
