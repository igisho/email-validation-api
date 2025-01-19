package main

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents the configuration structure.
type Config struct {
    SMTPEmail            string        `json:"smtp_email"`
    Regex                string        `json:"regex"`
    CacheEnabled         bool          `json:"cache_enabled"`
    CacheMaxTimeMX       time.Duration `json:"cache_max_time_mx"`
    CacheMaxTimeSMTP     time.Duration `json:"cache_max_time_smtp"`
    LoggingEnabled       bool          `json:"logging_enabled"`
    APIKey               string        `json:"api_key"`
    RequestsPerMinute    int           `json:"requests_per_minute"`
    Fail2BanDuration     time.Duration `json:"fail2ban_duration"`
    Fail2BanAllow        []string      `json:"fail2ban_allow"`    
    ListName             string        `json:"list_name"`
    ListUpdateInterval   time.Duration `json:"list_update_interval"`
    ListSourceURLs       []string      `json:"list_source_urls"`
}

// Global variable to hold the configuration
var config Config

// loadConfig loads the configuration from the config.json file.
func loadConfig() error {
    file, err := os.Open("config/config.json")
    if err != nil {
        return err
    }
    defer file.Close()
    return json.NewDecoder(file).Decode(&config)
}
