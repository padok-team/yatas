# YATAS
Yet Another Testing &amp; Auditing Solution 

## Features
YATAS is a simple and easy to use tool to audit your infrastructure for misconfiguration or potential security issues.

## Installation

```bash
brew tap stangirard/tap
brew install yatas
```

```bash
cp .yatas.yml.example .yatas.yml
```

Modify .yatas.yml to your needs.

## Usage

```bash
yatas ## --details 
```

Flags:
- `--details`: Show details of the issues found.

## Plugins

| Name | Description | Checks |
|------|-------------|--------|
| *AWS* | AWS checks | EC2 and S3 encryption |

