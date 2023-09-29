# Changelog
All notable changes to this add-on will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## Unreleased

### Added
- Max scan duration support.

### Fixed
- `ScanScopeCoverage` not picking from env bug

### Changed
- Excluded HTTP and HTTPS services from port-scanner results.
- Fixed passing whole URL to port-scanner for main domain instead of just host.

## [0.0.5] - 2023-09-26
- Docker image tag: `0.0.5-ca65421`
#### Fixed
- `?for_complete=true` added in stop api call
- `defer` function in main.go 

#### Added
- Support for sending start scan request

## [0.0.4] - 2023-08-18
- Docker image tag: `0.0.4-68c1832`
#### Fixed
- Don't panic if `furious` command exited with error, instead print the error.
#### Added
- Health update on webhook

## [0.0.3] - 2023-08-07
- Docker image tag: `0.0.3-7d578c6`
#### Added
- Send all subdomains found on a custom webhook on `subdomains.found`
- Run port scanning on all the sub-domains is `scopeCoverage` allows.

#### Removed
- Removed the support for `ioutil` as it is a deprecated package

## [0.0.2] - 2023-07-20
- Docker image tag: `0.0.2-b163ef5`
#### Added
- Initializes subdomain takeover.

## [0.0.1] - 2023-06-30
- Docker image tag: `0.0.1-d3bb22b`
#### Added
- Initial support for port scanner added
- Send alert to slack if any panic alert is faced
- Send start scan and end scan updates to dast api server