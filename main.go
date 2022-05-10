package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const ENDPOINT_URL = "https://login.microsoftonline.com/common/GetCredentialType"
const CONTENT_TYPE = "application/json"

type Request struct {
	Username string
}

type Result struct {
	Username       string
	Display        string
	IfExistsResult int
}

func valid_user(username string) int {
	request := Request{Username: username}
	json_request, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	response, err := http.Post(
		ENDPOINT_URL,
		CONTENT_TYPE,
		bytes.NewBuffer(json_request),
	)
	defer response.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var r Result
	if json.Unmarshal(body, &r); err != nil {
		panic(err)
	}
	return r.IfExistsResult
}

func main() {
	file := flag.String("f", "", "Specify a file")
    username := flag.String("u", "", "Specify one username")    
	flag.Parse()
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			username := scanner.Text()
			if valid_user(username) == 0 {
				fmt.Println(username)
			}
		}
	} else if *username != "" {
        if valid_user(*username) == 0 {
            fmt.Println(*username)
        }
    }
}
