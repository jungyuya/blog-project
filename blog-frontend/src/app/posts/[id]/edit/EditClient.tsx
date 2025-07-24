'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Layout from '@/components/Layout';
import { BACKEND_BASE_URL } from '@/config/backend_config';
import Link from 'next/link';

interface Post {
  postId: string;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
}

interface EditClientProps {
  postId: string;
}

export default function EditClient({ postId }: EditClientProps) {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [author, setAuthor] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 기존 게시물 불러오기
  useEffect(() => {
    async function fetchPost() {
      try {
        const res = await fetch(`${BACKEND_BASE_URL}/posts/${postId}`);
        if (!res.ok) throw new Error(`Failed to fetch: ${res.status}`);
        const data: Post = await res.json();
        setTitle(data.title);
        setContent(data.content);
        setAuthor(data.author);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error loading post');
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, [postId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim() || !author.trim()) {
      setError('모든 필드를 입력해주세요.');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${BACKEND_BASE_URL}/posts/${postId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, content, author }),
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.message || 'Update failed');
      }
      router.push(`/posts/${postId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error updating');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Layout>
        <div className="text-center py-10">
          <p>로딩 중...</p>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="text-center text-red-600 py-10">
          <p>{error}</p>
          <Link href="/" className="mt-4 inline-block text-blue-500">
            목록으로 돌아가기
          </Link>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-3xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-3xl font-bold mb-6 text-center">게시물 수정</h1>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="title" className="block mb-2 font-bold">제목</label>
            <input
              id="title"
              value={title}
              onChange={e => setTitle(e.target.value)}
              className="w-full border rounded p-2"
              disabled={loading}
              required
            />
          </div>

          <div className="mb-4">
            <label htmlFor="content" className="block mb-2 font-bold">내용</label>
            <textarea
              id="content"
              value={content}
              onChange={e => setContent(e.target.value)}
              className="w-full border rounded p-2 h-48"
              disabled={loading}
              required
            />
          </div>

          <div className="mb-6">
            <label htmlFor="author" className="block mb-2 font-bold">작성자</label>
            <input
              id="author"
              value={author}
              onChange={e => setAuthor(e.target.value)}
              className="w-full border rounded p-2"
              disabled={loading}
              required
            />
          </div>

          <div className="flex justify-between">
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
              disabled={loading}
            >
              {loading ? '수정 중...' : '수정'}
            </button>
            <button
              type="button"
              onClick={() => router.back()}
              className="bg-gray-300 hover:bg-gray-400 text-gray-800 py-2 px-4 rounded"
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
