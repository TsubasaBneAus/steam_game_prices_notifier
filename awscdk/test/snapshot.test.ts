import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { NotifierStack } from "../lib/notifier-stack";

test("Snapshot test", () => {
  const app = new cdk.App();
  const stack = new NotifierStack(app, "NotifierStack");
  const template = Template.fromStack(stack);
  expect(template.toJSON()).toMatchSnapshot();
});
