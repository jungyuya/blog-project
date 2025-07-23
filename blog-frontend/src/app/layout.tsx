// blog-frontend/src/app/layout.tsx
import type { Metadata } from "next";
import { Inter } from "next/font/google"; // 폰트 설정을 위해 Next.js가 자동으로 추가
import "./globals.css"; // 전역 CSS 파일 임포트

const inter = Inter({ subsets: ["latin"] }); // Next.js 폰트 최적화 기능

export const metadata: Metadata = {
  title: "Jungyu's Blog",
  description: "Jungyu's personal blog application built with Next.js and Go Lambda.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ko">
      <body className={inter.className}>
        {children}
      </body>
    </html>
  );
}