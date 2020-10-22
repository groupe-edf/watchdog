---
title: "File Handler"
linkTitle: "File Handler"
weight: 9
description: >
  Checks files pushed by developers (extensions, size ...).
---

## Extension
This rule allows developers to be prohibited from committing files that have certain extensions.

``` yaml
version: "1.0.0"
hooks:
  - name: files
    type: pre-receive
    rules:
      - type: file
        conditions:
          - type: extension
            rejection_message: '*.{{ .Condition.Condition }}' files are not allowed
            condition: exe
```

## Size
This rule is used to prevent developers from committing files that exceed the size defined in the condition.

``` yaml
version: "1.0.0"
hooks:
  - name: files
    rules:
      - type: file
        conditions:
          - type: size
            condition: lt 5mb
            rejection_message: File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}
```

Currently only the `lt` operator is supported. the size is defined in natural language and in English and in lowercase, example:
* `lt 5mb`
* `lt 5 megabytes`
* `lt 512 kb`

To fix this, move media files (.mp4, .mp3, .jpg, .png) from repositories that are only “supposed” to contain text to a repository hold media files that humans can’t read. The job of (git-lfs (Large File System))[https://wilsonmar.github.io/git-hooks/#git-lfs] is to move and replace binary files with a (texual) link to binary repositories.