# YATAS
Yet Another Testing &amp; Auditing Solution 

## Features
YATAS is a simple and easy to use tool to audit your infrastructure for misconfiguration or potential security issues.

<p align="center">
<img src="docs/demo.png" alt="demo" width="60%">
<p align="center">

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
| *AWS* | AWS checks | Good practices and security checks|


## Checks 

### AWS

- Check if S3 encryption is enabled
- Check if EC2 encryption is enabled
- Check if RDS encryption is enabled
- Check if RDS backup is enabled
- Check if RDS auto minor upgrade is enabled
-  Check if RDS private is enabled
-  Check if VPC CIDR is /20 or bigger
-  Check if VPC Flow Logs are enabled

