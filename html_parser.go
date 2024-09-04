package main

import (
    //"fmt"
    "io"
    "golang.org/x/net/html"
    "maps"
    "slices"
    "iter"
    "strings"
)

func hasPrefix(text, pattern string) bool {
    return len(pattern) <= len(text) && text[:len(pattern)] == pattern
}

func getURLSFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
    root, err := html.Parse(strings.NewReader(htmlBody))
    if err != nil {
        return []string{}, nil
    }

    cache := make(map[string]struct{})
    
    var dfs func(*html.Node)

    dfs = func(n *html.Node) {
        if n == nil {
            return
        }

        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    u := attr.Val
                    if !hasPrefix(u, "http") {
                        u = rawBaseURL + u
                    }
                    cache[u] = struct{}{}
                }
            }
        }

        for c := n.FirstChild; c != nil; c = c.NextSibling {
            dfs(c)
        }
    }

    dfs(root)

    found := make([]string, 0, len(cache))

    for key := range cache {
        found = append(found, key)
    }

    slices.Sort(found)

    return found, nil
}

func CollectURL(r io.Reader) ([]string, error) {
    root, err := html.Parse(r)

    if err != nil {
        return nil, err
    }

    //var urls []string
    
    cache := make(map[string]struct{})
    //seen := make(map[*html.Node]struct{})

    var dfs func(*html.Node)

    dfs = func(n *html.Node) {
        if n == nil {
            return
        }
        //seen[n] = struct{}{}
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if url := attr.Val; attr.Key == "href" {
                    //url, _ = normalizeURL(url)
                    if _, ok := cache[url]; !ok {
                        cache[url] = struct{}{}
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            dfs(c)
        }
    }

    dfs(root)
    
    return slices.Collect(maps.Keys(cache)), nil
}

func dfs(n *html.Node, yield func(string) bool) {
    if n == nil {
        return
    }

    if n.Type == html.ElementNode && n.Data == "a" {
        for _, attr := range n.Attr {
            if url := attr.Val; attr.Key == "href" {
                if !yield(url) {
                    return
                }
            }
        }
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        dfs(c, yield)
    }
}

func URLS(root *html.Node) iter.Seq[string] {
    return func(yield func(string) bool) {
        dfs(root, yield)
    }
}

func Unique[T comparable](iterator iter.Seq[T]) iter.Seq[T] {
    seen := make(map[T]struct{})
    return func(yield func(T) bool) {
        for val := range iterator {
            if _, ok := seen[val]; !ok {
                seen[val] = struct{}{}
                if !yield(val) {
                    return
                }
            }
        }
    }
}

type URLCollector func(*html.Node) iter.Seq[string]

func (u URLCollector) Unique(root *html.Node) iter.Seq[string] {
    seen := make(map[string]struct{})
    return func(yield func(string) bool) {
        for val := range u(root) {
            if _, ok := seen[val]; !ok {
                seen[val] = struct{}{}
                if !yield(val) {
                    return
                }
            }
        }
    }

}

func CollectURL2(s string) ([]string, error) {
    root, err := html.Parse(strings.NewReader(s))

    if err != nil {
        return nil, err
    }

    return slices.Collect(Unique(URLS(root))), nil

}

func CollectURL3(s string) ([]string, error) {
    root, err := html.Parse(strings.NewReader(s))

    if err != nil {
        return nil, err
    }

    return slices.Collect(URLCollector.Unique(URLS, root)), nil
}
