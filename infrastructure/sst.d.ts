// sst.d.ts
declare module "sst" {
  import { Construct } from "constructs";

  export interface StackContext {
    stack: any;
    app: any;
  }

  export class Api {
    constructor(scope: Construct, id: string, props: any);
  }

  export class NextjsSite {
    constructor(scope: Construct, id: string, props: any);
  }
}
