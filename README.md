# easyreq
[![GoDoc](https://godoc.org/github.com/zengyifei/easyreq?status.svg)](https://godoc.org/github.com/zengyifei/easyreq)

This Project aims at providing an easiest way to send get and post requests for Gophers.

Document
===
[中文](README.CN.md)

Why easyreq?
===
1. Imagine now that you want to send a simple Get request, output a response string, can't ignore the error, how many lines of code do you need? Then send a post request? Post request with upload files?It should be quite a lot.
2. I don't want to write such a long, ugly and hard-to-remember code every time. I believe other people are like this. So holding a simple, easy to use, easy to remember idea, I wrote this project.
3. If you want to customize the various headers or parameters for your request, please take a detour and use official built-in library net/http.  
4. I saw someone on github who had the same idea as me. He wrote [goreq](https://github.com/franela/goreq), but unfortunately, in my opinion, it is still a bit complicated. It still allows the user to define some of the options for the request, sometimes it means not easy to remember. So this is not allowed in my mind.And the usage of his lib doesn't seem to be as simple as mine.

Install
===
``` sh
go get github.com/zengyifei/easyreq
```

Usage
===
## GET:

```Golang
// The type of easyreq.Params is map[string]interface{}
// send request to http://localhost:5000/?a=1&b=2
resp, err := easyreq.Get("http://localhost:5000/",easyreq.Params{
    "a": 1,
    "b": "2",
})
if err != nil {
    log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　   // Unmarshal data into YourStruct
```

## POST:
```Golang
// get file data
data, _ := ioutil.ReadFile("filepath")

url := "http://localhost:5000"
params := easyreq.Params{
    "a": 1,
    "b": 2,
}

// add two fields and two files to the form 
formdata := easyreq.NewForm().AddField("field1", "value1").
                              AddField("field2", "value2").
                              AddFile("fileField1", "test1.txt", data).
                              AddFile("fileField2", "test2.txt", data)

// post form data to http://localhost:5000?a=1&b=2
resp, err := easyreq.Post(url, params, formdata)

// pass an io.Reader param to post binary data
resp, err = easyreq.PostBinary(url, params, bytes.NewReader(data))

if err != nil {
	log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　   // Unmarshal data into YourStruct
```
