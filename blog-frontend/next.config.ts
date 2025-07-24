// next.config.ts (blog-frontend)

import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  // app 라우터 사용 시 정적 export를 위한 설정
  output: 'export', // 이 줄을 추가합니다.

  // 이미 존재하거나 필요한 다른 설정들 (예시)
  reactStrictMode: true,
  images: {
    unoptimized: true, // 만약 Next/Image 컴포넌트 사용 시 최적화를 위해 서버가 필요하다면 이 줄을 제거하거나, S3/CloudFront에 맞게 loader 설정
  },
  // SWC (Speedy Web Compiler) 설정 (선택 사항)
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
};

export default nextConfig;