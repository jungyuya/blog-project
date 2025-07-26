# infrastructure/outputs.tf

output "backend_api_url" {
  value       = aws_api_gateway_rest_api.blog_api.execution_arn
  description = "API Gateway ARN or URL"
}

output "cloudfront_distribution_id" {
  value       = aws_cloudfront_distribution.blog_frontend_distribution.id
  description = "CloudFront Distribution ID"
}
