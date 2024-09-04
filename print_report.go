package main

import (
    "cmp"
    "io"
    "text/tabwriter"
    "fmt"
    "slices"
)


type mapEntry struct {
    key string
    value int
}

func printResults(c *config, w io.Writer) {
    entries := getEntrySlice(c.pages)

    writer := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

    fmt.Fprintln(writer, "\x1B[38;5;207mpage\tvisited\x1B[0m")

    for _, entry := range entries {
        fmt.Fprintf(writer, "%s\t%d\n", entry.key, entry.value)
    }

    writer.Flush()
}

func totalVisited(entries []mapEntry) int {
    total := 0
    for _, entry := range entries {
        total += entry.value
    }
    return total
}

func getEntrySlice(m map[string]int) []mapEntry {
    entries := make([]mapEntry, 0, len(m))
    for key, value := range m {
        entries = append(entries, mapEntry{ key, value })
    }

    sortEntriesVisited(entries)
    
    slices.Reverse(entries)

    return entries
}

func sortEntriesVisited(entries []mapEntry) {
    slices.SortFunc(entries, func(a, b mapEntry) int {
        return cmp.Compare(a.value, b.value)
    })
}

func printReport(pages map[string]int, baseURL string) {
    fmt.Println("=============================")
    fmt.Printf("REPORT for %s\n", baseURL)
    fmt.Println("=============================")

    entries := getEntrySlice(pages)

    for _, entry := range entries {
        fmt.Printf("Found %d internal links to %s\n", entry.value, entry.key)
    }

}

func printReport2(c *config, w io.Writer) {
    fmt.Fprintln(w, "=============================")
    s := c.baseURL.String()
    fmt.Fprintf(w, "REPORT for %s\n", s)
    fmt.Fprintln(w, "=============================")

    entries := getEntrySlice(c.pages)

    for _, entry := range entries {
        fmt.Fprintf(w, "Found %d internal links to %s\n", entry.value, entry.key)
    }


}
