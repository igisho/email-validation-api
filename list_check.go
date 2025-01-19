package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var listCache = struct {
    sync.RWMutex
    data map[string][]string
}{data: make(map[string][]string)}


// clearListCache safely clears the listCache data.
func clearListCache() {
    listCache.Lock()
    defer listCache.Unlock()
    listCache.data = make(map[string][]string)
}


// loadAllLists loads all lists from the file system into memory at startup.
func loadAllLists() error {
    files, err := os.ReadDir("config")
    if err != nil {
        return fmt.Errorf("error reading config directory: %v", err)
    }

    // Sort files by order
    sort.Slice(files, func(i, j int) bool {
        return files[i].Name() < files[j].Name()
    })

    clearListCache()

    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
            continue
        }

        filename := file.Name()
        parts := strings.Split(filename, "-")
        if len(parts) < 3 {
            continue
        }

        list, err := loadList(fmt.Sprintf("config/%s", filename))
        if err != nil {
            logError("Error loading list file:", err)
            continue
        }

        listCache.Lock()
        listCache.data[filename] = list
        listCache.Unlock()
    }

    return nil
}

// loadList loads a list from file and caches it in memory.
func loadList(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var list []string
    for scanner.Scan() {
        list = append(list, scanner.Text())
    }

    return list, nil
}

// checkList checks if the email, name, or domain is in the allow or deny list based on the configured order.
func checkList(email string) (bool, string) {
    email = strings.ToLower(email)
    parts := strings.Split(email, "@")
    name := regexp.MustCompile("[^a-z]").ReplaceAllString(parts[0], "")
    domain := parts[1]

    listCache.RLock()
    defer listCache.RUnlock()

    for filename, list := range listCache.data {
        parts := strings.Split(filename, "-")
        if len(parts) < 3 {
            continue
        }

        action := parts[1]
        checkType := strings.TrimSuffix(parts[2], ".txt")

        for _, line := range list {
            switch checkType {
            case "email":
                if line == email {
                    return action == "allow", fmt.Sprintf("Email %s is %s by list", email, action)
                }
            case "name":
                if line == name {
                    return action == "allow", fmt.Sprintf("Name %s is %s by list", name, action)
                }
            case "domain":
                if line == domain {
                    return action == "allow", fmt.Sprintf("Domain %s is %s by list", domain, action)
                }
            }
        }
    }

    return false, ""
}
