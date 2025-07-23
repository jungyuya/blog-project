// blog-frontend/src/components/Header.tsx
import Link from 'next/link';

export default function Header() {
  return (
    <header className="bg-blue-600 text-white p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-2xl font-bold">
          JUNGYU'S BLOG
        </Link>
        <nav>
          <Link href="/posts/new" className="bg-blue-700 hover:bg-blue-800 text-white font-bold py-2 px-4 rounded transition duration-300">
            새 글 작성
          </Link>
        </nav>
      </div>
    </header>
  );
}