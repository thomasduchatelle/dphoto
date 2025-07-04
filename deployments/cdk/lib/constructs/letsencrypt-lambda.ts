import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as logs from 'aws-cdk-lib/aws-logs';
import {Construct} from 'constructs';

export interface LetsEncryptLambdaConstructProps {
    environmentName: string;
    domainName: string;
    certificateEmail: string;
}

export class LetsEncryptLambdaConstruct extends Construct {
    public readonly lambdaFunction: lambda.Function;

    constructor(scope: Construct, id: string, props: LetsEncryptLambdaConstructProps) {
        super(scope, id);

        // Create IAM role for Let's Encrypt lambda
        const lambdaRole = new iam.Role(this, 'LetsEncryptRole', {
            roleName: `dphoto-${props.environmentName}-letsencrypt-role`,
            path: `/dphoto/${props.environmentName}/`,
            assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
            managedPolicies: [
                iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole')
            ],
            inlinePolicies: {
                'lambda-certs': new iam.PolicyDocument({
                    statements: [
                        new iam.PolicyStatement({
                            effect: iam.Effect.ALLOW,
                            actions: [
                                'acm:AddTagsToCertificate',
                                'acm:DescribeCertificate',
                                'acm:ImportCertificate',
                                'acm:ListCertificates',
                                'acm:ListTagsForCertificate',
                                'acm:RemoveTagsFromCertificate',
                                'route53:ListHostedZonesByName',
                                'route53:ListResourceRecordSets',
                                'route53:ChangeResourceRecordSets',
                                'route53:GetChange'
                            ],
                            resources: ['*']
                        }),
                        new iam.PolicyStatement({
                            effect: iam.Effect.ALLOW,
                            actions: [
                                'ssm:GetParameter',
                                'ssm:PutParameter',
                                'ssm:AddTagsToResource',
                                'ssm:RemoveTagsFromResource'
                            ],
                            resources: [
                                `arn:aws:ssm:${cdk.Stack.of(this).region}:${cdk.Stack.of(this).account}:parameter/dphoto/${props.environmentName}/*`
                            ]
                        })
                    ]
                })
            }
        });

        // Create Lambda function
        this.lambdaFunction = new lambda.Function(this, 'Function', {
            functionName: `dphoto-${props.environmentName}-sys-letsencrypt`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset('../../bin/sys-letsencrypt.zip'),
            role: lambdaRole,
            timeout: cdk.Duration.minutes(15),
            memorySize: 128,
            environment: {
                DPHOTO_DOMAIN: props.domainName,
                DPHOTO_CERTIFICATE_EMAIL: props.certificateEmail,
                DPHOTO_ENVIRONMENT: props.environmentName
            },
            logRetention: logs.RetentionDays.ONE_WEEK
        });

        // Create EventBridge rule for scheduled execution (Tuesdays at 9:42 AM)
        const scheduleRule = new events.Rule(this, 'ScheduleRule', {
            ruleName: `dphoto-${props.environmentName}-letsencrypt-schedule`,
            schedule: events.Schedule.cron({
                minute: '42',
                hour: '9',
                weekDay: '2'
            })
        });

        scheduleRule.addTarget(new targets.LambdaFunction(this.lambdaFunction));

        // Trigger lambda once on creation using custom resource
        const triggerLambda = new lambda.Function(this, 'TriggerFunction', {
            functionName: `dphoto-${props.environmentName}-letsencrypt-trigger`,
            runtime: lambda.Runtime.NODEJS_18_X,
            handler: 'index.handler',
            code: lambda.Code.fromInline(`
                const AWS = require('aws-sdk');
                const lambda = new AWS.Lambda();
                
                exports.handler = async (event, context) => {
                    console.log('Event:', JSON.stringify(event, null, 2));
                    
                    if (event.RequestType === 'Create' || event.RequestType === 'Update') {
                        try {
                            const params = {
                                FunctionName: '${this.lambdaFunction.functionName}',
                                InvocationType: 'Event'
                            };
                            
                            await lambda.invoke(params).promise();
                            console.log('Successfully triggered Let\\'s Encrypt lambda');
                            
                            await sendResponse(event, context, 'SUCCESS', {});
                        } catch (error) {
                            console.error('Error triggering lambda:', error);
                            await sendResponse(event, context, 'FAILED', {});
                        }
                    } else {
                        await sendResponse(event, context, 'SUCCESS', {});
                    }
                };
                
                async function sendResponse(event, context, responseStatus, responseData) {
                    const responseBody = JSON.stringify({
                        Status: responseStatus,
                        Reason: 'See CloudWatch Log Stream: ' + context.logStreamName,
                        PhysicalResourceId: context.logStreamName,
                        StackId: event.StackId,
                        RequestId: event.RequestId,
                        LogicalResourceId: event.LogicalResourceId,
                        Data: responseData
                    });
                    
                    const https = require('https');
                    const url = require('url');
                    
                    const parsedUrl = url.parse(event.ResponseURL);
                    const options = {
                        hostname: parsedUrl.hostname,
                        port: 443,
                        path: parsedUrl.path,
                        method: 'PUT',
                        headers: {
                            'content-type': '',
                            'content-length': responseBody.length
                        }
                    };
                    
                    return new Promise((resolve, reject) => {
                        const request = https.request(options, (response) => {
                            resolve();
                        });
                        
                        request.on('error', (error) => {
                            reject(error);
                        });
                        
                        request.write(responseBody);
                        request.end();
                    });
                }
            `),
            timeout: cdk.Duration.minutes(5)
        });

        triggerLambda.addToRolePolicy(new iam.PolicyStatement({
            effect: iam.Effect.ALLOW,
            actions: ['lambda:InvokeFunction'],
            resources: [this.lambdaFunction.functionArn]
        }));

        new cdk.CustomResource(this, 'TriggerCustomResource', {
            serviceToken: triggerLambda.functionArn,
            properties: {
                // Change this value to re-trigger the lambda
                Dummy: `trigger-${Date.now()}`
            }
        });
    }
}
