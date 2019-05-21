package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"requests"
)

func main() {

	// http GET demo
	{
		url := "http://localhost:5000"
		resp, err := requests.Get(url, requests.Params{
			"a": 1,
			"b": 2,
		})
		if err != nil {
			//handle error
			log.Fatal(err)
		}
		fmt.Println(resp.String())
		fmt.Println(resp.Bytes())
		fmt.Println(resp.Reader())
		fmt.Println(resp.Unmarshal(nil))
	}

	// http Post demo
	{
		url := "http://localhost:5000"
		params := requests.Params{
			"a": 1,
			"b": 2,
		}

		data, _ := ioutil.ReadFile("test.txt")

		formdata := requests.NewForm().AddField("field1", "value1").
			AddField("field2", "value2").
			AddFile("file1", "test.txt", data).
			AddFile("file2", "test.txt", data)
		// post data to http://localhost:5000?a=1&b=2
		resp, err := requests.Post(url, params, formdata)

		if err != nil {
			//handle error
			log.Fatal(err)
		}
		fmt.Println(resp.String())
		fmt.Println(resp.Bytes())
		fmt.Println(resp.Reader())
		fmt.Println(resp.Unmarshal(nil))
	}

}
