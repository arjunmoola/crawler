package main

import (
    "net/http"
    "io"
    "fmt"
)

func getHTML(rawURL string) (string, error) {
    res, err := http.Get(rawURL)

    if err != nil {
        return "", err
    }

    defer res.Body.Close()

    if res.StatusCode >= 400 {
        return "", fmt.Errorf("error: received http status code error: %v", res.StatusCode)
    }

    contentType := res.Header.Get("Content-Type")

    if !hasPrefix(contentType, "text/html"){
        return "", fmt.Errorf("error: Content-Type is not text/html. got=%s", contentType)
    }

    data, err := io.ReadAll(res.Body)

    if err != nil {
        return "", err
    }

    return string(data), err

}
