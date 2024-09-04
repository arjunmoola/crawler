package main

import (
    "os"
    "io"
    "fmt"
    "net/url"
    "sync"
)

type config struct {
    pages map[string]int
    maxPages int
    baseURL *url.URL
    mu *sync.Mutex
    concurrencyControl chan struct{}
    wg *sync.WaitGroup
    logger io.Writer
    totalVisited int
}

func NewConfig(baseURL string, maxConcurrency int, maxPages int, logger io.Writer) (*config, error) {
    u, err := url.Parse(baseURL)
    if err != nil {
        return nil, err
    }

    c := &config{}
    c.baseURL = u
    c.wg = &sync.WaitGroup{}
    c.concurrencyControl = make(chan struct{}, maxPages)
    c.pages = make(map[string]int)
    c.mu = &sync.Mutex{}
    c.maxPages = maxPages
    c.logger = logger
    return c, nil
}

func (c *config) getPageLength() (int, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    n := len(c.pages)

    return n, n < c.maxPages
}

func (c *config) compareHost(url1 string) (bool, error) {
    parsedUrl, err := url.Parse(url1)
    if err != nil {
        return false, err
    }

    return parsedUrl.Hostname() == c.baseURL.Hostname(), nil
}

func (c *config) addPageVisit(normalizedURL string) (isFirst bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if _, ok := c.pages[normalizedURL]; !ok {
        c.pages[normalizedURL]++
        return true
    }
    c.pages[normalizedURL]++
    return
}

func (c *config) incVisit() {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.totalVisited++
}

func (c *config) crawlPage(rawCurrentURL string) {
    normalizedURL, err := normalizeURL(rawCurrentURL)
    if err != nil {
        return
    }

    c.incVisit()

    if c.logger != nil {
        fmt.Fprintf(c.logger, "\x1B[38;5;47mcrawling page:\x1B[0m \x1B[38;5;166m\x1B[4m%s\x1B[0m\n", rawCurrentURL)
    }

    _, more := c.getPageLength()

    if !more {
        return
    }

    sameHost, err := c.compareHost(rawCurrentURL)

    if err != nil || !sameHost {
        return
    }

    hasNotVisited := c.addPageVisit(normalizedURL)

    if !hasNotVisited {
        return
    }

    //if c.logger != nil {
    //    fmt.Fprintf(c.logger, "\x1B[38;5;47mcrawling page: \x1B[4m%s\x1B[0m\n", rawCurrentURL)
    //}

    htmlBody, err := getHTML(rawCurrentURL)
    
    if err != nil {
        return
    }

    foundURLS, err := getURLSFromHTML(htmlBody, c.baseURL.String())

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    for _, url := range foundURLS {
        c.wg.Add(1)
        c.concurrencyControl <- struct{}{}
        go func() {
            defer c.wg.Done()
            <-c.concurrencyControl
            c.crawlPage(url)
        }()
    }

}

func getHostname(u string) (string, error) {
    parsedUrl, err := url.Parse(u)
    if err != nil {
        return "", err
    }
    return parsedUrl.Hostname(), nil
}

