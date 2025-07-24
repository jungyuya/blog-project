// src/app/posts/[id]/edit/page.tsx
import { BACKEND_BASE_URL } from '@/config/backend_config';
import { notFound } from 'next/navigation';
import EditClient from './EditClient';

interface PostIdParam {
  id: string;
}

// export 시 동적 경로를 미리 생성
export async function generateStaticParams() {
  try {
    const res = await fetch(`${BACKEND_BASE_URL}/posts`, { cache: 'no-store' });
    if (!res.ok) return [];
    const posts: { postId: string }[] = await res.json();
    return posts.map(p => ({ id: p.postId }));
  } catch {
    return [];
  }
}

// async 로직만 두고, metadata는 일단 빼 버립니다.
export default async function EditPage({
  params,
}: {
  params: Promise<PostIdParam>;
}) {
  const { id } = await params;
  if (!id) notFound();
  return <EditClient postId={id} />;
}
