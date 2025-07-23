// blog-frontend/src/app/posts/new/page.tsx
"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation'; // 페이지 이동을 위해
import Layout from '@/components/Layout';
import { BACKEND_BASE_URL } from '@/config/backend_config';

export default function NewPostPage() {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [author, setAuthor] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter(); // useRouter 훅 초기화

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); // 기본 폼 제출 동작 방지

    // 간단한 유효성 검사
    if (!title.trim() || !content.trim() || !author.trim()) {
      setError('제목, 내용, 작성자를 모두 입력해주세요.');
      return;
    }

    setLoading(true);
    setError(null); // 이전 에러 메시지 초기화

    try {
      const response = await fetch(`${BACKEND_BASE_URL}/posts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title, content, author }),
      });

      if (!response.ok) {
        const errorData = await response.json(); // 백엔드에서 JSON 에러 응답이 올 경우
        throw new Error(errorData.message || `게시물 작성 실패: ${response.status} ${response.statusText}`);
      }

      const newPost = await response.json(); // 생성된 게시물 정보 (예: postId 포함)
      alert('게시물이 성공적으로 작성되었습니다!');
      router.push(`/posts/${newPost.postId}`); // 생성된 게시물의 상세 페이지로 이동
      // 또는 router.push('/'); // 목록 페이지로 이동
    } catch (err) {
      console.error('게시물 작성 중 오류 발생:', err);
      setError(err instanceof Error ? err.message : '게시물 작성 중 알 수 없는 오류가 발생했습니다.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-3xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-900 mb-6 text-center">새 게시물 작성</h1>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="title" className="block text-gray-700 text-sm font-bold mb-2">
              제목:
            </label>
            <input
              type="text"
              id="title"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={loading} // 로딩 중에는 입력 비활성화
              required
            />
          </div>

          <div className="mb-4">
            <label htmlFor="content" className="block text-gray-700 text-sm font-bold mb-2">
              내용:
            </label>
            <textarea
              id="content"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-48 resize-y"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={loading}
              required
            ></textarea>
          </div>

          <div className="mb-6">
            <label htmlFor="author" className="block text-gray-700 text-sm font-bold mb-2">
              작성자:
            </label>
            <input
              type="text"
              id="author"
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              disabled={loading}
              required
            />
          </div>

          {error && (
            <p className="text-red-500 text-sm italic mb-4 text-center">{error}</p>
          )}

          <div className="flex items-center justify-between">
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={loading}
            >
              {loading ? '작성 중...' : '게시물 작성'}
            </button>
            <button
              type="button"
              onClick={() => router.back()} // 이전 페이지로 돌아가기
              className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={loading}
            >
              취소
            </button>
          </div>
        </form>
      </div>
    </Layout>
  );
}