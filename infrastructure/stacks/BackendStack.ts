import { StackContext, Api } from "sst/constructs";

export function BackendStack({ stack }: StackContext) {
  const api = new Api(stack, "BlogApi", {
    routes: {
      "GET    /posts": "blog-backend/getPosts",
      "GET    /post/{id}": "blog-backend/getPost",
      "POST   /posts": "blog-backend/createPost",
      "PUT    /post/{id}": "blog-backend/updatePost",
      "DELETE /post/{id}": "blog-backend/deletePost"
    },
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
