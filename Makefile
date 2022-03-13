export AWS_DEFAULT_REGION=us-west-2
export AWS_DEFAULT_PROFILE=metis
S3_BUCKET=testing-metis
ENVIRONMENT=test
STACK_NAME=metis01

$(eval TMP_DIR=$(shell mktemp -d))

default:

deploy: upload validate stack application

all: generateSelfSignedCert createSSHKeys deploy

application:
	cd ./application/src/metis && make deploy

generateSelfSignedCert:
	openssl genrsa -des3 -passout pass:x -out $(TMP_DIR)/server.pass.key 2048
	openssl rsa -passin pass:x -in $(TMP_DIR)/server.pass.key -out $(TMP_DIR)/priv.key
	openssl req -new -key $(TMP_DIR)/priv.key -out $(TMP_DIR)/server.csr -subj "/C=AU/ST=NSW/L=Sydney/O=Metis/OU=Website/CN=metis.qc.to"
	openssl x509 -req -days 365 -in $(TMP_DIR)/server.csr -signkey $(TMP_DIR)/priv.key -out $(TMP_DIR)/public.crt
	aws iam upload-server-certificate --server-certificate-name metis.qc.to --certificate-body file://$(TMP_DIR)/public.crt --private-key file://$(TMP_DIR)/priv.key
	echo "Update Certificate ARN in the parameters/Environment.json file, press enter to continue."
	read

createSSHKeys:
	ssh-keygen -t rsa -b 4096 -f ~/.ssh/metis.qc.to-us-west-2-$(ENVIRONMENT) && aws ec2 import-key-pair --key-name metis.qc.to-us-west-2-$(ENVIRONMENT) --public-key-material file://~/.ssh/metis.qc.to-us-west-2-$(ENVIRONMENT).pub

upload:
	#Upload Templates to S3
	cd cloudformation/lambda && zip -r $(TMP_DIR)/create_table_func.zip *
	aws s3 cp $(TMP_DIR)/create_table_func.zip s3://$(S3_BUCKET)/functions/
	aws s3 cp ./cloudformation/templates/VPC.cft.yml s3://$(S3_BUCKET)/templates/
	aws s3 cp ./cloudformation/templates/Environment.cft.yml s3://$(S3_BUCKET)/templates/

stack:
	aws cloudformation update-stack --stack-name $(STACK_NAME) --template-url http://s3-$(AWS_DEFAULT_REGION).amazonaws.com/$(S3_BUCKET)/templates/Environment.cft.yml --parameters $$(cat ./cloudformation/parameters/Environment.json) --capabilities CAPABILITY_IAM || aws cloudformation create-stack --stack-name $(STACK_NAME) --template-url http://s3-$(AWS_DEFAULT_REGION).amazonaws.com/$(S3_BUCKET)/templates/Environment.cft.yml --parameters $$(cat ./cloudformation/parameters/Environment.json) --capabilities CAPABILITY_IAM
	echo "Monitor Route53 HostedZone creation to update FreeDNS for NS, press enter to continue"
	read

validate:
	aws cloudformation validate-template --template-url http://s3-$(AWS_DEFAULT_REGION).amazonaws.com/$(S3_BUCKET)/templates/Environment.cft.yml
	aws cloudformation validate-template --template-url http://s3-$(AWS_DEFAULT_REGION).amazonaws.com/$(S3_BUCKET)/templates/VPC.cft.yml
