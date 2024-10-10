# AWS Vault EC2 and S3 Manager

Welcome to the **AWS Vault EC2 and S3 Manager**, a simple CLI tool to manage your AWS EC2 instances and S3 buckets with AWS credentials securely stored in HashiCorp Vault.

This CLI tool allows you to:

- List and describe EC2 instances
- List S3 buckets
- Securely retrieve AWS credentials from Vault
- Easily extend for more AWS operations in the future

## Features

- **EC2 Management**: List EC2 instances and display their details such as ID, type, and state.
- **S3 Bucket Management**: List all S3 buckets in your account.
- **Vault Integration**: Store and retrieve AWS credentials from Vault, ensuring that sensitive information is kept secure.

## Requirements

- Go 1.16+
- AWS Account
- AWS SDK for Go v2
- HashiCorp Vault setup with AWS credentials stored
- AWS CLI configured locally

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/aws-vault-ec2-s3-manager.git
    cd aws-vault-ec2-s3-manager
    ```

2. Install dependencies (you can use `go mod` if necessary):

    ```bash
    go mod tidy
    ```

3. Build the CLI tool:

    ```bash
    go build -o ec2s3-cli
    ```

4. Run the tool:

    ```bash
    ./ec2s3-cli
    ```

## Usage

Once you run the tool, you will see a menu like this:

```bash
Welcome to EC2 and S3 CLI Manager
Please select an option:
1. Manage EC2 Instances
2. Manage S3 Buckets
3. Exit
