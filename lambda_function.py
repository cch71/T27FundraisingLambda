import json
import boto3

def lambda_handler(event, context):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('Orders')
    response = table.scan()
    items = response['Items']
    print(items)

    return {
        'statusCode': 200,
        'body': json.dumps({'items': items}),
        'headers': {
            'Access-Control-Allow-Origin': '*',
        },
    }
