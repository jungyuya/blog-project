# infrastructure/backend.tf
terraform {
  backend "s3" {
    bucket         = "blog-terraform-state-bucket"
    key            = "blog-project/terraform.tfstate"
    region         = "ap-northeast-2" # 서울 리전
    encrypt        = true # 데이터 암호화
  }
}