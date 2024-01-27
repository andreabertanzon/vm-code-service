package services

type TerraformStateHandler interface {
	GetTerraformState() ([]byte, error)
	PutTerraformState([]byte) error
}

type S3Service interface {
	TerraformStateHandler
	DowloadBucketFolderToZip(bucketName string, folder string) ([]byte, error)
}
