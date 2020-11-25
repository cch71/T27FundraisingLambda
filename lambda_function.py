import json
import boto3
from decimal import Decimal
from boto3.dynamodb.conditions import Key

def json_default_encoder(obj):
     if isinstance(obj, Decimal):
        if obj % 1 > 0:
           return float(obj)
        else:
            return int(obj)
     raise TypeError("Object of type '%s' is not JSON serializable" % type(obj).__name__)

def lambda_handler(event, context):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('T27FundraiserOrders')


    req = event['body']
    if isinstance(req, str):
        req = json.loads(req)
    
    retvals=[]
    if 'orderOwner' in req:
        # Expression attribute names can only reference items in the projection expression.
        #ProjectionExpression=", ".join(req['fields']),
        query_args = {}
        if 'orderId' in req:
            query_args['KeyConditionExpression'] = Key('orderOwner').eq(req['orderOwner']) & Key('orderId').eq(req['orderId'])
        else:
            query_args['KeyConditionExpression'] = Key('orderOwner').eq(req['orderOwner'])

        if 'fields' in req:
            query_args['ProjectionExpression'] = ", ".join(req['fields'])

        response = table.query(**query_args)
                
        retvals=response['Items']

    retvals = json.dumps(retvals, default=json_default_encoder)
    print(f"Query Resp: {retvals}")

    return {
        'statusCode': 200,
        'body': retvals,
        'headers': {
            'Access-Control-Allow-Origin': '*',
        },
    }
    
