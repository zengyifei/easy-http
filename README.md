# requests
This Project aims at providing an easiest way to send get and post requests for Gophers.

Document
===
[中文](README.CN.md)

Install
===
``` sh
go get github.com/zengyifei/requests
```

Usage
===
## GET:

```Golang
// The type of requests.Params is map[string]interface{}
// send request to http://localhost:5000/?a=1&b=2
resp, err := requests.Get("http://localhost:5000/",requests.Params{
    "a": 1,
    "b": "2",
})
if err != nil {
    log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　 // Unmarshal data into YourStruct
```

## POST:
```Golang
// get file data
fp, _ := filepath.Abs("example/test.txt")
data, _ := ioutil.ReadFile(fp)

url := "http://localhost:5000"
params := requests.Params{
    "a": 1,
    "b": 2,
}

// add two fields and two files to the form 
formdata := requests.NewForm().AddField("field1", "value1").
                               AddField("field2", "value2").
                               AddFile("fileField1", "test.txt", data).
                               AddFile("fileField2", "test.txt", data)

// post form data to http://localhost:5000?a=1&b=2
resp, err := requests.Post(url, params, formdata)

if err != nil {
	log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　 // Unmarshal data into YourStruct
```
