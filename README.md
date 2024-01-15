# vm-code-service
Simple service that is used to read preseed files and
other useful configuration files you need inside your template or terraform created vms.

The idea of a service is to prevent embedding the preseed file inside the .iso file when working with scenarios where you have no dhcp, you have configs you need in the template but cannot go public.

The code expects a MinIO configuration to read from a bucket and serve the files from that bucket.

The bucket could be:
- minio
- s3

so any compatible s3 source.