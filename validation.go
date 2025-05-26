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
func validateEmail(email string, noCache bool) (bool, string, string, bool, map[string]bool) {
    checks := make(map[string]bool)
    cached := false

    // 1. Syntax check (regex)
    checks["syntax"] = validateEmailFormat(email)
    if !checks["syntax"] {
        logError("Invalid email format for email:", email)
        return false, "Invalid email format", "Email format does not match standard", cached, checks
    }

    // 2. Allow/Deny list
    valid, message := checkList(email)
    if valid {
        checks["allow_deny_list"] = true
        logInfo("Email is valid by list for email:", email)
        return true, message, "", cached, checks
    } else if message != "" {
        checks["allow_deny_list"] = false
        logError("Email is denied by list for email:", email)
        return false, message, "Email is denied by list", cached, checks
    } else {
        // Neither explicitly allowed nor denied
        checks["allow_deny_list"] = true
    }

    // 3. MX records
    domain := strings.Split(email, "@")[1]
    mxRecords, mxCached, err := getFirstMXRecord(domain, noCache)
    checks["mx"] = err == nil
    if err != nil {
        logError("No MX record found for domain:", domain)
        return false, err.Error(), "Domain has no mail server", mxCached, checks
    }

    // 4. SMTP check
    smtpValid, smtpCached := smtpCheck(mxRecords[0], email, domain, noCache)
    checks["smtp"] = smtpValid
    if !smtpValid {
        logError("Failed to verify SMTP server for domain:", domain)
        return false, fmt.Sprintf("Failed to verify SMTP server for domain %s", domain), "SMTP server not responsive", smtpCached || mxCached, checks
    }

    logInfo("Email validation successful for email:", email)
    return true, fmt.Sprintf("Email %s is valid", email), "", smtpCached || mxCached, checks
}
