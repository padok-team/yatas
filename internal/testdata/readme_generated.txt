
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

### LoadBalancer
- AWS_ELB_001 ELB Access Logs Enabled

### GD
- AWS_GD_001 GuardDuty Enabled

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

### SecurityHub
- AWS_SHU_001 SecurityHub Activated
- AWS_SHU_002 SecurityHub Auto Enabled

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