package main

import (
    //"io"
    "fmt"
    "os"
    //"text/tabwriter"
    //"slices"
    //"cmp"
    "time"
    "strconv"
)


func main() {
    if n := len(os.Args); n < 3 {
        fmt.Println("not enough arguments")
        os.Exit(1)
    } else if  n > 4 {
        fmt.Println("too many arguments provided")
        os.Exit(1)
    }

    baseURL := os.Args[1]
    maxConcurrency, err  := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    maxPages, _ := strconv.Atoi(os.Args[3])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    con, err := NewConfig(baseURL, maxConcurrency, maxPages, os.Stdout)

    if err != nil {
        fmt.Println("\x1B[38;5;124merror: unable to set up config. got error %v\x1B[0m", err)
    }

    startTime := time.Now()

    con.wg.Add(1)
    con.concurrencyControl <- struct{}{}
    go func() {
        defer con.wg.Done()

        <-con.concurrencyControl

        con.crawlPage(baseURL)

    }()
    con.wg.Wait()

    totalTime := time.Since(startTime)

    fmt.Printf("\x1B[38;5;5mTotal time: %v\x1B[0m\n", totalTime)

    printReport2(con, os.Stdout)

    //fmt.Printf("total sites visited %d\n", con.totalVisited)
}
