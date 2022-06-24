import { StackContext, Table, Bucket } from "@serverless-stack/resources";
import { RemovalPolicy } from "aws-cdk-lib";
import { BucketAccessControl } from "aws-cdk-lib/aws-s3";
import { removalPolicy } from "./utils";

export function StorageStack({ stack, app }: StackContext) {
  const table = new Table(stack, "main", {
    fields: {
      pk: "string",
      sk: "string",
    },
    primaryIndex: {
      partitionKey: "pk",
      sortKey: "sk",
    },
  });

  const policy = removalPolicy(app);
  const s3 = new Bucket(stack, `${app.stage}-bucket`, {
    cors: [
      {
        allowedMethods: ["GET", "PUT"],
        allowedOrigins: ["*"],
        allowedHeaders: ["*"],
        maxAge: "12 hours",
      },
    ],
    name: app.logicalPrefixedName("main"),
    cdk: {
      bucket: {
        accessControl: BucketAccessControl.PUBLIC_READ,
        removalPolicy: policy,
        autoDeleteObjects: policy == RemovalPolicy.DESTROY ? true : false,
        publicReadAccess: true,
      },
    },
  });

  return { table, s3 };
}
