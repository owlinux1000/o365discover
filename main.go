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

func valid_users(usernames_conn chan string, valid_user_conn chan string) {
	for username := range usernames_conn {
		if valid_user(username) != "" {
			valid_user_conn <- username
		} else {
			valid_user_conn <- ""
		}
	}
}

func valid_user(username string) string {
	request := Request{Username: username}
	jm, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	response, err := http.Post(
		ENDPOINT_URL,
		CONTENT_TYPE,
		bytes.NewBuffer(jm),
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
	if r.IfExistsResult == 0 {
		return username
	} else {
		return ""
	}
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

		var line_count int
		var usernames []string

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			usernames = append(usernames, scanner.Text())
			line_count += 1
		}

		usernames_conn := make(chan string, line_count)
		valid_user := make(chan string)
		go func() {
			for _, username := range usernames {
				usernames_conn <- username
			}
		}()
		for i := 0; i < line_count; i++ {
			go valid_users(usernames_conn, valid_user)
		}
		for i := 0; i < line_count; i++ {
			username := <-valid_user
			if username != "" {
				fmt.Println(username)
			}
		}
		close(usernames_conn)
		close(valid_user)
	} else if *username != "" {
		if valid_user(*username) != "" {
			fmt.Println(*username)
		}
	}
}
