// Package: services
// File: minio.go
// Code snippet:
//
//	type MinioConfig struct {
//		Endpoint        string
//		AccessKeyID     string
//		SecretAccessKey string
//		Region          string
//		DisableSSL      bool
//	}
package services

import (
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	DisableSSL      bool
}

type MinioService struct {
	config   MinioConfig
	s3Config *aws.Config
}

func NewMinioService() (*MinioService, error) {
	cfg, err := prepareConfig()
	if err != nil {
		return nil, err
	}

	minioService := &MinioService{
		config:   cfg,
		s3Config: addAwsConfig(cfg),
	}
	return minioService, nil
}

func (m *MinioService) GetTerraformState() ([]byte, error) {
	newSession, err := session.NewSession(m.s3Config)
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(newSession)

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("vm-files"),
		Key:    aws.String("terraform.tfstate"),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

// initializes viper to read from configuration files
func prepareConfig() (MinioConfig, error) {
	minioConfig := MinioConfig{}
	viper.SetConfigName("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Access the variables using viper directly
	accessKeyId := viper.GetString("ACCESS_KEY_ID")
	secretAccessKey := viper.GetString("SECRET_ACCESS_KEY")
	bucketServerEndpoint := viper.GetString("BUCKET_SERVER_ENDPOINT")
	region := viper.GetString("REGION")
	disableSSL := viper.GetBool("DISABLE_SSL")
	forcePathStyleAwsUrl := viper.GetBool("FORCE_PATHSTYLE_AWS_URL")

	log.Printf("Config:ACCESS_KEY_ID: %s\nSECRET_ACCESS_KEY: %s\nBUCKET_SERVER_ENDPOINT: %s\nREGION: %s\nDISABLE_SSL: %v\nFORCE_PATHSTYLE_URL: %v\n",
		accessKeyId, secretAccessKey,
		bucketServerEndpoint, region,
		disableSSL, forcePathStyleAwsUrl)

	if accessKeyId == "" || secretAccessKey == "" || bucketServerEndpoint == "" || region == "" {
		return minioConfig, fmt.Errorf("missing config voices in config file")
	}

	if err := viper.ReadInConfig(); err != nil {
		return minioConfig, err
	}
	minioConfig.AccessKeyID = accessKeyId
	minioConfig.SecretAccessKey = secretAccessKey
	minioConfig.Endpoint = bucketServerEndpoint
	minioConfig.Region = region
	minioConfig.DisableSSL = disableSSL

	return minioConfig, nil
}

func addAwsConfig(c MinioConfig) *aws.Config {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
		Endpoint:         aws.String(c.Endpoint),
		Region:           aws.String(c.Region),   // Dummy value, not used by Minio
		DisableSSL:       aws.Bool(c.DisableSSL), // Set to true if your Minio server is not using SSL
		S3ForcePathStyle: aws.Bool(true),         // Important for Minio, enforces path-style url instead of virtual style
	}
	return s3Config
}
