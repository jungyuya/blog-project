import { StackContext, NextjsSite } from "sst/constructs";

export function FrontendStack({ stack }: StackContext) {
  const site = new NextjsSite(stack, "FrontendSite", {
    path: "../blog-frontend",
  });

  stack.addOutputs({
    FrontendUrl: site.url,
  });
}
