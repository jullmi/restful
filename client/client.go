package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

var u1 = User{"admin", "admin"}
var u2 = User{"user", "password"}
var u3 = User{"", "bbb"}

const addEndPoint = "/add"
const getEndPoint = "/get"
const deleteEndPoint = "/delete"
const timeEndPoint = "/time"

func deleteEndpoint(server string, user User) int {
	userMarshal, _ := json.Marshal(user)
	u := bytes.NewReader(userMarshal)

	req, err := http.NewRequest("DELETE", server+deleteEndPoint, u)
	if err != nil {
		fmt.Println("Error in req:", err)
		return http.StatusInternalServerError
	}

	req.Header.Set("Content-type", "application/json")

	c := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if resp == nil {
		return http.StatusNotFound
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	fmt.Print("/delete returned: ", string(data))

	if err != nil {
		fmt.Println(err)
	}

	return resp.StatusCode

}

func getEndpoint(server string, user User) int {
	userMarshal, _ := json.Marshal(user)

	u := bytes.NewReader(userMarshal)

	req, err := http.NewRequest("GET", server+getEndPoint, u)
	if err != nil {
		fmt.Println("Error:", err)
		return http.StatusInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if resp == nil {
		return http.StatusNotFound
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	fmt.Print("/get returned: ", string(data))
	if err != nil {
		fmt.Println("Error:", err)
	}

	return resp.StatusCode
}

func addEndpoint(server string, user User) int {
	userMarshal, _ := json.Marshal(user)
	u := bytes.NewReader(userMarshal)

	req, err := http.NewRequest("POST", server+addEndPoint, u)
	if err != nil {
		fmt.Println("Error in req", err)
		return http.StatusInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)

	}

	if resp == nil || (resp.StatusCode == http.StatusNotFound) {
		return resp.StatusCode
	}

	defer resp.Body.Close()
	return resp.StatusCode

}

func timeEndpoint(server string) (int, string) {
	req, err := http.NewRequest("POST", server+timeEndPoint, nil)

	if err != nil {
		fmt.Println("Error in req", err)
		return http.StatusInternalServerError, ""
	}

	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if resp == nil || (resp.StatusCode == http.StatusNotFound) {
		return resp.StatusCode, ""
	}

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(data)

}

func slashEndpoint(server string, URL string) (int, string) {
	req, err := http.NewRequest("POST", server+URL, nil)

	if err != nil {
		fmt.Println("Error in req:", err)
	}

	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if resp == nil {
		return resp.StatusCode, ""
	}

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(data)

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments")
		fmt.Println("Need: server")
		return
	}

	server := os.Args[1]

	fmt.Println("/add")
	HTTPCode := addEndpoint(server, u1)
	if HTTPCode != http.StatusOK {
		fmt.Println("u1 return code", HTTPCode)
	} else {
		fmt.Println("u1 Data added:", u1, HTTPCode)
	}

	HTTPCode = addEndpoint(server, u2)
	if HTTPCode != http.StatusOK {
		fmt.Println("u2 return code:", HTTPCode)
	} else {
		fmt.Println("u2 Data added", u2, HTTPCode)
	}

	HTTPCode = addEndpoint(server, u3)
	if HTTPCode != http.StatusOK {
		fmt.Println("u3 return code:", HTTPCode)
	} else {
		fmt.Println("u3 Data added", u3, HTTPCode)
	}

	fmt.Println("/get")
	HTTPCode = getEndpoint(server, u1)
	fmt.Println("/get u1 return code:", HTTPCode)
	HTTPCode = getEndpoint(server, u2)
	fmt.Println("/get u2 return code:", HTTPCode)
	HTTPCode = getEndpoint(server, u3)
	fmt.Println("/get u3 return code:", HTTPCode)

	fmt.Println("/delete")
	HTTPCode = deleteEndpoint(server, u1)
	fmt.Println("/delete u1 return code:", HTTPCode)
	HTTPCode = deleteEndpoint(server, u1)
	fmt.Println("/delete u1 return code:", HTTPCode)
	HTTPCode = deleteEndpoint(server, u2)
	fmt.Println("/delete u2 return code:", HTTPCode)
	HTTPCode = deleteEndpoint(server, u3)
	fmt.Println("/delete u3 return code:", HTTPCode)


	fmt.Println("/time")
	HTTPCode, myTime := timeEndpoint(server)
	fmt.Println("/time return:", HTTPCode, "time:", myTime)
	time.Sleep(time.Second)
	HTTPCode, myTime = timeEndpoint(server)
	fmt.Println("/time return:", HTTPCode, "time:", myTime)

	fmt.Println("/")
	URL := "/"
	HTTPCode, response := slashEndpoint(server, URL)
	fmt.Println("/ return: ", HTTPCode, "with response: ", response)
}
