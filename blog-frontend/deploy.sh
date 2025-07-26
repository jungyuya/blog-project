#!/usr/bin/env bash
# blog-frontend/deploy.sh

set -e

FRONTEND_DIR="$(cd "$(dirname "$0")" && pwd)"
S3_BUCKET="jungyu-blog-frontend-20250722"
DISTRIBUTION_ID=$(aws ssm get-parameter \
    --name "/jungyu/blog/dev/CLOUDFRONT_DISTRIBUTION_ID" \
    --query "Parameter.Value" \
    --output text)

# Next.js 빌드
cd "$FRONTEND_DIR"
npm install
npm run build

# S3에 업로드
aws s3 sync ./out "s3://${S3_BUCKET}" --delete

# CloudFront 캐시 무효화
aws cloudfront create-invalidation \
  --distribution-id "${DISTRIBUTION_ID}" \
  --paths "/*"
