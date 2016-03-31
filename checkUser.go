package main

import "net/http"

func checkIfUserExists(username string) (bool, error) {
	resp, err := http.Get("https://keybase.io/" + username)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode <= 300, nil
}
