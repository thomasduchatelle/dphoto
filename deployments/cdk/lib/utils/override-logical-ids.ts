import {CfnResource} from "aws-cdk-lib";
import {Construct} from "constructs/lib/construct";

export function pinLogicalId(construct: Construct, newLogicalId: string) {
    const resource = construct.node.defaultChild as CfnResource;
    resource.overrideLogicalId(newLogicalId)
}