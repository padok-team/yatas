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

- AWS_S3_001  S3 Encryption
- AWS_S3_002  S3 Account Zone Only
- AWS_VOL_001 EC2 Volumes Encryption
- AWS_RDS_001 RDS Encryption
- AWS_RDS_002 RDS Backup
- AWS_RDS_003 RDS Minor Auto Upgrade
- AWS_RDS_004 RDS Private
- AWS_VPC_001 VPC CIDR
- AWS_VPC_002 VPC Flow Logs
- AWS_VPC_003 VPC Gateway
- AWS_VPC_004 VPC Only One

