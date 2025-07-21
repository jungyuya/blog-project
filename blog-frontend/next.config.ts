

import type { NextConfig } from "next";

const nextConfig = {
  output: 'export', // 이 줄을 추가합니다.
  reactStrictMode: true,
  images: {
    unoptimized: true, // S3/CloudFront 배포 시 이미지 최적화 이슈 방지를 위해 설정 (필요 시)
  },
  // 다른 Next.js 설정들...
};

export default nextConfig;
