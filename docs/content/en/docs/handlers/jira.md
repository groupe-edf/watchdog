---
title: "Jira Handler"
linkTitle: "Jira Handler"
weight: 9
description: >
  Checks the compliance of commits with projects on Jira
---

## Issue
This rule allows to check if the commit message contains a reference of the associated Jira ticket, this condition can be ignored with the `skip: TECH` attribute.

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: jira
        conditions:
          - type: issue
            skip: TECH
            rejection_message: Commit message is missing the JIRA Issue JIRA-123
```