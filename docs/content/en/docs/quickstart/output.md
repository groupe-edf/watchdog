---
title: "Output"
linkTitle: "Output"
weight: 3
description: >
  Issues output and available formatters.
---

When the developer pushes his changes to the remote repository, the watchdog analysis of the various commits will return successfully with a detailed list of low severity issues.

``` bash
remote:
remote: -----BEGIN REJECTION MESSAGES-----
remote: GL-HOOK-ERR: severity=low handler=file condition=extension commit=eda373cc message="'*.exe' files are not allowed"
remote: -----END REJECTION MESSAGES-----
remote:
remote:
remote: #########################################
remote: #                                       #
remote: #  Your push was successfully accepted  #
remote: #                                       #
remote: #########################################
remote:
remote: Operation took 43.319478ms
```
Either fail with a detailed list of high severity issues
``` bash
remote:
remote: -----BEGIN REJECTION MESSAGES-----
remote: GL-HOOK-ERR: severity=high handler=file condition=extension commit=eda373cc message="'*.exe' files are not allowed"
remote: -----END REJECTION MESSAGES-----
remote:
remote:
remote: ####################################################
remote: #                                                  #
remote: #  Your push was rejected because previous errors  #
remote: #                                                  #
remote: ####################################################
remote:
remote: Operation took 43.319478ms
```

Format
--------------------

### logfmt

The default output format of the messages is [logfmt] (https://brandur.org/logfmt) prefixed by `GL-HOOK-ERR:` to have the possibility of uploading these messages on the Gitlab graphical interface in the case of a direct modification of the code on Gitlab.

```bash
severity=high handler=file condition=extension commit=eda373cc message="'*.exe' files are not allowed"
```

* **severity** Issue severity `high`, `medium` or `low`
* **handler** Name of the handler that detected the problem
* **condition** Name of the condition that detected the problem
* **commit** The current hash of the commit (On 8 characters)
* **message** Issue description

Only high severity issues block commits from being persisted in the Git repository.

### json
```bash
[
  {
    "author": "Habib MAALEM",
    "commit": "9560bbeb3b93d9a6d545133dea3e26e0f1fd7a66",
    "condition": "secret",
    "email": "habib.maalem@gmail.com",
    "handler": "security",
    "leaks": [
      {
        "file": "src/main/resources/rsa_server.key",
        "line_number": 1,
        "rule": "ASYMMETRIC_PRIVATE_KEY",
        "severity": "MAJOR",
        "tags": [
          "key"
        ]
      }
    ],
    "message": "Secrets, token and passwords are forbidden, `src/main/resources/rsa_server.key:----***********************`",
    "severity": "low"
  }
]
```