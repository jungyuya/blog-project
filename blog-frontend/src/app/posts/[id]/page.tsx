// blog-frontend/src/app/posts/[id]/page.tsx
"use client";

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
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

export default function PostDetailPage() {
  const { id } = useParams(); // URL에서 'id' 파라미터 가져오기 (이것이 postId)
  const router = useRouter();

  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      setError("게시물 ID가 제공되지 않았습니다.");
      return;
    }

    const fetchPost = async () => {
      try {
        const response = await fetch(`${BACKEND_BASE_URL}/posts/${id}`);
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('게시물을 찾을 수 없습니다.');
          }
          const errorText = await response.text();
          throw new Error(`게시물 불러오기 실패: ${response.status} ${response.statusText} - ${errorText}`);
        }
        const data: Post = await response.json();
        setPost(data);
      } catch (err) {
        console.error("Error fetching post:", err);
        setError(err instanceof Error ? err.message : "게시물 정보를 불러오는 중 오류가 발생했습니다.");
      } finally {
        setLoading(false);
      }
    };

    fetchPost();
  }, [id]);

  // ⭐ 게시물 삭제 핸들러 추가 ⭐
  const handleDelete = async () => {
    if (!post) return;

    // 사용자에게 삭제 여부 확인
    const confirmDelete = confirm(`정말로 "${post.title}" 게시물을 삭제하시겠습니까?`);
    if (!confirmDelete) {
      return; // 사용자가 취소하면 아무것도 하지 않음
    }

    setLoading(true); // 로딩 상태 시작
    setError(null); // 이전 오류 메시지 초기화

    try {
      const response = await fetch(`${BACKEND_BASE_URL}/posts/${post.postId}`, {
        method: 'DELETE', // DELETE 메서드 사용
      });

      if (!response.ok) {
        // DELETE 요청은 보통 본문이 없거나 간단한 메시지이므로, text()로 받아서 에러 메시지 파싱 시도
        const errorText = await response.text();
        let errorMessage = `게시물 삭제 실패: ${response.status} ${response.statusText}`;
        try {
          const errorJson = JSON.parse(errorText);
          errorMessage = errorJson.message || errorMessage;
        } catch (e) {
          // JSON 파싱 실패 시 원본 텍스트 사용
        }
        throw new Error(errorMessage);
      }

      // alert('게시물이 성공적으로 삭제되었습니다.'); // 사용자에게 알림 (선택 사항)
      router.push('/'); // 삭제 성공 후 목록 페이지로 이동
    } catch (err) {
      console.error('게시물 삭제 중 오류 발생:', err);
      setError(err instanceof Error ? err.message : '게시물 삭제 중 알 수 없는 오류가 발생했습니다.');
    } finally {
      setLoading(false); // 로딩 상태 종료
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

  if (!post) {
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
      <div className="max-w-4xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-4xl font-extrabold text-gray-900 mb-4 break-words">{post.title}</h1>
        <p className="text-sm text-gray-500 mb-2">
          작성자: {post.author}
        </p>
        <p className="text-sm text-gray-500 mb-6">
          작성일: {new Date(post.createdAt).toLocaleString()} | 최종 수정일: {new Date(post.updatedAt).toLocaleString()}
        </p>
        <div className="prose prose-lg max-w-none text-gray-800 leading-relaxed mb-8 break-words" style={{ whiteSpace: 'pre-wrap' }}>
          {post.content}
        </div>

        <div className="flex justify-end gap-4 mt-6">
          <Link href="/" className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed">
            목록으로
          </Link>
          <Link href={`/posts/${post.postId}/edit`} className="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-2 px-4 rounded transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed">
            수정
          </Link>
          {/* ⭐ 여기에 삭제 버튼을 추가합니다. ⭐ */}
          <button
            onClick={handleDelete}
            className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={loading} // 삭제 중에는 버튼 비활성화
          >
            {loading ? '삭제 중...' : '삭제'}
          </button>
        </div>
      </div>
    </Layout>
  );
}
