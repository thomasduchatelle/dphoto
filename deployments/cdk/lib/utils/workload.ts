import * as iam from 'aws-cdk-lib/aws-iam';

export interface Workload {
    role: iam.IGrantable;
    function?: { addEnvironment(key: string, value: string): void };
}
