package main

import (
    "net/url"
)

func normalizeURL(input string) (string, error) {
    parsedURL, err := url.Parse(input)

    if err != nil {
        return "", err
    }



    host := parsedURL.Hostname()
    path := parsedURL.Path
    //query := parsedURL.Query().Encode()

    x, err := url.JoinPath(host, path)
    
   if err != nil {
       return "", err
   }

    x = stripTrailingSlash(x)


    if query := parsedURL.RawQuery; query != "" {
        x += "?" + query
    }

    x = stripTrailingSlash(x)

    if fragment := parsedURL.Fragment; fragment != "" {
        x += "#" + fragment
    }

    x = stripTrailingSlash(x)

    return x, nil
}

func stripTrailingSlash(s string) string {
    if s == "" {
        return ""
    }
    if c := s[len(s)-1]; c == '/' {
        return s[:len(s)-1]
    }
    return s
}
