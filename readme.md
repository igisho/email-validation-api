# Email Validation API (EVA)

This is a Go-based API for validating email addresses. The API includes several validation checks and security features such as rate limiting, fail2ban-like functionality, and API key authentication. The application is Dockerized for easy deployment.


## Example

GET /name@domain.tld

`
{
email: "name@domain.tld",
valid: true,
message: "Email name@domain.tld is valid",
cached: false
}
`

## Features

- **Email Format Validation**: Checks if the email format is valid using a regex.
- **Allow/Deny Lists**: Validates emails, names, or domains against allow or deny lists.
- **Disposal Check**: Check against open source lists of disposable email services.
- **MX Record Check**: Checks if the domain has valid MX records.
- **SMTP Check**: Verifies the SMTP server of the domain.
- **Security Headers**: Adds security headers to responses.
- **API Key Authentication**: Protects the API with an API key.
- **Rate Limiting**: Limits the number of requests a client can make.
- **Fail2ban-like Functionality**: Temporarily bans IP addresses with too many failed login attempts.
- **Dockerized**: Easy deployment with Docker and Docker Compose.

## Endpoints

### Validate Email Address (GET)
- Method: GET
- Path: <name@domain.tld> (or as query param)
- Query
  - email (string): The email address to validate (if not provided as path)
  - nochache (boolean): If "true", bypass the cache.
  - key(string): (Optional) The API key for authentication.

### Validate Email Address (POST)
- Method: POST
- Path: <name@domain.tld> (or as form data)
- Form Data / Payload
  - email (string): The email address to validate (if not provided as path)
  - nochache (boolean): If "true", bypass the cache.
  - key(string): (Optional) The API key for authentication.


## Installation

git clone https://github.com/gf78/email-validation-api.git
cd email-validation-api
chmod +x build.sh
./build.sh
docker-compose logs -f

## Configuration

The configuration is managed through a JSON file located at `config/config.json`.
