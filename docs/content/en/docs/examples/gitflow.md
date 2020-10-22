---
title: "Gitflow"
linkTitle: "Gitflow"
weight: 1
---

```yaml
version: "1.1.0"
hooks:
  - name: generic
    type: pre-receive
    rules:
      - type: branch
        conditions:
          - type: pattern
            condition: (feature|release|hotfix)\/[a-z\d-_.]+
            rejection_message: Branch name must match gitflow naming
```