// blog-frontend/src/app/posts/[id]/edit/page.tsx
"use client";

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Layout from '@/components/Layout';
import { BACKEND_BASE_URL } from '@/config/backend_config';
import Link from 'next/link';

// 게시물 데이터 타입 정의 (백엔드와 일치)
interface Post {
  postId: string;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
}

export default function EditPostPage() {
  const { id } = useParams(); // URL에서 'id' 파라미터 가져오기 (이것이 postId)
  const router = useRouter();

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [author, setAuthor] = useState('');
  const [loading, setLoading] = useState(true); // 초기 로딩 상태는 true
  const [error, setError] = useState<string | null>(null);
  const [postFetched, setPostFetched] = useState(false); // 게시물 데이터를 가져왔는지 여부

  // 1. 컴포넌트 마운트 시 기존 게시물 데이터 불러오기
  useEffect(() => {
    if (!id || postFetched) { // id가 없거나 이미 데이터를 가져왔으면 실행하지 않음
      if (!id) setLoading(false); // id가 없으면 로딩 종료
      return;
    }

    const fetchPost = async () => {
      try {
        const response = await fetch(`${BACKEND_BASE_URL}/posts/${id}`);
        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Failed to fetch post for editing: ${response.status} ${response.statusText} - ${errorText}`);
        }
        const data: Post = await response.json();
        
        // 가져온 데이터로 폼 필드 초기화
        setTitle(data.title);
        setContent(data.content);
        setAuthor(data.author); // 작성자 필드도 초기화
        setPostFetched(true); // 데이터 가져옴 표시
      } catch (err) {
        console.error("Error fetching post for editing:", err);
        setError(err instanceof Error ? err.message : "게시물 정보를 불러오는 중 오류가 발생했습니다.");
      } finally {
        setLoading(false);
      }
    };

    fetchPost();
  }, [id, postFetched]); // id 또는 postFetched 상태가 변경될 때마다 실행

  // 2. 폼 제출 시 게시물 업데이트
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!title.trim() || !content.trim() || !author.trim()) {
      setError('제목, 내용, 작성자를 모두 입력해주세요.');
      return;
    }
    if (!id) {
      setError('게시물 ID를 찾을 수 없습니다.');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${BACKEND_BASE_URL}/posts/${id}`, {
        method: 'PUT', // PUT 메서드 사용
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title, content, author }), // 수정된 데이터 전송
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `게시물 수정 실패: ${response.status} ${response.statusText}`);
      }

      // alert('게시물이 성공적으로 수정되었습니다!'); // 사용자에게 알림
      router.push(`/posts/${id}`); // 수정된 게시물 상세 페이지로 이동
    } catch (err) {
      console.error('게시물 수정 중 오류 발생:', err);
      setError(err instanceof Error ? err.message : '게시물 수정 중 알 수 없는 오류가 발생했습니다.');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="text-center py-10">
          <p className="text-xl text-gray-600">게시물 정보를 불러오는 중입니다...</p>
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mt-4"></div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="text-center py-10 text-red-600">
          <p className="text-xl font-bold">오류 발생:</p>
          <p className="mt-2 text-lg">{error}</p>
          <Link href="/" className="mt-6 inline-block bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition duration-300">
            목록으로 돌아가기
          </Link>
        </div>
      </Layout>
    );
  }

  // 게시물 정보를 가져오지 못했거나 id가 없는 경우 (404 처리)
  if (!postFetched && !loading) { // 로딩이 끝났는데도 데이터가 없으면
    return (
      <Layout>
        <div className="text-center py-10 text-gray-600">
          <p className="text-xl font-bold">게시물을 찾을 수 없습니다.</p>
          <Link href="/" className="mt-6 inline-block bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded transition duration-300">
            목록으로 돌아가기
          </Link>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-3xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold text-gray-900 mb-6 text-center">게시물 수정</h1>
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
              disabled={loading}
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
              {loading ? '수정 중...' : '게시물 수정'}
            </button>
            <button
              type="button"
              onClick={() => router.back()}
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
