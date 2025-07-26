// src/app/page.tsx
export const revalidate = 60;

import Layout from '@/components/Layout';
import PostItem from '@/components/PostItem';

interface Post {
  postId: string;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
}

async function getPosts(): Promise<Post[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/posts`, {
    next: { revalidate: 60 }
  });

  if (!res.ok) throw new Error('게시물 데이터를 가져오는데 실패했습니다.');
  return res.json();
}

export default async function Home() {
  let posts: Post[] = [];

  try {
    posts = await getPosts();
  } catch (error) {
    return (
      <Layout>
        <p className="text-red-500 text-center py-10">데이터 로딩 실패: {(error as Error).message}</p>
      </Layout>
    );
  }

  return (
    <Layout>
      <h1 className="text-4xl font-bold text-center mb-8">최신 게시물</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {posts.length > 0 ? (
          posts.map((post) => (
            <PostItem key={post.postId} post={post} />
          ))
        ) : (
          <p className="text-center text-gray-600 text-lg col-span-full">게시물이 없습니다. 새로운 게시물을 작성해보세요!</p>
        )}
      </div>
    </Layout>
  );
}