package main

import (
	"fmt"
	"regexp"
	"strings"
)

// validateEmailFormat validates the email format using regex from the config.
func validateEmailFormat(email string) bool {
    return regexp.MustCompile(config.Regex).MatchString(email)
}

// validateEmail performs all email validation checks.
func validateEmail(email string, noCache bool) (bool, string, string, bool) {
    // Check email format.
    if !validateEmailFormat(email) {
        logError("Invalid email format for email:", email)
        return false, "Invalid email format", "Email format does not match standard", false
    }

    // Check against allow/deny lists.
    valid, message := checkList(email)
    if valid {
        logInfo("Email is valid by list for email:", email)
        return true, message, "", false
    } else if message != "" {
        logError("Email is denied by list for email:", email)
        return false, message, "Email is denied by list", false
    }

    domain := strings.Split(email, "@")[1]

    // Get first MX record.
    mxRecords, mxCached, err := getFirstMXRecord(domain, noCache)
    if err != nil {
        logError("No MX record found for domain:", domain)
        return false, err.Error(), "Domain has no mail server", false
    }

    // SMTP check.
    smtpValid, smtpCached := smtpCheck(mxRecords[0], email, domain, noCache)
    if !smtpValid {
        logError("Failed to verify SMTP server for domain:", domain)
        return false, fmt.Sprintf("Failed to verify SMTP server for domain %s", domain), "SMTP server not responsive", smtpCached
    }

    logInfo("Email validation successful for email:", email)
    return true, fmt.Sprintf("Email %s is valid", email), "", smtpCached || mxCached
}
