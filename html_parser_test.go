package main

import (
    "testing"
    //"strings"
    "slices"
)

func TestGetURLSFromHTML(t *testing.T) {
    testCases := []struct{
        name string
        inputURL string
        inputBody string
        expected []string
    }{
        {
        name:     "absolute and relative URLs",
        inputURL: "https://blog.boot.dev",
        inputBody: `
    <html>
        <body>
            <a href="/path/one">
                <span>Boot.dev</span>
            </a>
            <a href="https://other.com/path/one">
                <span>Boot.dev</span>
            </a>
            <a href="/path/one">
                <span>Boot2.dev</span>
            </a>
            <a href="/path/one/a">
                <span>Boot2.dev</span>
            </a>
            <a href="https://google.com/search"></a>
            <p> Hello World
            <a href="https://google.com/search"></a>
            </p>
            <a href="https://www.boot.dev">Learn Backend Development</a>
        </body>
    </html>
    `,
        expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/one/a", "https://other.com/path/one", "https://google.com/search", "https://www.boot.dev"},
        },
        {
            name: "same links",
            inputURL: "https://boot.dev",
            inputBody:`
<html>
    <body>
        <a href="/path/one"></a>
        <a href="/path/two"></a>
        <a href="/path/two"></a>
    </body>
</html>`,
            expected: []string{ "https://boot.dev/path/one", "https://boot.dev/path/two" },
        },
        {
            name: "deep nested",
            inputURL: "https://boot.dev",
            inputBody: `
<html>
    <body>
        <p>
            <div>
                <a href="/path/one"></a>
            </div>

            <div>
                <div>
                    <p><a href="/path/one"></a></p>
                </div>
            </div>
        </p>
    </body>
</html>`,
            expected: []string{ "https://boot.dev/path/one" },

        },
        {
            name: "react andy",
            inputURL: "https://blog.boot.dev",
            inputBody: `
<html>
    <body>
        <a href="https://blog.boot.dev"><span>Go to Boot.dev, you React Andy</span></a>
    </body>
</html>`,
    expected: []string{ "https://blog.boot.dev" },
        },
    }

    for i, tt := range testCases {
        t.Run(tt.name, func(t *testing.T) {
            found, err := getURLSFromHTML(tt.inputBody, tt.inputURL)
            if err != nil {
                t.Fatal(err)
            }
            slices.Sort(tt.expected)
            if len(found) != len(tt.expected) {
                t.Fatalf("Test %d %s, incorrect number of urls found. got=%d, want=%d", i, tt.name,  len(found), len(tt.expected))
            }

            for j := range len(found) {
                if found[j] != tt.expected[j] {
                    t.Errorf("Test %d %s: incorrect val at index %d. got=%s, want=%s", i, tt.name, j, found[j], tt.expected[j])
                }
            }
        })
    }
}
