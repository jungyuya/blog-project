// blog-frontend/src/app/page.tsx
"use client";

import { useState, useEffect } from 'react';
import Layout from '@/components/Layout';
import PostItem from '@/components/PostItem';
import { BACKEND_BASE_URL } from '@/config/backend_config';

// Post 인터페이스를 백엔드 응답 형식에 맞게 수정
interface Post {
  postId: string; // 'id' 대신 'postId'로 변경
  title: string;
  content: string;
  author: string; // author 필드는 백엔드에서 아직 처리 안 했지만, 프론트엔드에서 표시를 위해 추가
  createdAt: string;
  updatedAt: string;
}

export default function Home() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await fetch(`${BACKEND_BASE_URL}/posts`);
        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Failed to fetch posts: ${response.status} ${response.statusText} - ${errorText}`);
        }
        // 백엔드에서 반환되는 데이터가 배열임을 가정
        const data: Post[] = await response.json();
        setPosts(data);
      } catch (err) {
        console.error("Error fetching posts:", err);
        setError(err instanceof Error ? err.message : "An unexpected error occurred.");
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, []);

  if (loading) {
    return (
      <Layout>
        <div className="text-center py-10">
          <p className="text-xl text-gray-600">게시물을 불러오는 중입니다...</p>
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
          <p className="mt-4 text-gray-500">API Gateway URL이 올바른지 확인해주세요.</p>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <h1 className="text-4xl font-bold text-center mb-8">최신 게시물</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {posts.length > 0 ? (
          posts.map((post) => (
            // key와 post prop에 post.postId를 사용하도록 수정
            <PostItem key={post.postId} post={post} />
          ))
        ) : (
          <p className="text-center text-gray-600 text-lg col-span-full">게시물이 없습니다. 새로운 게시물을 작성해보세요!</p>
        )}
      </div>
    </Layout>
  );
}