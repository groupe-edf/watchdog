---
title: "Tag Handler"
linkTitle: "Tag Handler"
weight: 9
description: >
  Check all Git actions on tags.
---

## SemVer
This rule ensures that Tags pushed by developers respect the naming conventions described by `Semantic Versioning` V2.0.0. See [https://semver.org/](https://semver.org/).

``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
    - type: tag
      conditions:
      - type: semver
        rejection_message: Tag version `{{ .Tag }}` must respect semantic versionning v2.0.0 https://semver.org/
```