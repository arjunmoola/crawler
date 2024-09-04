package main

import (
    "testing"
)

func TestCommonDomain(t *testing.T) {
    testCases := []struct{
        url1 string
        url2 string
        expectedRes bool
    }{
        { "https://boot.dev", "https://twitter.com", false },
    }

    for _, tt := range testCases {
        t.Run("test1", func(t *testing.T) {
            res, err := compareHost(tt.url1, tt.url2)
            if err != nil {
                t.Fatal(err)
            }

            t.Log(getHostname(tt.url1))
            t.Log(getHostname(tt.url2))

            if res != tt.expectedRes {
                t.Fatalf("incorrect result. got=%v, want=%v", res, tt.expectedRes)
            }
        })
    }
}
