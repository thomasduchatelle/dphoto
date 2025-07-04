import * as cdk from 'aws-cdk-lib';
import {Duration, triggers} from 'aws-cdk-lib';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53_targets from 'aws-cdk-lib/aws-route53-targets';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import {Construct, IDependable} from 'constructs';
import {AwsCustomResource, AwsCustomResourcePolicy, PhysicalResourceId} from 'aws-cdk-lib/custom-resources';

export interface ApiGatewayConstructProps {
    environmentName: string;
    domainName: string;
    rootDomain: string;
    certificateEmail: string;
}

export class ApiGatewayConstruct extends Construct {
    public readonly httpApi: apigatewayv2.HttpApi;
    public readonly domainName: apigatewayv2.DomainName;

    constructor(scope: Construct, id: string, props: ApiGatewayConstructProps) {
        super(scope, id);

        const letsEncryptLambdaTrigger: triggers.Trigger = this.installCertificateRenewalMechanism(props)

        const certificateArn = this.readCertificateARN(letsEncryptLambdaTrigger, props.environmentName);

        const {httpApi, domainName} = this.createAPIGateway(props, certificateArn)
        this.httpApi = httpApi;
        this.domainName = domainName;
    }

    private installCertificateRenewalMechanism({environmentName, certificateEmail, domainName}: ApiGatewayConstructProps) {
        const lambdaRole = new iam.Role(this, 'LetsEncryptRole', {
            roleName: `dphoto-${environmentName}-letsencrypt-role`,
            path: `/dphoto/${environmentName}/`,
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
                                `arn:aws:ssm:${cdk.Stack.of(this).region}:${cdk.Stack.of(this).account}:parameter/dphoto/${environmentName}/*`
                            ]
                        })
                    ]
                })
            }
        });

        const letsEncryptLambda = new lambda.Function(this, 'LetsEncryptRenewal', {
            functionName: `dphoto-${environmentName}-sys-letsencrypt`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset('../../bin/sys-letsencrypt.zip'),
            role: lambdaRole,
            timeout: cdk.Duration.minutes(15),
            memorySize: 128,
            environment: {
                DPHOTO_DOMAIN: domainName,
                DPHOTO_CERTIFICATE_EMAIL: certificateEmail,
                DPHOTO_ENVIRONMENT: environmentName,
                SSM_KEY_CERTIFICATE_ARN: this.getSsmKeyCertificateArn(environmentName),
            },
            logRetention: logs.RetentionDays.ONE_WEEK
        });

        new events.Rule(this, 'LetsEncryptRenewalSchedule', {
            ruleName: `dphoto-${environmentName}-letsencrypt-schedule`,
            schedule: events.Schedule.cron({
                minute: '42',
                hour: '9',
                weekDay: '2'
            })
        }).addTarget(new targets.LambdaFunction(letsEncryptLambda));

        return new triggers.Trigger(this, 'LetsEncryptRenewalTrigger', {
            handler: letsEncryptLambda,
            timeout: Duration.minutes(5),
            invocationType: triggers.InvocationType.REQUEST_RESPONSE,
        });
    }

    private readCertificateARN(letsEncryptLambdaTrigger: IDependable, environmentName: string): string {
        const certificateLookup = new AwsCustomResource(this, 'CertificateLookup', {
            onCreate: {
                service: 'SSM',
                action: 'getParameter',
                parameters: {
                    Name: this.getSsmKeyCertificateArn(environmentName)
                },
                physicalResourceId: PhysicalResourceId.of('cert-arn-lookup')
            },
            policy: AwsCustomResourcePolicy.fromSdkCalls({
                resources: [`arn:aws:ssm:${cdk.Stack.of(this).region}:${cdk.Stack.of(this).account}:parameter/dphoto/${environmentName}/*`]
            })
        });

        certificateLookup.node.addDependency(letsEncryptLambdaTrigger);
        return certificateLookup.getResponseField('Parameter.Value');
    }

    private createAPIGateway(props: ApiGatewayConstructProps, certificateArn: string) {
        const httpApi = new apigatewayv2.HttpApi(this, 'HttpApi', {
            apiName: `dphoto-${props.environmentName}-api`,
            description: `DPhoto API for ${props.environmentName} environment`,
            corsPreflight: {
                allowOrigins: ['*'],
                allowMethods: [apigatewayv2.CorsHttpMethod.ANY],
                allowHeaders: ['*']
            }
        });
        const domainName = new apigatewayv2.DomainName(this, 'DomainName', {
            domainName: props.domainName,
            certificate: cdk.aws_certificatemanager.Certificate.fromCertificateArn(
                this,
                'Certificate',
                certificateArn
            )
        });
        const hostedZone = route53.HostedZone.fromLookup(this, 'HostedZone', {
            domainName: props.rootDomain
        });

        new route53.ARecord(this, 'DnsRecord', {
            zone: hostedZone,
            recordName: props.domainName,
            target: route53.RecordTarget.fromAlias(
                new route53_targets.ApiGatewayv2DomainProperties(
                    domainName.regionalDomainName,
                    domainName.regionalHostedZoneId
                )
            )
        });

        return {httpApi, domainName};
    }

    private getSsmKeyCertificateArn(environmentName: string) {
        return `/dphoto/${environmentName}/acm/domainCertificationArn`;
    }
}