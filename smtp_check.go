package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// smtpCheck performs the SMTP check using the first MX record with caching.
func smtpCheck(mxRecord, email, domain string, noCache bool) (bool, bool) {
    cacheKey := domain + "-" + email
    if !noCache && config.CacheEnabled {
        if valid, found := cache.GetSMTP(cacheKey); found {
            return valid, true
        }
    }

    timeout := 10 * time.Second
    conn, err := net.DialTimeout("tcp", net.JoinHostPort(mxRecord, "25"), timeout)
    if err != nil {
        logError("Failed to connect to MX record", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }
    defer conn.Close()

    // Read server's initial response.
    buff := make([]byte, 1024)
    n, err := conn.Read(buff)
    //logInfo("SMTP-Response for email: ", email, string(buff[:n]))

	// 554 Your access to this mail system has been rejected due to the sending MTA's poor reputation. If you believe that this failure is in error, please contact the intended recipient via alternate means.
	if strings.HasPrefix(string(buff[:n]), "554") {
        logError("Your access to this mail system has been rejected due to the sending MTA's poor reputation\n")
        cache.SetSMTP(cacheKey, true)
        return true, false
    }


    if err != nil || !strings.HasPrefix(string(buff[:n]), "220") {
        logError("Initial response error", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }

	
    // Send HELO command.
    heloCmd := fmt.Sprintf("HELO %s\r\n", domain)
    if _, err := conn.Write([]byte(heloCmd)); err != nil {
        fmt.Printf("HELO command error: %v\n", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }

    // Send MAIL FROM command.
    mailFromCmd := fmt.Sprintf("MAIL FROM:<%s>\r\n", config.SMTPEmail)
    if _, err := conn.Write([]byte(mailFromCmd)); err != nil {
        fmt.Printf("MAIL FROM command error: %v\n", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }

    // Send RCPT TO command.
    rcptToCmd := fmt.Sprintf("RCPT TO:<%s>\r\n", email)
    if _, err := conn.Write([]byte(rcptToCmd)); err != nil {
        fmt.Printf("RCPT TO command error: %v\n", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }


    // Read response to RCPT TO command.
    n, err = conn.Read(buff)

        logInfo("SMTP-Response", email, string(buff[:n]))

    if err != nil || !strings.HasPrefix(string(buff[:n]), "250") {
        logError("RCPT TO response error", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }

    // Send QUIT command.
    if _, err := conn.Write([]byte("QUIT\r\n")); err != nil {
        logError("QUIT command error", err)
        cache.SetSMTP(cacheKey, false)
        return false, false
    }

    logInfo("SMTP check successful for email:", email)
    cache.SetSMTP(cacheKey, true)
    return true, false
}
