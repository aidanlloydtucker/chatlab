package main

import (
	"fmt"
	"net/http"
)

func main() {
    // Will work
    firstExists, err := checkIfUserExists("leijurv")
    if err != nil {
		panic(err)
	}
    fmt.Println("leijurv exists:", firstExists)

    // Will not work
    secondExists, err := checkIfUserExists("eververververvcewdx")
    if err != nil {
		panic(err)
	}
    fmt.Println("eververververvcewdx exists:", secondExists)
}

func checkIfUserExists(username string) (bool, error){
    resp, err := http.Get("https://keybase.io/" + username)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

    return resp.StatusCode <= 300, nil
}
