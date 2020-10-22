---
title: "Usage"
linkTitle: "Usage"
weight: 1
description: >
  Usage.
---

## Usage client-side

On the developer side, to invoke customs hooks, the user must add a configuration file under the name of `.githooks.yml` or `.githook.yaml` in the root of the project.

``` bash
$ ll /workspaces/project-name
  |--docs/
  |--src/
  |--tests/
  |--.gitignore
  |--.githooks.yml
```

.githooks.yml
--------------------

In this file the developer or rather the maintainer of the project will define the different rules that will be executed in the remote Git repository. The following table describes the schema and the different attributes of the file.

| Attribut | Description | Default | Optional | Type |
|:--------------------|:--------------------|:----------:|:----------:|----------:|
| **version** | Version of the .githooks file, the tool will allow you to accept different versions of the diagram |  |  | String |
| **hooks[]** | Liste des hooks à déclencher |  |  | String |
| **hook.name** | Nom du hook |  |  | String |
| **hook.type** | Git hook type, it can have one of the values `pre-receive`, `update`, `post-receive`. See [Customizing Git - Git Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) | pre-receive |  | String |
| **hook.rejection_message** | Global rejection message |  | Yes | String |
| **hook.rules[]** | List of custom hook rules |  |  | String |
| **hook.rule.type** | Type of the rule to be triggered, it can have one of the following values `branch`, `commit`, `file`, `jira`, `security`. See [Handlers](#handlers) |  |  | String |
| **hook.rule.enabled** | This attribute enables or disables the rule | true | Yes | String |
| **hook.rule.conditions[]** | List of conditions that will trigger the rule. They are specific to each rule |  |  | String |
| **hook.rule.condition.type** | Type of condition, the value of this attribute may vary from one rule to another |  |  | String |
| **hook.rule.condition.pattern** | It is a regular expression that describes the condition to be satisfied. See [Syntaxe](https://github.com/google/re2/wiki/Syntax) |  |  | String |
| **hook.rule.condition.rejection_message** | Rule-specific rejection message. By default all customs hooks offer rejection messages. See [Acceptance and rejection messages](#issues) |  | Yes | String |
| **hook.rule.condition.ignore** | This option allows you to ignore not to block the push even if the rule is not satisfied | false | Yes | String |
| **hook.rule.condition.skip** | A regular expression that allows the rule to be ignored if it is satisfied. See [Syntaxe](https://github.com/google/re2/wiki/Syntax). Useful in the case of automatic commits generated during Merge Requests |  | Yes | String |