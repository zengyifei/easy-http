# requests
This Project aims at providing an easiest way to send http requests for Gophers.

Comparison
===
At past,if I want to send a request and get the response text,the code will look like below:
```Golang
resp, err := http.Get("http://localhost:5000/?a=1&b=2)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
data, err := ioutil.ReadAll()
if err != nil {
    log.Fatal(err)
}
log.Println(string(data))
```
How ugly and long it looks.Now it can become like this:
```Golang
resp, err := requests.Get("http://localhost:5000/?a=1&b=2")
if err != nil {
    log.Fatal(err)
}
log.Println(resp.String())
```
It looks much better, isn't it?

Usage
===
## Get
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
log.Println(resp.String())      　　　　　　　// get response string
log.Println(resp.Bytes())       　　　　　　　// get response bytes
log.Println(resp.Reader())      　　　　　　　// get response reader
log.Println(resp.Unmarshal(&YourStruct{}))　 // Unmarshal data into YourStruct
```
