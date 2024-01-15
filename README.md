# vm-code-service
Simple service that is used to read preseed files and other useful configuration files you need inside your template or terraform created vms.

The idea of a service is to **prevent** embedding the preseed file inside the .iso file (or any other config file or script) when working with scenarios where you have no dhcp, you have configs you need in the template but cannot go public.

The bucket could be:
- minio
- s3
- any other compatible s3 source.

A .env file is used to parse the necessary configuration for your s3 compatible service (like   minio) the parsing is done via Viper library