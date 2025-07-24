import { BACKEND_BASE_URL } from '@/config/backend_config';
import Link from 'next/link';
import { notFound } from 'next/navigation';
import DeleteButton from './DeleteButton';

interface Post {
  postId: string;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
}

interface PostIdParam {
  id: string;
}

// 1) generateStaticParams에 개수 제한을 두어 미리 렌더링할 경로를 줄임 (최대 20개)
export async function generateStaticParams() {
  try {
    const res = await fetch(`${BACKEND_BASE_URL}/posts`);
    if (!res.ok) return [];
    const posts: Post[] = await res.json();

    // 최대 20개 포스트만 미리 렌더링
    return posts.slice(0, 20).map((p) => ({ id: p.postId }));
  } catch {
    return [];
  }
}

// 2) 페이지 컴포넌트: 동적 경로에 따른 상세 페이지
export default async function PostDetailPage({
  params,
}: {
  params: Promise<PostIdParam>;
}) {
  const { id } = await params;
  if (!id) notFound();

  let post: Post | null = null;
  let error: string | null = null;

  try {
    const res = await fetch(`${BACKEND_BASE_URL}/posts/${id}`);
    if (!res.ok) {
      if (res.status === 404) notFound();
      const txt = await res.text();
      throw new Error(`Fetch Error: ${res.status} ${txt}`);
    }
    post = await res.json();
  } catch (err) {
    console.error(err);
    error = err instanceof Error ? err.message : '알 수 없는 오류';
  }

  if (error) {
    return (
      <div className="text-center py-10 text-red-600">
        <p className="text-xl font-bold">오류 발생:</p>
        <p className="mt-2 text-lg">{error}</p>
        <Link
          href="/"
          className="mt-6 inline-block bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
        >
          목록으로 돌아가기
        </Link>
      </div>
    );
  }

  if (!post) notFound();

  return (
    <div className="max-w-4xl mx-auto bg-white p-8 rounded-lg shadow-lg">
      <h1 className="text-4xl font-extrabold text-gray-900 mb-4 break-words">
        {post.title}
      </h1>
      <p className="text-sm text-gray-500 mb-2">작성자: {post.author}</p>
      <p className="text-sm text-gray-500 mb-6">
        작성일: {new Date(post.createdAt).toLocaleString()} | 최종 수정일:{' '}
        {new Date(post.updatedAt).toLocaleString()}
      </p>
      <div
        className="prose prose-lg max-w-none text-gray-800 leading-relaxed mb-8 break-words"
        style={{ whiteSpace: 'pre-wrap' }}
      >
        {post.content}
      </div>
      <div className="flex justify-end gap-4 mt-6">
        <Link
          href="/"
          className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded"
        >
          목록으로
        </Link>
        <Link
          href={`/posts/${post.postId}/edit`}
          className="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-2 px-4 rounded"
        >
          수정
        </Link>
        <DeleteButton postId={post.postId} />
      </div>
    </div>
  );
}
