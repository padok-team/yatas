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
- `--compare`: Compare the results of the previous run with the current run and show the differences.

## Plugins

| Name | Description | Checks |
|------|-------------|--------|
| *AWS* | AWS checks | Good practices and security checks|


## Checks 

### Ignore results for known issues
You can ignore results of checks by add the following to your `.yatas.yml` file:

```yaml
ignore:
  - id: "AWS_VPC_004"
    regex: true
    values: 
      - "VPC Flow Logs are not enabled on vpc-.*"
  - id: "AWS_VPC_003"
    regex: false
    values: 
      - "VPC has only one gateway on vpc-08ffec87e034a8953"
```

### Exclude a test
You can exclude a test by adding the following to your `.yatas.yml` file:

```yaml
plugins:
  - name: "aws"
    enabled: true
    description: "Check for AWS good practices"
    exclude:
      - AWS_S3_001
```

### Specify which tests to run 

To only run a specific test, add the following to your `.yatas.yml` file:

```yaml
plugins:
  - name: "aws"
    enabled: true
    description: "Check for AWS good practices"
    include:
      - "AWS_VPC_003"
      - "AWS_VPC_004"
```

### AWS - 25 Checks

- AWS_S3_001 S3 Encryption
- AWS_S3_002 S3 Bucket in one zone
- AWS_VOL_001 EC2 Volumes Encryption
- AWS_VOL_002 EC2 Volumes Type
- AWS_VOL_003 EC2 Volumes Encryption
- AWS_VOL_003 EC2 Volumes Snapshots
- AWS_VOL_003 EC2 Volumes Type
- AWS_VOL_003 EC2 Snapshots Encryption
- AWS_VOL_004 EC2 Snapshots Age
- AWS_RDS_001 RDS Encryption
- AWS_RDS_002 RDS Backup
- AWS_RDS_003 RDS Minor Auto Upgrade
- AWS_RDS_004 RDS Private
- AWS_VPC_001 VPC CIDR
- AWS_VPC_002 VPC Only One
- AWS_VPC_003 VPC Gateway
- AWS_VPC_004 VPC Flow Logs
- AWS_CLD_001 Cloudtrails Encryption
- AWS_CLD_002 Cloudtrails Global Service Events Activated
- AWS_CLD_003 Cloudtrails Multi Region
- AWS_CLD_001 Image Scanning Enabled
- AWS_LMD_001 Lambda Private
- AWS_LMD_002 Lambda In Security Group
- AWS_DYN_001 Dynamodb Encryption
- AWS_DYN_002 Dynamodb Continuous Backups
