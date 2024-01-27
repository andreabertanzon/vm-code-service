# vm-code-service
This is a simple service in Golang that can be used as a support for handling Terraform state from http via MinIO. The actual S3 provider from Terraform should actually allow you to use MinIO, but I had no luck with it.

## Functionalities
### Bucket Folder Dowload
Allows you to dowload a the content of an entire MinIO folder in zip format. This is helpful if you need packer templates, or Terraform provisioners to prepare content on the machine (Kubernetes scripts, various config files)
You basically define your what your machine template is (for example K8s-node) and make sure there is a folder on your vm-files bucket in minio called K8s-node containig all the files you want pre provisioned in the machine right from the get go.

### Terraform state
NB: The feature is not yet fully complete, missing lock on file

This Functionality allows you to store state files on a bucket (on ANY s3 compatible services). The state file can be used by terraform to check the state of your infrastructure.