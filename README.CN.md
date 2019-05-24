# easyreq
[![GoDoc](https://godoc.org/github.com/zengyifei/easyreq?status.svg)](https://godoc.org/github.com/zengyifei/easyreq)

这个项目主要是为Gophers提供一种最简单的方式去发送get和post请求。

文档
===
[English](README.md)

项目初衷
===
1. 设想一下，现在要你发送一个最简单的Get请求，输出响应的字符串，不能忽略错误，你需要用多少行代码？那发送一个post请求呢？post请求带上上传文件呢？应该还蛮多的吧。  
2. 为了能够最简单地发送 GET 和 POST 请求，因为日常生活中这两种用得最多。可是我又不想每次都写那么冗长又丑又难记的代码，怎么办？自己写一个简单，易使用，易记的啰。  
3. 如果你想要自定义你的请求各种头部或参数的话，您请绕道而行，使用官方提供的net/http吧。  
4. 我看到github上有人和我有同样的想法，他写了[goreq](https://github.com/franela/goreq)，但遗憾的是，在我看来，它仍有些复杂了。它仍允许用户自己定义请求的一些选项，但这是我不允许的，灵活有时也意味着不好记。而且他的用法看上去也没我的简单。  
5. 该项目后续也不会提供太多东西。总之，简单，易用，好记是该项目想达到的目标。

安装
===
``` sh
go get github.com/zengyifei/easyreq
```

用法
===
## GET:

```Golang
// easyreq.Params 的类型是 map[string]interface{}
// 向 http://localhost:5000/?a=1&b=2 发送请求
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

// 向表单添加了两个字段和两个文件
formdata := easyreq.NewForm().AddField("field1", "value1").
                              AddField("field2", "value2").
                              AddFile("fileField1", "test1.txt", data).
                              AddFile("fileField2", "test2.txt", data)

// 向 http://localhost:5000?a=1&b=2 发送表单数据
resp, err := easyreq.Post(url, params, formdata)

// 传递一个 io.Reader 参数发送二进制数据
resp, err = easyreq.PostBinary(url, params, bytes.NewReader(data))

if err != nil {
	log.Fatal(err)
}

log.Println(resp.String())      　　　　　　　 // get response string
log.Println(resp.Bytes())       　　　　　　　 // get response bytes
log.Println(resp.Reader())      　　　　　　　 // get response reader
log.Println(resp.Unmarshal(&YourStruct))　   // Unmarshal data into YourStruct
```
