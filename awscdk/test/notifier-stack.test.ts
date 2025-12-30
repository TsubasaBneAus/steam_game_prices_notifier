import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { NotifierStack } from "../lib/notifier-stack";

describe("Assertion tests", () => {
  let template: Template;
  beforeAll(() => {
    const app = new cdk.App();
    const stack = new NotifierStack(app, "NotifierStack");
    template = Template.fromStack(stack);
  });

  test("1 Log Group exists", () => {
    template.resourcePropertiesCountIs(
      "AWS::Logs::LogGroup",
      {
        LogGroupName: "steam-game-prices-notifier-log-group",
      },
      1
    );
  });

  test("Retention period for the Log Group is 1 week", () => {
    template.hasResourceProperties("AWS::Logs::LogGroup", {
      LogGroupName: "steam-game-prices-notifier-log-group",
      RetentionInDays: 7,
    });
  });

  test("1 Lambda function exists", () => {
    template.resourcePropertiesCountIs(
      "AWS::Lambda::Function",
      {
        FunctionName: "steam-game-prices-notifier-lambda",
      },
      1
    );
  });

  test("The CPU architecture of the Lambda function is ARM 64", () => {
    template.hasResourceProperties("AWS::Lambda::Function", {
      FunctionName: "steam-game-prices-notifier-lambda",
      Architectures: ["arm64"],
    });
  });

  test("The runtime of the Lambda function is 'provided.al2023'", () => {
    template.hasResourceProperties("AWS::Lambda::Function", {
      FunctionName: "steam-game-prices-notifier-lambda",
      Runtime: "provided.al2023",
    });
  });

  test("The handler of the Lambda function is 'bootstrap'", () => {
    template.hasResourceProperties("AWS::Lambda::Function", {
      FunctionName: "steam-game-prices-notifier-lambda",
      Handler: "bootstrap",
    });
  });

  test("The timeout of the Lambda function is 2 minutes", () => {
    template.hasResourceProperties("AWS::Lambda::Function", {
      FunctionName: "steam-game-prices-notifier-lambda",
      Timeout: 120,
    });
  });

  test("The logging format of the Lambda function is 'JSON'", () => {
    template.hasResourceProperties("AWS::Lambda::Function", {
      FunctionName: "steam-game-prices-notifier-lambda",
      LoggingConfig: {
        LogFormat: "JSON",
      },
    });
  });

  test("1 EventBridge rule exists", () => {
    template.resourcePropertiesCountIs(
      "AWS::Events::Rule",
      {
        Name: "steam-game-prices-notifier-rule",
      },
      1
    );
  });

  test("The schedule of the EventBridge rule is 'cron(0 9 * * ? *)'", () => {
    template.hasResourceProperties("AWS::Events::Rule", {
      Name: "steam-game-prices-notifier-rule",
      ScheduleExpression: "cron(0 9 * * ? *)",
    });
  });

  test("1 OIDC provider exists", () => {
    template.resourceCountIs("AWS::IAM::OIDCProvider", 1);
  });

  test("The URL of the OIDC provider is 'https://token.actions.githubusercontent.com'", () => {
    template.hasResourceProperties("AWS::IAM::OIDCProvider", {
      Url: "https://token.actions.githubusercontent.com",
    });
  });

  test("The client IDs of the OIDC provider are 'sts.amazonaws.com'", () => {
    template.hasResourceProperties("AWS::IAM::OIDCProvider", {
      ClientIdList: ["sts.amazonaws.com"],
    });
  });

  test("1 IAM role exists", () => {
    template.resourcePropertiesCountIs(
      "AWS::IAM::Role",
      {
        RoleName: "steam-game-prices-notifier-github-actions-role",
      },
      1
    );
  });

  test("The conditions of the IAM role are correct", () => {
    template.hasResourceProperties("AWS::IAM::Role", {
      RoleName: "steam-game-prices-notifier-github-actions-role",
      AssumeRolePolicyDocument: {
        Statement: [
          {
            Condition: {
              StringEquals: {
                "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
              },
              StringLike: {
                "token.actions.githubusercontent.com:sub": "repo:TsubasaBneAus/steam_game_prices_notifier:*",
              },
            },
          },
        ],
      },
    });
  });
});
