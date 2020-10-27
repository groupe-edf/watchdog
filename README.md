# Watchdog - Git server-side custom hooks
[![Actions Status](https://github.com/groupe-edf/watchdog/workflows/test/badge.svg)](https://github.com/groupe-edf/watchdog/actions)
[![Codecov](https://codecov.io/gh/groupe-edf/watchdog/branch/feature/migrate-to-github/graph/badge.svg?token=IDCWJIZ156)](https://github.com/groupe-edf/watchdog)
[![Go Report Card](https://goreportcard.com/badge/github.com/groupe-edf/watchdog)](https://goreportcard.com/report/github.com/groupe-edf/watchdog)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4370/badge)](https://bestpractices.coreinfrastructure.org/projects/4370)

Watchdog allows to define custom hooks in YAML format. When attached to the official repository, some of these can serve as a way to enforce policy by rejecting certain commits or branches.

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)

## Example
``` yaml
version: "1.0.0"
hooks:
  - name: global
    rules:
      - type: commit
        conditions:
        - type: length
          condition: lt 120
          rejection_message: Commit message longer than {{ .Operand }}
      - type: security
        conditions:
        - type: secret
          skip: docs/.*|.*.json|tests/.*
```

## Features

- Accept or reject commits message
- Check files size and extension
- Scan for secrets
- Check branch and tags names
- JSON reporting
- Public git repository scans

## Installation
Please refer to [Documentation](https://groupe-edf.github.io/watchdog/docs/deployment/install)

### From source
```bash
$ go get -ldflags="-X github.com/groupe-edf/watchdog/internal/version.Version=$(cat VERSION)" github.com/groupe-edf/watchdog
```
Or
```bash
$ git clone https://github.com/groupe-edf/watchdog && cd watchdog
$ go install -ldflags="-X github.com/groupe-edf/watchdog/internal/version.Version=$(cat VERSION)"
```

## Usage
First of all, you must create a `.githooks.yml` with some rules. See [Configuration](https://groupe-edf.github.io/watchdog/docs/quickstart/usage)


### Local git repository
```bash
$ watchdog --hook-file=".githooks.yml" \
    --hook-type="" \
    --hook-input="" \
    --logs-level="info" \
    --logs-format="json" \
    --output-format="json" \
    --logs-path="watchdog.log"
```

### Remote git repository
```bash
$ watchdog --hook-file=".githooks.yml" \
    --hook-type="" \
    --hook-input="" \
    --logs-level="info" \
    --logs-format="json" \
    --output-format="json" \
    --logs-path="watchdog.log" \
    --uri="https://github.com/{group}/{repository}"
```

## Roadmap
Watchdog roadmap uses [Github milestones](https://github.com/groupe-edf/watchdog/milestones) to track the progress of the project.

## Contributing
We would love you to contribute to `groupe-edf/watchdog`, pull requests are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## License
The scripts and documentation in this project are released under the [GPL License](LICENSE)