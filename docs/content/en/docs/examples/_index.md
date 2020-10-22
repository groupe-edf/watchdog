---
title: "Examples"
linkTitle: "Examples"
weight: 3
---

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: commit
        conditions:
          - type: length
            condition: gt 20
            rejection_message: Commit message should exceed {{ .Operand }} characters
```