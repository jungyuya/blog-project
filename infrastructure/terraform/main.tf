# infrastructure/main.tf

# 서울 리전 (ap-northeast-2)을 기본 프로바이더로 설정
provider "aws" {
  region = "ap-northeast-2"
}

# us-east-1 (버지니아 북부) 리전에서 ACM 인증서 발급을 위한 Provider
# CloudFront는 us-east-1 리전의 인증서만 지원합니다.
provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

# jungyu.store 도메인의 Route 53 호스팅 영역 정보를 가져옵니다.
# Terraform이 자동으로 해당 호스팅 영역 ID를 참조할 수 있도록 합니다.
data "aws_route53_zone" "jungyu_store_hosted_zone" {
  name = "jungyu.store." # 도메인 이름 끝에 반드시 .을 붙여야 합니다.
}

# S3 버킷 생성 (블로그 프론트엔드 정적 파일 호스팅용)
resource "aws_s3_bucket" "blog_frontend_bucket" {
  # 이 버킷 이름은 AWS 내에서 전역적으로 고유해야 합니다.
  bucket = "jungyu-blog-frontend-20250722" 

  tags = {
    Name        = "JungyuBlogFrontendBucket"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# S3 버킷 웹사이트 호스팅 설정
resource "aws_s3_bucket_website_configuration" "blog_frontend_bucket_website_config" {
  bucket = aws_s3_bucket.blog_frontend_bucket.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html" # SPA (Single Page Application)에서 라우팅 처리 시 유용
  }
}

# CloudFront Origin Access Identity (OAI) 생성
# OAI는 CloudFront가 S3 버킷의 private 콘텐츠에 접근할 수 있도록 하는 역할입니다.
resource "aws_cloudfront_origin_access_identity" "blog_frontend_oai" {
  comment = "OAI for Jungyu Blog Frontend S3 Bucket"
}

# S3 버킷 정책 (CloudFront OAI가 접근할 수 있도록 허용)
resource "aws_s3_bucket_policy" "blog_frontend_bucket_policy" {
  bucket = aws_s3_bucket.blog_frontend_bucket.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.blog_frontend_oai.iam_arn
        },
        Action   = "s3:GetObject",
        Resource = "${aws_s3_bucket.blog_frontend_bucket.arn}/*"
      }
    ]
  })
}

# S3 버킷 퍼블릭 접근 차단 설정 (보안 강화)
resource "aws_s3_bucket_public_access_block" "blog_frontend_bucket_public_access_block" {
  bucket = aws_s3_bucket.blog_frontend_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# 블로그 이미지 저장용 S3 버킷 (별도)
resource "aws_s3_bucket" "blog_image_bucket" {
  # 이 버킷 이름 또한 AWS 내에서 전역적으로 고유해야 합니다.
  bucket = "jungyu-blog-image-storage-20250722" 

  tags = {
    Name        = "BlogImageStorage"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# 이미지 S3 버킷의 퍼블릭 접근 차단 설정
resource "aws_s3_bucket_public_access_block" "blog_image_bucket_public_access_block" {
  bucket = aws_s3_bucket.blog_image_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}


# ACM 인증서 발급 (HTTPS를 위해 필수이며, CloudFront는 us-east-1 리전 인증서만 지원)
resource "aws_acm_certificate" "blog_domain_cert" {
  provider          = aws.us_east_1 # us-east-1 리전에서 발급
  domain_name       = "blog.jungyu.store" # 실제 사용할 서브도메인 
  # subject_alternative_names = ["jungyu.store"] # 루트 도메인도 HTTPS로 사용하려면 이 줄의 주석을 해제합니다.
  validation_method = "DNS" # Route 53을 통한 DNS 유효성 검사

  tags = {
    Name = "JungyuBlogCertificate"
  }

  lifecycle {
    create_before_destroy = true # 인증서 교체 시 다운타임 방지
  }
}

# ACM 인증서 유효성 검사를 위한 Route 53 레코드 생성
# 이 레코드는 ACM이 인증서 소유권을 확인하기 위해 필요하며, Route 53에서 자동으로 생성/관리됩니다.
resource "aws_route53_record" "blog_domain_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.blog_domain_cert.domain_validation_options : dvo.domain_name => dvo
  }

  zone_id = data.aws_route53_zone.jungyu_store_hosted_zone.zone_id # 위에서 정의된 data 소스 참조
  name    = each.value.resource_record_name
  type    = each.value.resource_record_type
  records = [each.value.resource_record_value]
  ttl     = 60
}

# ACM 인증서 유효성 검증 (위의 레코드가 생성되고 전파된 것을 확인)
resource "aws_acm_certificate_validation" "blog_domain_cert_validation" {
  provider        = aws.us_east_1
  certificate_arn = aws_acm_certificate.blog_domain_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.blog_domain_cert_validation : record.fqdn]
  # depends_on = [aws_route53_record.blog_domain_cert_validation]
}


# CloudFront 배포 (CDN) 설정
resource "aws_cloudfront_distribution" "blog_frontend_distribution" {
  origin {
    domain_name = aws_s3_bucket.blog_frontend_bucket.bucket_regional_domain_name
    origin_id   = "S3-BlogFrontendBucket"

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.blog_frontend_oai.cloudfront_access_identity_path
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "CloudFront distribution for Jungyu Blog Frontend"
  default_root_object = "index.html" # SPA의 진입점

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-BlogFrontendBucket"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https" # HTTPS 강제
    min_ttl                = 0
    default_ttl            = 3600 # 1시간 캐시
    max_ttl                = 86400 # 24시간 캐시
    compress               = true
  }

  # 커스텀 도메인(blog.jungyu.store)과 ACM 인증서를 연결합니다.
  aliases = ["blog.jungyu.store"] # 실제 사용할 서브도메인을 여기에 추가합니다.
  viewer_certificate {
    cloudfront_default_certificate = false # 기본 인증서 대신 ACM 사용
    acm_certificate_arn      = aws_acm_certificate.blog_domain_cert.arn # 발급받은 ACM 인증서 ARN 연결
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2019"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none" # 모든 지역 허용
    }
  }

  tags = {
    Name        = "JungyuBlogFrontendDistribution"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# Route 53 Alias Record (커스텀 도메인 -> CloudFront 배포 연결)
# blog.jungyu.store 요청이 CloudFront로 라우팅되도록 설정합니다.
resource "aws_route53_record" "blog_frontend_alias" {
  zone_id = data.aws_route53_zone.jungyu_store_hosted_zone.zone_id # 위에서 가져온 호스팅 영역 ID 참조
  name    = "blog.jungyu.store" # 실제 사용할 서브 도메인을 여기에 추가합니다.
  type    = "A" # Alias Record는 A 타입으로 설정

  alias {
    name                   = aws_cloudfront_distribution.blog_frontend_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.blog_frontend_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

# Cognito User Pool (사용자 인증을 위한)
resource "aws_cognito_user_pool" "blog_user_pool" {
  name = "JungyuBlogUserPool"

  # 이메일로 로그인하도록 설정
  username_attributes = ["email"]
  # 이메일로 계정 확인하도록 설정
  auto_verified_attributes = ["email"]

  # 비밀번호 정책
  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
    require_uppercase = true
  }

  schema {
    name                = "email"
    attribute_data_type = "String"
    mutable             = true
    required            = true
  }

  tags = {
    Name        = "JungyuBlogUserPool"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# Cognito User Pool Client (웹/모바일 앱에서 Cognito와 통신하기 위한 클라이언트)
resource "aws_cognito_user_pool_client" "blog_user_pool_client" {
  name         = "JungyuBlogWebClient"
  user_pool_id = aws_cognito_user_pool.blog_user_pool.id

  # ID 토큰, 액세스 토큰, 리프레시 토큰 발급 활성화
  allowed_oauth_flows_user_pool_client = true
  # 모든 인증 흐름에 'ALLOW_' 접두사를 사용해야 합니다.
  explicit_auth_flows = ["ALLOW_ADMIN_USER_PASSWORD_AUTH", "ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
  supported_identity_providers         = ["COGNITO"]

  # Callback URL 및 Logout URL (Next.js 앱이 배포될 URL 또는 로컬 개발 URL)
  # 나중에 프론트엔드 배포 후 정확한 URL로 업데이트해야 합니다.
  # 현재는 임시로 blog.jungyu.store 도메인을 사용하도록 설정했습니다.
  callback_urls = ["http://localhost:3000/callback", "https://blog.jungyu.store/callback"]
  logout_urls   = ["http://localhost:3000/logout", "https://blog.jungyu.store/logout"]

  # OAuth2 흐름 활성화
  allowed_oauth_flows  = ["code", "implicit"]
  allowed_oauth_scopes = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]
}

# Lambda가 사용할 기본적인 IAM Role (Serverless Framework에서 이 Role을 참조할 예정)
resource "aws_iam_role" "lambda_exec_role" {
  name = "JungyuBlogLambdaExecutionRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "JungyuBlogLambdaExecutionRole"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# Lambda 실행 Role에 기본 CloudWatch Logs 권한 부여
resource "aws_iam_role_policy_attachment" "lambda_logs_attachment" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# DynamoDB 테이블 (블로그 게시물 저장용)
resource "aws_dynamodb_table" "blog_posts_table" {
  name         = "JungyuBlogPosts"
  billing_mode = "PAY_PER_REQUEST" # 온디맨드 용량 (초기 개발에 적합)
  hash_key     = "postId" # 파티션 키

  attribute {
    name = "postId"
    type = "S" # String
  }

  tags = {
    Name        = "JungyuBlogPostsTable"
    Environment = "Dev"
    Project     = "Blog"
  }
}

# Lambda 실행 Role에 DynamoDB 권한 부여
resource "aws_iam_role_policy" "lambda_dynamodb_policy" {
  name   = "JungyuBlogLambdaDynamoDBAccess"
  role   = aws_iam_role.lambda_exec_role.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ],
        Resource = [
          aws_dynamodb_table.blog_posts_table.arn
        ]
      }
    ]
  })
}

# Terraform outputs (배포된 리소스의 정보를 쉽게 확인하기 위함)
output "blog_frontend_bucket_name" {
  value       = aws_s3_bucket.blog_frontend_bucket.id
  description = "Name of the S3 bucket for blog frontend."
}

output "blog_image_bucket_name" {
  value       = aws_s3_bucket.blog_image_bucket.id
  description = "Name of the S3 bucket for blog images."
}

output "cognito_user_pool_id" {
  value       = aws_cognito_user_pool.blog_user_pool.id
  description = "ID of the Cognito User Pool."
}

output "cognito_user_pool_client_id" {
  value       = aws_cognito_user_pool_client.blog_user_pool_client.id
  description = "ID of the Cognito User Pool Client."
}

output "lambda_execution_role_arn" {
  value       = aws_iam_role.lambda_exec_role.arn
  description = "ARN of the IAM role for Lambda execution."
}

output "cloudfront_domain_name" {
  value       = aws_cloudfront_distribution.blog_frontend_distribution.domain_name
  description = "CloudFront distribution domain name (AWS-generated)."
}

output "cloudfront_id" {
  value       = aws_cloudfront_distribution.blog_frontend_distribution.id
  description = "CloudFront distribution ID."
}

output "blog_cloudfront_url" {
  value       = "https://${aws_cloudfront_distribution.blog_frontend_distribution.domain_name}"
  description = "URL of the blog frontend via CloudFront's default domain."
}

output "blog_custom_domain_url" {
  value       = "https://blog.jungyu.store"
  description = "Custom domain URL for the blog frontend (blog.jungyu.store)."
}

output "acm_certificate_arn" {
  value       = aws_acm_certificate.blog_domain_cert.arn
  description = "ARN of the ACM certificate for blog.jungyu.store."
}

output "dynamodb_posts_table_name" {
  value       = aws_dynamodb_table.blog_posts_table.name
  description = "Name of the DynamoDB table for blog posts."
}