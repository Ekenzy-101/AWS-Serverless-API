import { Api, StackContext, use } from "@serverless-stack/resources";
import { AuthStack } from "./AuthStack";
import { StorageStack } from "./StorageStack";

export function ApiStack({ stack, app }: StackContext) {
  const { table } = use(StorageStack);
  const { userPoolClientId, userPoolClientSecret, userPoolId } = use(AuthStack);

  const environment = {
    TABLE_NAME: table.tableName,
    APP_ENV: app.stage,
    COGNITO_USER_POOL_ID: userPoolId,
    COGNITO_APP_CLIENT_ID: userPoolClientId,
    COGNITO_APP_CLIENT_SECRET: userPoolClientSecret,
  };

  const api = new Api(stack, "main", {
    authorizers: {
      cognito: {
        type: "user_pool",
        userPool: {
          id: userPoolId,
          clientIds: [userPoolClientId],
        },
      },
    },
    cors: {
      allowHeaders: ["Authorization"],
      allowMethods: ["GET", "POST", "PUT", "OPTIONS"],
      allowOrigins: ["*"],
    },
    defaults: {
      function: {
        timeout: "30 seconds",
        runtime: "go1.x",
        srcPath: "backend",
        environment,
        permissions: [table],
      },
    },
    routes: {
      "GET /": { function: "functions/main.go" },
      "POST /accounts/register": { function: "functions/registerUser/main.go" },
      "POST /accounts/login": { function: "functions/loginUser/main.go" },
      "POST /accounts/login/verify": {
        function: "functions/verifyLogin/main.go",
      },
      "POST /accounts/logout": {
        authorizer: "cognito",
        function: "functions/logoutUser/main.go",
      },
      "PATCH /accounts/password/change": {
        authorizer: "cognito",
        function: "functions/changePassword/main.go",
      },
      "POST /accounts/password/forgot": {
        function: "functions/forgotPassword/main.go",
      },
      "PATCH /accounts/password/reset": {
        function: "functions/resetPassword/main.go",
      },
      "POST /accounts/token": { function: "functions/refreshTokens/main.go" },
      "PATCH /accounts/verify": { function: "functions/verifyUser/main.go" },
      "POST /accounts/mfa": {
        authorizer: "cognito",
        function: "functions/setupUserMFA/main.go",
      },
      "POST /accounts/mfa/disable": {
        authorizer: "cognito",
        function: "functions/disableUserMFA/main.go",
      },
      "POST /accounts/mfa/verify": {
        authorizer: "cognito",
        function: "functions/verifyUserMFA/main.go",
      },
    },
  });

  stack.addOutputs({
    API_ENDPOINT: api.url,
    API_ID: api.httpApiId,
    ...environment,
  });
}
