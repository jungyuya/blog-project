// src/app/posts/[id]/page.tsx
export const revalidate = 60;

import { notFound } from 'next/navigation';
import Layout from '@/components/Layout';

interface Post {
  postId: string;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
}

async function getPostById(id: string): Promise<Post | null> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/posts/${id}`, {
    next: { revalidate: 60 }
  });

  if (!res.ok) return null;
  return res.json();
}

export default async function PostPage({ params }: { params: { id: string } }) {
  const post = await getPostById(params.id);

  if (!post) return notFound();

  return (
    <Layout>
      <h1 className="text-3xl font-bold">{post.title}</h1>
      <p className="text-gray-500 text-sm mb-4">작성자: {post.author} | 작성일: {post.createdAt}</p>
      <div>{post.content}</div>
    </Layout>
  );
}