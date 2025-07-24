// src/app/posts/[id]/edit/page.tsx
import { BACKEND_BASE_URL } from '@/config/backend_config';
import { notFound } from 'next/navigation';
import EditClient from './EditClient';

interface PostIdParam {
  id: string;
}

// 최대 20개만 미리 생성하도록 제한 추가
export async function generateStaticParams() {
  try {
    const res = await fetch(`${BACKEND_BASE_URL}/posts`);
    if (!res.ok) return [];
    const posts: { postId: string }[] = await res.json();

    return posts.slice(0, 20).map(p => ({ id: p.postId }));
  } catch {
    return [];
  }
}

export default async function EditPage({
  params,
}: {
  params: Promise<PostIdParam>;
}) {
  const { id } = await params;
  if (!id) notFound();
  return <EditClient postId={id} />;
}
