// blog-frontend/src/components/Layout.tsx
import Header from './Header';
import Footer from './Footer';
import React from 'react'; // React.ReactNode 타입을 위해 임포트

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="flex-grow container mx-auto px-4 py-8">
        {children}
      </main>
      <Footer />
    </div>
  );
}