import { StackContext, Table } from "@serverless-stack/resources";

export function StorageStack({ stack }: StackContext) {
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

  return { table };
}
