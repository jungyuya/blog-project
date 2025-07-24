'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function PostCreatedPage() {
  const router = useRouter();

  useEffect(() => {
    const timeout = setTimeout(() => {
      router.push('/'); // 글 목록으로 이동
    }, 3000); // 3초 후 자동 이동

    return () => clearTimeout(timeout); // cleanup
  }, [router]);

  return (
    <main className="max-w-xl mx-auto px-4 py-8 text-center">
      <h1 className="text-2xl font-bold mb-4">✅ 글이 성공적으로 작성되었습니다!</h1>
      <p className="mb-2 text-gray-700">
        글 목록으로 이동 중입니다. 잠시만 기다려주세요...
      </p>
      <p className="text-sm text-gray-500">
        글 작성 후 서버에 반영되기까지 <strong>1~2분 내외의 시간이 소요됩니다.</strong><br />
        이유가 궁금하다면&nbsp;
        <a
          href="https://jungyu.store"
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-500 underline"
        >
          이 링크
        </a>
        를 참조하세요.
      </p>
    </main>
  );
}
