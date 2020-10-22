---
title: "Gitlab"
linkTitle: "Gitlab"
weight: 3
description: >
  Deploy Watchdog in Gitlab instance.
---

## Usage server-side

Watchdog can be configured only in the server side of the Git repository (For the moment a **bare** repository), it is compatible with different Git solutions like [Gitea](https://gitea.io/), [Gogs](https://gogs.io/), [Gitlab](https://gitlab.com/).

In our case we will focus on installing Watchdog on the Gitlab platform. See [instructions](https://docs.gitlab.com/ee/administration/server_hooks.html).

We have the possibility to install the tool in different levels of Gitlab
* Global
* By group
* By project

For a global deployment, you must create one of the `pre-receive.d`,` post-receive.d` or `update.d` folders in the `/opt/gitlab/embedded/service/gitlab-shell/hooks` directory and add a script that will launch note CLI

``` bash
#!/bin/sh
HOOKS=("pre-receive" "post-receive" "update")
HOOK_TYPE=$(cd $(dirname "${BASH_SOURCE[0]}") >/dev/null 2>&1 && echo ${PWD##*/})
while read -r OLDREV NEWREV REFNAME; do
  /usr/local/bin/watchdog \
    --docs-link="https://groupe-edf.github.io/watchdog/docs/" \
    --hook-type="pre-receive" \
    --hook-input="$OLDREV $NEWREV $REFNAME" \
    --logs-path="watchdog.log" \
    --logs-level="info" \
    --verbose="true"
  status=$?
  if [ $status -eq 1 ]; then
    exit 1
  fi
done
```

Options
--------------------

| Option            | Description                      |
|-------------------|----------------------------------|
| `--docs-link`     | Allows you to define the link to the Custom Hooks documentation           |
| `--hook-input`    | This data is automatically filled in upon receipt of a Git `git push` package        |
| `--hook-type`     | Allows you to specify to the watchdog CLI the type of hooks it has just received to pass it to the different managers. By default `pre-receive`     |
| `--logs-path`     | Allows you to define the path for writing logs. By default it is `/var/log/watchdog/watchdog.log` |
| `--logs-level`    | Allows to define the level of logs to write, can take one of the following values: `debug`, `error`, `fatal`, `info`, `panic`, `trace` and `warning`. By default it is `info`        |
| `--verbose`       | Allows you to view more information about running Custom Hooks in the Developer Console. Default `true`          |
| `--help`, `-h`    | Show help/usage                  |
| `--version`, `-v` | Print version                    |

The `WATCHDOG_*` environment variables are used to pass the adjustment parameters of certain services