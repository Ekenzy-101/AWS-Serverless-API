import { Auth, StackContext } from "@serverless-stack/resources";
import {
  AccountRecovery,
  Mfa,
  StandardAttributesMask,
  ClientAttributes,
  UserPoolClientIdentityProvider,
} from "aws-cdk-lib/aws-cognito";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import {
  AwsCustomResource,
  PhysicalResourceId,
  AwsCustomResourcePolicy,
} from "aws-cdk-lib/custom-resources";
import { removalPolicy } from "./utils";

export function AuthStack({ stack, app }: StackContext) {
  const standardAttributes: StandardAttributesMask = {
    email: true,
    emailVerified: true,
    fullname: true,
  };
  const clientReadAttributes = new ClientAttributes().withStandardAttributes(
    standardAttributes
  );
  const clientWriteAttributes = new ClientAttributes().withStandardAttributes({
    ...standardAttributes,
    emailVerified: false,
  });

  const auth = new Auth(stack, "main", {
    identityPoolFederation: false,
    login: ["email", "username"],
    cdk: {
      userPool: {
        accountRecovery: AccountRecovery.EMAIL_ONLY,
        mfa: Mfa.OPTIONAL,
        mfaSecondFactor: { otp: true, sms: false },
        removalPolicy: removalPolicy(app),
      },
      userPoolClient: {
        userPoolClientName: app.logicalPrefixedName("server"),
        authFlows: {
          custom: true,
          userSrp: true,
          userPassword: true,
        },
        supportedIdentityProviders: [UserPoolClientIdentityProvider.COGNITO],
        readAttributes: clientReadAttributes,
        writeAttributes: clientWriteAttributes,
        generateSecret: true,
      },
    },
    defaults: {
      function: {
        srcPath: "backend",
        runtime: "go1.x",
      },
    },
    triggers: {
      preSignUp: {
        handler: "triggers/preSignup/main.go",
        environment: {
          APP_ENV: app.stage,
          COGNITO_USER_POOL_ID: "",
        },
      },
    },
  });

  auth.attachPermissionsForTriggers([
    new PolicyStatement({
      actions: ["cognito-idp:*"],
      resources: ["*"],
    }),
  ]);

  const userPoolClientId = auth.userPoolClientId;
  const userPoolId = auth.userPoolId;
  const userPoolClientSecret = new AwsCustomResource(
    stack,
    "DescribeCognitoUserPoolClient",
    {
      resourceType: "Custom::DescribeCognitoUserPoolClient",
      onCreate: {
        region: app.region,
        service: "CognitoIdentityServiceProvider",
        action: "describeUserPoolClient",
        parameters: {
          UserPoolId: userPoolId,
          ClientId: userPoolClientId,
        },
        physicalResourceId: PhysicalResourceId.of(userPoolClientId),
      },
      // TODO: can we restrict this policy more?
      policy: AwsCustomResourcePolicy.fromSdkCalls({
        resources: AwsCustomResourcePolicy.ANY_RESOURCE,
      }),
    }
  ).getResponseField("UserPoolClient.ClientSecret");

  return {
    userPoolId,
    userPoolClientId,
    userPoolClientSecret,
  };
}
