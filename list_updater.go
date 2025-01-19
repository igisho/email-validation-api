package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// updateList retrieves all sources, combines them into one list, removes duplicates and comments, and sorts it.
func updateList() error {
    combinedList := make(map[string]struct{})
    
    for _, source := range config.ListSourceURLs {
        resp, err := http.Get(source)
        if err != nil {
            logError("Error retrieving source:", err)
            continue
        }
        defer resp.Body.Close()

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "'") {
                continue
            }
            combinedList[line] = struct{}{}
        }
    }

    sortedList := make([]string, 0, len(combinedList))
    for line := range combinedList {
        sortedList = append(sortedList, line)
    }
    sort.Strings(sortedList)

    file, err := os.Create(fmt.Sprintf("config/%s", config.ListName))
    if err != nil {
        return fmt.Errorf("error creating target list file: %v", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range sortedList {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return fmt.Errorf("error writing to target list file: %v", err)
        }
    }
    writer.Flush()

    logInfo("List updated successfully")
    return nil
}

// scheduleListUpdate schedules the list update based on the update interval.
func scheduleListUpdate() {
    ticker := time.NewTicker(config.ListUpdateInterval * time.Second)
    go func() {
        for {
            select {
            case <-ticker.C:
                err := updateList()
                if err != nil {
                    logError("Error updating list:", err)
                } else {
                    err = loadAllLists()
                    if err != nil {
                        logError("Error loading lists:", err)
                        return
                    }
                }
            }
        }
    }()
}
