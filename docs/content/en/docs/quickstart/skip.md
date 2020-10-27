---
title: "Skip"
linkTitle: "Skip"
weight: 4
description: >
  Skiping all custom hooks, certain handlers or conditions
---

Developers have the option of skipping checks on certain commits with the expression `[skip hooks.rule.condition]`. The expression is valid at

* Rule, for example: Ignore security rules of type `Secret` **[skip hooks.security.secret]**
* Handler, for example: Skip all `security` rules **[skip hooks.security]**
* All Handlers **[skip hooks]**

Example
----------
```bash
$ git commit -am "chore: add dummy configuration file [skip hooks.security.secret]"
$ git commit -am "chore: add dummy configuration file [skip hooks.security]"
$ git commit -am "chore: add dummy configuration file [skip hooks]"
```

Or even better with [Push Options](https://git-scm.com/docs/git-push#Documentation/git-push.txt--oltoptiongt) to keep git history clean **Under development**
```bash
$ git push -o security.secrets.skip
$ git push -o hooks.skip="security.secrets,tag.semver" # Comma separated handler.rule
```