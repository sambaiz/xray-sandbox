import * as cdk from '@aws-cdk/core';
import { CfnSamplingRule, CfnGroup } from '@aws-cdk/aws-xray';
export class XraySandBoxStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new CfnSamplingRule(this, 'TestAppSamplingRule', {
      samplingRule: {
        ruleName: "test-app",
        resourceArn: "*",
        priority: 10,
        fixedRate: 0,
        reservoirSize: 10,
        serviceName: "*",
        serviceType: "*",
        host: "localhost:8080",
        httpMethod: "*",
        urlPath: "*",
        version: 1,
      }
    })

    new CfnGroup(this, 'TestAppGroup', {
      groupName: "test-app-group",
      filterExpression: 'http.url CONTAINS "localhost:8080"',
      insightsConfiguration: {
        insightsEnabled: true,
        notificationsEnabled: false
      }
    })

    // The code that defines your stack goes here
  }
}
