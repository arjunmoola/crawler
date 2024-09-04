package main

import (
    "testing"
)

func TestNormalizeURL(t *testing.T) {
    testCases := []struct{
        name string
        inputURL string
        expectedURL string
    }{
        { "remove scheme", "https://blog.boot.dev/path", "blog.boot.dev/path" },
        { "remove scheme 2", "https://blog.boot.dev/path/", "blog.boot.dev/path" },
        { "remove scheme 3", "http://blog.boot.dev/path/", "blog.boot.dev/path" },
        { "remove schem3 4", "http://blog.boot.dev/path", "blog.boot.dev/path" },
        { "query", "https://blog.boot.dev/path?a=b&b=c/", "blog.boot.dev/path?a=b&b=c" },

    }

    for i, tt := range testCases {
        t.Run(tt.name, func(t *testing.T) {
            normalized, err := normalizeURL(tt.inputURL)

            if err != nil {
                t.Fatalf("error: %v", err)
            }

            if normalized != tt.expectedURL {
                t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tt.name, tt.expectedURL, normalized)
            }
        })
    }
}
