---
title: "Conventional Commits"
linkTitle: "Conventional Commits"
weight: 1
---

A good commit message should tell you what has changed and why. The how, that is, how to make these changes, does not have to be explained. Reading the code and highlighting changes via a diff is self-explanatory. For example checks if the commit messages respect the conventional validation format: https://www.conventionalcommits.org/en/v1.0.0/

```yaml
version: "1.1.0"
hooks:
  - name: generic
    rules:
      - type: commit
        conditions:
          - type: length
            condition: le 120
      - type: commit
        conditions:
          # Angular guidelines https://github.com/angular/angular/blob/master/CONTRIBUTING.md
          - type: pattern
            condition: (?m)^(build|ci|docs|feat|fix|perf|refactor|style|test)\([a-z]+\):\s([a-z\.\-\s]+)
            rejection_message: >
              Message must be formatted like type(scope): subject, type Must be one of the following:
              * build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
              * ci: Changes to our CI configuration files and scripts (example scopes: Circle, BrowserStack, SauceLabs)
              * docs: Documentation only changes
              * feat: A new feature
              * fix: A bug fix
              * perf: A code change that improves performance
              * refactor: A code change that neither fixes a bug nor adds a feature
              * style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
              * test: Adding missing tests or correcting existing tests
```