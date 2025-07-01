export interface EnvironmentConfig {
    domainName?: string;
    enableMonitoring: boolean;
    lambdaMemory: number;
    dynamoDbBillingMode: 'PAY_PER_REQUEST' | 'PROVISIONED';
    tableName: string;
    bucketName: string;
}

export const environments: Record<string, EnvironmentConfig> = {
    dev: {
        domainName: undefined,
        enableMonitoring: false,
        lambdaMemory: 512,
        dynamoDbBillingMode: 'PAY_PER_REQUEST',
        tableName: 'dphoto-dev',
        bucketName: 'dphoto-dev'
    },
    live: {
        domainName: 'photos.duchatelle.net', // Update with actual domain
        enableMonitoring: true,
        lambdaMemory: 1024,
        dynamoDbBillingMode: 'PAY_PER_REQUEST',
        tableName: 'dphoto-live',
        bucketName: 'dphoto-live'
    }
};