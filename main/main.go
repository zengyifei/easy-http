package main

import (
	"fmt"
	"log"
	"requests"
)

func main() {
	resp, err := requests.Get("http://localhost:5000", requests.Params{
		"a": 1,
		"b": 2,
		"c": 3,
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
