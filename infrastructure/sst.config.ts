import { SSTConfig } from "sst";

export default {
  config() {
    return {
      name: "jungyu-blog",
      region: "ap-northeast-2",
    };
  },
  stacks(app) {
    app.stack(BackendStack);
    app.stack(FrontendStack);
  },
} satisfies SSTConfig;

// 스택 import
import { BackendStack } from "./stacks/BackendStack";
import { FrontendStack } from "./stacks/FrontendStack";
