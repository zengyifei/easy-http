# requests
这个项目主要是为Gophers提供一种最简单的方式去发送get和post请求。

文档
===
[English](README.md)

安装
===
``` sh
go get github.com/zengyifei/requests
```

用法
===
## GET:

```Golang
// requests.Params 的类型是 map[string]interface{}
// 向 http://localhost:5000/?a=1&b=2 发送请求
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
log.Println(resp.Unmarshal(&YourStruct))　   // Unmarshal data into YourStruct
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

// 向表单添加了两个字段和两个文件
formdata := requests.NewForm().AddField("field1", "value1").
                               AddField("field2", "value2").
                               AddFile("fileField1", "test.txt", data).
                               AddFile("fileField2", "test.txt", data)

// 向 http://localhost:5000?a=1&b=2 发送表单数据
resp, err := requests.Post(url, params, formdata)

if err != nil {
	log.Fatal(err)
}
log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　   // Unmarshal data into YourStruct
```
