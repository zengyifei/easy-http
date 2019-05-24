# easyreq
[![GoDoc](https://godoc.org/github.com/zengyifei/easyreq?status.svg)](https://godoc.org/github.com/zengyifei/easyreq)

This Project aims at providing an easiest way to send get and post requests for Gophers.

Document
===
[中文](README.CN.md)

Why easyreq?
===
Golang net/http lib gives us the flexibility to customize requests to handle the various situations that occur.But it takes me several lines when I just want to send an easiest GET request. Like this:
``` Golang
resp, err := http.Get("http://localhost:5000")
if err != nil {
    fmt.Println(err)
}
defer resp.Body.Close()
data, err := ioutil.ReadAll(resp.Body)
if err != nil {
    fmt.Println(err)
}
fmt.Println("Response string：", string(data))
```
Then what about a POST request with multiple fields and files? More lines. And maybe you also need to baidu or Google to get the usage about how to post files.  
That's the reason why I write easyreq which aims at providing an easiest way to help gophers to send get and post requests.  
It should be easy to use and easy to remember and maybe will not add too much funcitons in the future in order not to complicate it. 


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
