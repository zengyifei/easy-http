# requests
This Project aims at providing an easiest way to send get and post requests for Gophers.

Usage
===
## GET
The type of `requests.Params` is `map[string]interface{}`
```Golang
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
log.Println(resp.Unmarshal(&YourStruct{}))　 // Unmarshal data into YourStruct
```

## POST
```Golang
fp, _ := filepath.Abs("example/test.txt")
// get file data
data, _ := ioutil.ReadFile(fp)

url := "http://localhost:5000"
params := requests.Params{
    "a": 1,
    "b": 2,
}

formdata := requests.NewForm().AddField("field1", "value1").
                               AddField("field2", "value2").
                               AddFile("fileField1", "test.txt", data).
                               AddFile("fileField2", "test.txt", data)

// post data to http://localhost:5000?a=1&b=2
resp, err := requests.Post(url, params, formdata)

if err != nil {
	log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct{}))　 // Unmarshal data into YourStruct
```
