import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';

export interface Workload {
    role: iam.IGrantable;
    function?: lambda.Function;
}
