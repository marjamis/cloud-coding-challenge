import sys
import pymysql as mysql

# AWS Code: Begin
#  Copyright 2016 Amazon Web Services, Inc. or its affiliates. All Rights Reserved.
#  A copy of the License is located at http://aws.amazon.com/agreement/ .
from botocore.vendored import requests
import json

SUCCESS = "SUCCESS"
FAILED = "FAILED"


def send(event, context, responseStatus, responseData, physicalResourceId):
    responseUrl = event['ResponseURL']

    print(responseUrl)

    responseBody = {}
    responseBody['Status'] = responseStatus
    responseBody['Reason'] = 'See the details in CloudWatch Log Stream: ' + \
        context.log_stream_name
    responseBody['PhysicalResourceId'] = physicalResourceId or context.log_stream_name
    responseBody['StackId'] = event['StackId']
    responseBody['RequestId'] = event['RequestId']
    responseBody['LogicalResourceId'] = event['LogicalResourceId']
    responseBody['Data'] = responseData

    json_responseBody = json.dumps(responseBody)

    print("Response body:\n" + json_responseBody)

    headers = {
        'content-type': '',
        'content-length': str(len(json_responseBody))
    }

    try:
        response = requests.put(responseUrl,
                                data=json_responseBody,
                                headers=headers)
        print("Status code: " + response.reason)
    except Exception as e:
        print("send(..) failed executing requests.put(..): " + str(e))
# AWS Code: End

# Handler which does the bulk of the work


def handler(event, context):
    responseData = {}
    print(event)  # Insecure but for testing it's really handy

    # CREATE
    if event['RequestType'] == 'Create':
        try:
            db = mysql.connect(event['ResourceProperties']['DBEndpoint'], event['ResourceProperties']
                               ['User'], event['ResourceProperties']['Password'], event['ResourceProperties']['DB'])
            cursor = db.cursor()
            cursor.execute(
                "create table metis (id int NOT NULL AUTO_INCREMENT PRIMARY KEY, instanceId VARCHAR(500), requesterIp VARCHAR(50), addedDate TIMESTAMP)")
            db.close()
            send(event, context, SUCCESS, responseData, 'LambdaFunction')
        except Exception:
            print("Unexpected error:", sys.exc_info())
            send(event, context, FAILED, responseData, 'LambdaFunction')

    # UPDATE - just sends success for CFN as nothing else is required
    if event['RequestType'] == 'Update':
        send(event, context, SUCCESS, responseData, 'LambdaFunction')

    # DELETE - just sends success for CFN as nothing else is required
    if event['RequestType'] == 'Delete':
        send(event, context, SUCCESS, responseData, 'LambdaFunction')
