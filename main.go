package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/viper"
)

var (
	accessKeyId          string
	secretAccessKey      string
	bucketServerEndpoint string
	region               string
	disableSSL           bool
	forcePathStyleAwsUrl bool
)

func main() {
	err := prepareConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = prepareAwsConfig()

	//TODO: instead of hardcoding the file, just find another option
	http.HandleFunc("/", handleVMFileRequest)

	/**TODO: instead of hardcoding minio stuff, add code that can interact with S3,
	    * the idea is to start the service and create the bucket if it does not exist.
		* use s3 if possible.
		* this means also setting up a way to config the code so that you can configure the
		* service for any user. Any bucket.
	*/
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareAwsConfig() *aws.Config {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		Endpoint:         aws.String(bucketServerEndpoint),
		Region:           aws.String(region),             // Dummy value, not used by Minio
		DisableSSL:       aws.Bool(disableSSL),           // Set to true if your Minio server is not using SSL
		S3ForcePathStyle: aws.Bool(forcePathStyleAwsUrl), // Important for Minio, enforces path-style url instead of virtual style
	}
	return s3Config
}

func handleVMFileRequest(w http.ResponseWriter, r *http.Request) {
	// handle the files by query parameter like file=
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, World")
}

// initializes viper to read from configuration files
func prepareConfig() error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Access the variables using viper directly
	accessKeyId = viper.GetString("ACCESS_KEY_ID")
	secretAccessKey = viper.GetString("SECRET_ACCESS_KEY")
	bucketServerEndpoint = viper.GetString("BUCKET_SERVER_ENDPOINT")
	region = viper.GetString("REGION")
	disableSSL = viper.GetBool("DISABLE_SSL")
	forcePathStyleAwsUrl = viper.GetBool("FORCE_PATHSTYLE_AWS_URL")

	log.Printf("Config:ACCESS_KEY_ID: %s\nSECRET_ACCESS_KEY: %s\nBUCKET_SERVER_ENDPOINT: %s\nREGION: %s\nDISABLE_SSL: %v\nFORCE_PATHSTYLE_URL: %v\n",
		accessKeyId, secretAccessKey,
		bucketServerEndpoint, region,
		disableSSL, forcePathStyleAwsUrl)

	if accessKeyId == "" || secretAccessKey == "" || bucketServerEndpoint == "" || region == "" {
		return fmt.Errorf("missing config voices in config file")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
