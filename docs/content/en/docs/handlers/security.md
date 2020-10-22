---
title: "Security Handler"
linkTitle: "Security Handler"
weight: 9
description: >
  Checks the security of content pushed by developers.
---

## Secrets
This rule allows you to analyze the presence of sensitive and confidential information in the content of commits files, currently the rule allows you to scan.

* AWS Access Key
* Authentification basique
* Google Access Token
* URL de connexion MySQL, Redis, PostgreSQL
* Private Key

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
    - type: security
      conditions:
      - type: secret
        rejection_message: Secrets, token and passwords are forbidden, `{{ .Object }}:{{ .Value }}`
        skip: .*.json|tests/.*
```

Providers
--------------------
List of supported providers

| Domain | Platform/API | Key Type  | Target Regular Expression | Source |
|:----------|:----------|:----------|:----------|:----------|
| Cloud | Amazon Web Services | Access Key ID | `AKIA[0-9A-Z]{16}` |  |
| Cloud | Amazon Web Services | Secret Key | `[0-9a-zA-Z/+]{40}` |  |
| Web | URI | Basic Authentication | `(http\|https)://[^{}[[:space:]]]+:([^{}[[:space:]]]+)@` | [The 'Basic' HTTP Authentication Scheme](https://tools.ietf.org/html/rfc7617) |
| Web | NPM | Base64 | `_auth[[:space:]]*=[[:space:]]*(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==\|[A-Za-z0-9+\\/]{3}=\|[A-Za-z0-9+\\/]{4})` | [The Base16, Base32, and Base64 Data Encodings](https://tools.ietf.org/html/rfc3548) |
| Cloud | Google Cloud Platform | OAuth 2.0 | `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}` |  |
| Cloud | Google Cloud Platform | API Key | `[A-Za-z0-9_]{21}--[A-Za-z0-9_]{8}` |  |
| Database | MySQL | Basic Authentication | `mysql://[^{}[[:space:]]]+:([^{}\[[:space:]]]+)@` |  |
| Database | Redis | Basic Authentication | `(redis\|rediss\|redis-socket\|redis-sentinel)://[^{}[[:space:]]]+:([^{}[[:space:]]]+)@` |  |
| Communication | Slack | API Key | `xox.-[0-9]{12}-[0-9]{12}-[0-9a-zA-Z]{24}` |  |
| Communication | Twilio | API Key | `55[0-9a-fA-F]{32}` |  |
| Social Media | Twitter | Access Token | `[1-9][0-9]+-[0-9a-zA-Z]{40}` |  |