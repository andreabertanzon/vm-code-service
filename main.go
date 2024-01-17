package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	s3Config := prepareAwsConfig()

	queryParams := r.URL.Query()
	fileName := queryParams.Get("file")
	if fileName == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "You must specify a file to query for, ?file=pippo.txt")
		return
	}

	download := queryParams.Get("download")
	if download == "true" {
		fmt.Println("filename:", fileName)
		fmt.Println("filename:", download)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		w.Header().Add("Content-Type", "text/plain")
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		http.Error(w, "Error establishing S3 connection", http.StatusInternalServerError)
	}

	s3Client := s3.New(newSession)

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("vm-files"),
		Key:    aws.String(fileName),
	})
	if err != nil {
		http.Error(w, "Error connecting to S3 bucket", http.StatusInternalServerError)
		return
	}
	defer result.Body.Close()

	// Serve the file
	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, "Error sending file.", http.StatusInternalServerError)
	}

}

// initializes viper to read from configuration files
func prepareConfig() error {
	viper.SetConfigName("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

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
