This is a simple server that scrapes AWS RDS Instances

### Build
```
make
```

### Run
```
./aws_rds_exporter --aws.region=eu-west-1
```

## Exposed metrics
The `aws_rds_exporter` exports various RDS Metrics


## Docker
You can deploy this exporter using the [jimdo/aws-health-exporter](https://hub.docker.com/r/jimdo/aws-health-exporter/) Docker Image.

Example
```
docker pull alec2435/aws_rds_exporter
docker run -p 9383:9383 alec2436/aws_rds_exporter
```

### Credentials
The `aws-rds_exporter` requires AWS credentials to access the AWS RDS API. For example you can pass them via env vars using `-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}` options.

