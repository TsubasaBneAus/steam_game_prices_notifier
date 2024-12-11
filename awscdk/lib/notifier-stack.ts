import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as dotenv from "dotenv";

dotenv.config();

// A stack for the Steam Game Prices Notifier
export class NotifierStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create a log group
    const logGroup = new cdk.aws_logs.LogGroup(this, "LogGroup", {
      logGroupName: "steam-game-prices-notifier-log-group",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      retention: cdk.aws_logs.RetentionDays.ONE_WEEK,
    });

    // Create a Lambda function
    const lambda = new cdk.aws_lambda.Function(this, "Lambda", {
      functionName: "steam-game-prices-notifier-lambda",
      code: cdk.aws_lambda.Code.fromAsset("../function.zip"),
      architecture: cdk.aws_lambda.Architecture.ARM_64,
      runtime: cdk.aws_lambda.Runtime.PROVIDED_AL2023,
      handler: "bootstrap",
      environment: {
        NOTION_API_KEY: process.env.NOTION_API_KEY ?? "",
        NOTION_DATABASE_ID: process.env.NOTION_DATABASE_ID ?? "",
        DISCORD_WEBHOOK_ID: process.env.DISCORD_WEBHOOK_ID ?? "",
        DISCORD_WEBHOOK_TOKEN: process.env.DISCORD_WEBHOOK_TOKEN ?? "",
        STEAM_USER_ID: process.env.STEAM_USER_ID ?? "",
      },
      timeout: cdk.Duration.minutes(2),
      logGroup: logGroup,
      loggingFormat: cdk.aws_lambda.LoggingFormat.JSON,
    });

    // Create a EventBridge rule (UTC)
    const rule = new cdk.aws_events.Rule(this, "Rule", {
      ruleName: "steam-game-prices-notifier-rule",
      schedule: cdk.aws_events.Schedule.cron({
        minute: "0",
        hour: "9",
        day: "*",
        month: "*",
        year: "*",
      }),
    });
    rule.addTarget(new cdk.aws_events_targets.LambdaFunction(lambda));

    // Create a OIDC provider
    const provider = new cdk.aws_iam.OpenIdConnectProvider(
      this,
      "OIDCProvider",
      {
        url: "https://token.actions.githubusercontent.com",
        clientIds: ["sts.amazonaws.com"],
      }
    );

    // Create a role for the OIDC provider
    const role = new cdk.aws_iam.Role(this, "Role", {
      roleName: "steam-game-prices-notifier-github-actions-role",
      assumedBy: new cdk.aws_iam.WebIdentityPrincipal(
        provider.openIdConnectProviderArn,
        {
          StringEquals: {
            "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
          },
          StringLike: {
            "token.actions.githubusercontent.com:sub":
              "repo:TsubasaBneAus/steam_game_prices_notifier:*",
          },
        }
      ),
    });
    role.addToPolicy(
      new cdk.aws_iam.PolicyStatement({
        actions: ["lambda:UpdateFunctionCode"],
        resources: [lambda.functionArn],
      })
    );
  }
}
