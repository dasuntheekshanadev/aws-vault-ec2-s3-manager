package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
)

func main() {
	color.Cyan("Welcome to EC2 and S3 CLI Manager")

	// Main loop for the menu
	for {
		displayMenu()
		choice := getUserInput()

		switch choice {
		case "1":
			handleEC2()
		case "2":
			handleS3()
		case "3":
			color.Yellow("Exiting the program. Goodbye!")
			os.Exit(0)
		default:
			color.Red("Invalid option. Please try again.")
		}
	}
}

// Display menu options
func displayMenu() {
	color.Cyan("Please select an option:")
	color.Cyan("1. Manage EC2 Instances")
	color.Cyan("2. Manage S3 Buckets")
	color.Cyan("3. Exit")
}

// Get user input from the terminal
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1] // Remove the newline character from input
}

// Handle EC2 management
func handleEC2() {
	awsCreds, err := getAWSCredentialsFromVault()
	if err != nil {
		log.Fatalf("Failed to get credentials: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(awsCreds))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	client := ec2.NewFromConfig(cfg)
	describeEC2Instances(client)
}

// Handle S3 bucket management
func handleS3() {
	awsCreds, err := getAWSCredentialsFromVault()
	if err != nil {
		log.Fatalf("Failed to get credentials: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(awsCreds))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	listS3Buckets(client)
}

// Fetches AWS credentials from Vault (same as in your original code)
func getAWSCredentialsFromVault() (aws.CredentialsProvider, error) {
	vaultClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %v", err)
	}

	secret, err := vaultClient.Logical().Read("secret/data/aws")
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from Vault: %v", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path 'secret/data/aws'")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format in 'secret/data/aws'")
	}

	creds, ok := data["creds"].(string)
	if !ok || creds == "" {
		return nil, fmt.Errorf("failed to find creds field in secret")
	}

	var credentials map[string]string
	err = json.Unmarshal([]byte(creds), &credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to parse creds JSON: %v", err)
	}

	accessKey := credentials["AWS_ACCESS_KEY_ID"]
	secretKey := credentials["AWS_SECRET_ACCESS_KEY"]

	return CustomCredentialsProvider{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
	}, nil
}

// Describe EC2 instances
func describeEC2Instances(client *ec2.Client) {
	color.Green("Describing EC2 instances...")

	output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatalf("Failed to describe EC2 instances: %v", err)
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			color.Green("Instance ID: %s", *instance.InstanceId)
			color.Green("Instance State: %s", instance.State.Name)
			color.Green("Instance Type: %s", instance.InstanceType)
		}
	}
}

// List S3 buckets
func listS3Buckets(client *s3.Client) {
	color.Blue("Listing S3 buckets...")

	output, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Failed to list S3 buckets: %v", err)
	}

	for _, bucket := range output.Buckets {
		color.Blue("Bucket Name: %s", *bucket.Name)
	}
}

// CustomCredentialsProvider implements aws.CredentialsProvider
type CustomCredentialsProvider struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (p CustomCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
	}, nil
}
