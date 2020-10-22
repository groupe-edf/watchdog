---
title: "Branch Handler"
linkTitle: "Branch Handler"
weight: 9
description: >
  Check all Git actions on the project branches.
---

## Regular Expression
This is a generic rule where you can use regular expressions to define the naming conventions for branches. See [Syntax](https://github.com/google/re2/wiki/Syntax)
``` yaml
version: "1.0.0"
hooks:
  - name: global
    type: pre-receive
    rules:
      - type: branch
        conditions:
          - type: pattern
            condition: (feature|release|hotfix)\/[a-z\d-_.]+
            rejection_message: Branch `{{ .Branch }}` must match Gitflow naming convention
```