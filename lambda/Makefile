STACK_NAME=kubesnapshot
BUCKET_NAME=kubesnapshot.$(KUBE_CLUSTER_NAME)

.PHONY: bucket deploy destroy
all: kubesnapshot

kubesnapshot: *.go ../*/*.go
	go get -v
	GOOS=linux GOARCH=amd64 go build -o $@

bucket:
	aws s3 ls s3://$(BUCKET_NAME) || aws s3 mb s3://$(BUCKET_NAME) --region $(AWS_REGION)

packaged.yaml: kubesnapshot template.yaml
	sam package --template-file template.yaml \
		--output-template-file packaged.yaml \
		--s3-bucket $(BUCKET_NAME)

deploy: packaged.yaml
	sam deploy --template-file packaged.yaml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--region $(AWS_REGION) \
		--parameter-overrides KubeClusterName=$(KUBE_CLUSTER_NAME)

destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME) --region $(AWS_REGION)
	aws s3 rb s3://$(BUCKET_NAME) --force --region $(AWS_REGION)
