import type { App } from "@serverless-stack/resources";
import { RemovalPolicy } from "aws-cdk-lib";

export function removalPolicy(app: App) {
  return app.local ? RemovalPolicy.DESTROY : RemovalPolicy.RETAIN;
}
