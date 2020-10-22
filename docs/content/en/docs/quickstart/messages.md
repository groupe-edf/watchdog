---
title: "Messages"
linkTitle: "Messages"
weight: 2
description: >
  Checks the security of content pushed by developers.
---

The rejection messages are written in the form of a template that can be enriched with the already predefined substitution variables

Variables
--------------------
| Variable           | Description        |
|--------------------|--------------------|
| **{{ .Branch }}** | The name of the branch. |
| **{{ .Commit.Author.Email }}** | Email of the author of the commit. |
| **{{ .Commit.Author.Name }}** | The name of the author of the commit. |
| **{{ .Commit.Subject }}** | The commit message. |
| **{{ .Commit.Date }}** | The date of the commit. |
| **{{ .Commit.Hash }}** | The hash of the commit. |
| **{{ .Commit.ShortHash }}** | The 8 character hash of the commit. |
| **{{ .Condition.Condition }}** | La condition définit. |
| **{{ .Object }}** | It can have several values, for example the name of the scanned file in the manager **File**. |
| **{{ .Operator }}** | The operator used in the condition. |
| **{{ .Operand }}** | The comparison value. |
| **{{ .Tag }}** | The name of the tag. |
| **{{ .Value }}** | It can have several values, for example the size of the file in the **File** manager or confidential information detected in the **Security** manager. |

**Example**
```yaml
rejection_message: File {{ .Object }} size {{ .Value }} greater or equal than {{ .Operand }}
```