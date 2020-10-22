---
title: "Commit Handler"
linkTitle: "Commit Handler"
weight: 9
description: >
  Check commits messages.
---

## Regular Expression
This is a generic rule where you can use regular expressions to define conditions to accept or reject commit messages. See [Syntax](https://github.com/google/re2/wiki/Syntax).

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: commit
        conditions:
          - type: pattern
            condition: (?m)^(build|ci|docs|feat|fix|perf|refactor|style|test)\([a-z]+\):\s([a-z\.\-\s]+)
            skip: Merge commit
```

## Commit message length
This rule forces users not to exceed a certain length of the commit message.

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: commit
        conditions:
          - type: length
            condition: le 20
            rejection_message: Commit message should not exceed '{{ .Operand }}' characters
```

In the `condition` attribute we can use one of the following predicates
| Predicate          | Description        |
|--------------------|--------------------|
| **eq** | eqal |
| **ne** | not equal |
| **lt** | lower than |
| **le** | lower than or equal |
| **ge** | greater than or equal |
| **gt** | greater than |

## Email
It is a generic rule of type regular expression which allows to define conditions on the emails of the users to accept or refuse the commits. See [Syntax](https://github.com/google/re2/wiki/Syntax).

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: commit
        conditions:
          - type: email
            condition: "[a-zA-Z\-]+@acme.com"
            rejection_message: Author email '{{ .Commit.Author.Email }}' is not valid email address
```