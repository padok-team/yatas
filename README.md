# YATAS
[![codecov](https://codecov.io/gh/StanGirard/YATAS/branch/main/graph/badge.svg?token=OFGny8Za4x)](https://codecov.io/gh/StanGirard/YATAS) [![goreport](https://goreportcard.com/badge/github.com/stangirard/yatas)](https://goreportcard.com/badge/github.com/stangirard/yatas)

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
- `--ci`: Exit code 1 if there are issues found, 0 otherwise.
- `--resume`: Only shows the number of tests passing and failing.

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

## AWS - 47 Checks

### APIGateway
- AWS_APG_001 Apigateway Cloudwatch Logs enabled
- AWS_APG_002 APIGateway stages protected, by ACL

### AutoScaling
- AWS_ASG_001 Autoscaling Desired Capacity vs Max Capacity below 80%

### Backup
- AWS_BAK_001 EC2 Snapshots Encryption
- AWS_BAK_002 EC2 Snapshots Age

### Cloudfront
- AWS_CFT_001 TLS 1.2 Minimum
- AWS_CFT_002 Cloudfront HTTPS Only
- AWS_CFT_003 Standard Logging Enabled
- AWS_CFT_004 Cookies Logging Enabled
- AWS_CFT_005 ACL Used

### CloudTrail
- AWS_CLD_001 Cloudtrails Encryption
- AWS_CLD_002 Cloudtrails Global Service Events Activated
- AWS_CLD_003 Cloudtrails Multi Region

### DynamoDB
- AWS_DYN_001 Dynamodb Encryption
- AWS_DYN_002 Dynamodb Continuous Backups

### EC2
- AWS_EC2_001 EC2 Public IP
- AWS_EC2_002 Monitoring Enabled

### ECR
- AWS_ECR_001 Image Scanning Enabled
- AWS_ECR_002 Encrypted
- AWS_ECR_003 TagImmutable

### LoadBalancer
- AWS_ELB_001 ELB Access Logs Enabled

### GuardDuty
- AWS_GDT_001 GuardDuty Enabled

### IAM
- AWS_IAM_001 IAM 2FA
- AWS_IAM_002 IAM Access Key Age
- AWS_IAM_003 IAM User Can Elevate Rights

### Lambda
- AWS_LMD_001 Lambda Private
- AWS_LMD_002 Lambda In Security Group

### RDS
- AWS_RDS_001 RDS Encryption
- AWS_RDS_002 RDS Backup
- AWS_RDS_003 RDS Minor Auto Upgrade
- AWS_RDS_004 RDS Private
- AWS_RDS_005 RDS Logging
- AWS_RDS_006 RDS Delete Protection

### S3 Bucket
- AWS_S3_001 S3 Encryption
- AWS_S3_002 S3 Bucket in one zone
- AWS_S3_003 S3 Bucket object versioning
- AWS_S3_004 S3 Bucket retention policy
- AWS_S3_005 S3 Public Access Block

### Volume
- AWS_VOL_001 EC2 Volumes Encryption
- AWS_VOL_002 EC2 Volumes Type
- AWS_VOL_003 EC2 Volumes Snapshots

### VPC
- AWS_VPC_001 VPC CIDR
- AWS_VPC_002 VPC Only One
- AWS_VPC_003 VPC Gateway
- AWS_VPC_004 VPC Flow Logs
- AWS_VPC_005 At least 2 subnets
- AWS_VPC_006 Subnets in different zone