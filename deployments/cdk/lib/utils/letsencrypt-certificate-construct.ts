import * as cdk from 'aws-cdk-lib';
import {Duration, triggers} from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as logs from 'aws-cdk-lib/aws-logs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import {Construct, IDependable} from 'constructs';
import {AwsCustomResource, AwsCustomResourcePolicy, PhysicalResourceId} from 'aws-cdk-lib/custom-resources';
import {ICertificate} from "aws-cdk-lib/aws-certificatemanager";
import * as fs from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';

export interface LetsEncryptCertificateConstructProps {
    environmentName: string;
    domainName: string;
    certificateEmail: string;
    ssmParameterSuffix?: string;
}

const hashes = {
    lambda: '',
}

export async function computeHash() {
    const rootDir = path.resolve(__dirname, '../../../../');
    const hash = crypto.createHash('sha256');
    const filesToHash: string[] = [];

    filesToHash.push(path.join(rootDir, 'go.mod'));
    filesToHash.push(path.join(rootDir, 'api/lambdas/go.mod'));

    const collectGoFiles = (dir: string) => {
        if (!fs.existsSync(dir)) {
            return;
        }
        const entries = fs.readdirSync(dir, {withFileTypes: true});
        for (const entry of entries) {
            const fullPath = path.join(dir, entry.name);
            if (entry.isDirectory()) {
                collectGoFiles(fullPath);
            } else if (entry.isFile() && entry.name.endsWith('.go')) {
                filesToHash.push(fullPath);
            }
        }
    };

    collectGoFiles(path.join(rootDir, 'api/lambdas/sys-letsencrypt'));
    collectGoFiles(path.join(rootDir, 'pkg/dns'));
    collectGoFiles(path.join(rootDir, 'pkg/dnsadapters'));

    filesToHash.sort();

    for (const file of filesToHash) {
        if (fs.existsSync(file)) {
            const content = fs.readFileSync(file);
            hash.update(file);
            hash.update(new Uint8Array(content));
        }
    }

    hashes.lambda = hash.digest('hex');
    return hashes.lambda;
}

export class LetsEncryptCertificateConstruct extends Construct {
    public readonly certificate: ICertificate;
    private readonly ssmParameterSuffix: string;

    constructor(scope: Construct, id: string, props: LetsEncryptCertificateConstructProps) {
        super(scope, id);
        
        this.ssmParameterSuffix = props.ssmParameterSuffix || 'domainCertificationArn';

        const letsEncryptLambdaTrigger: triggers.Trigger = this.installCertificateRenewalMechanism(props)
        const certificateArn = this.readCertificateARN(letsEncryptLambdaTrigger, props.environmentName);
        this.certificate = cdk.aws_certificatemanager.Certificate.fromCertificateArn(
            this,
            'Certificate',
            certificateArn
        )
    }

    private installCertificateRenewalMechanism({environmentName, certificateEmail, domainName}: LetsEncryptCertificateConstructProps) {
        const lambdaRole = new iam.Role(this, 'Role', {
            roleName: `dphoto-${environmentName}-letsencrypt`,
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

        const logGroup = new logs.LogGroup(this, 'LogGroup', {
            logGroupName: `/dphoto/${environmentName}/lambda/system-letsencrypt`,
            retention: logs.RetentionDays.ONE_WEEK,
            removalPolicy: cdk.RemovalPolicy.DESTROY
        });

        const letsEncryptLambda = new lambda.Function(this, 'RenewalLambda', {
            functionName: `dphoto-${environmentName}-system-letsencrypt`,
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
            logGroup: logGroup
        });

        new events.Rule(this, 'RenewalSchedule', {
            ruleName: `dphoto-${environmentName}-letsencrypt-renewal`,
            schedule: events.Schedule.cron({
                minute: '42',
                hour: '9',
                weekDay: '2'
            })
        }).addTarget(new targets.LambdaFunction(letsEncryptLambda));

        return new triggers.Trigger(this, 'RenewalTrigger', {
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

    private getSsmKeyCertificateArn(environmentName: string) {
        return `/dphoto/${environmentName}/acm/${this.ssmParameterSuffix}`;
    }
}
