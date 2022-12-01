package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {

	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	fmt.Printf("Response status : %s \n", resp.Status)
	fmt.Printf("Body : %s \n ", body)

}
