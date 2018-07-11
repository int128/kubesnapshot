BUCKET_NAME=kubesnapshot
STACK_NAME=kubesnapshot

.PHONY: bucket deploy destroy
all: kubesnapshot

kubesnapshot: *.go ../*/*.go
	go get -v
	GOOS=linux GOARCH=amd64 go build -o $@

bucket:
	aws s3 mb s3://$(BUCKET_NAME) --region $(AWS_REGION)

deploy: kubesnapshot
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket $(BUCKET_NAME)
	sam deploy --template-file packaged.yaml --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --region $(AWS_REGION)

destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME) --region $(AWS_REGION)
	aws s3 rb s3://$(BUCKET_NAME) --force --region $(AWS_REGION)
