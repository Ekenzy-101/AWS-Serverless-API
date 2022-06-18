import { StorageStack } from "./StorageStack";
import type { App } from "@serverless-stack/resources";
import { AuthStack } from "./AuthStack";
import { ApiStack } from "./ApiStack";
import { removalPolicy } from "./utils";

export default function (app: App) {
  app.setDefaultRemovalPolicy(removalPolicy(app));
  app.stack(StorageStack).stack(AuthStack).stack(ApiStack);
}
