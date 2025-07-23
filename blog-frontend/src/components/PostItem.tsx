// blog-frontend/src/components/PostItem.tsx
import Link from 'next/link';

// PostItemProps 인터페이스를 백엔드 응답 형식에 맞게 수정
interface PostItemProps {
  post: {
    postId: string; // 'id' 대신 'postId'로 변경
    title: string;
    content: string;
    author: string;
    createdAt: string;
  };
}

export default function PostItem({ post }: PostItemProps) {
  // DEBUG: console.log 추가 (이제 postId가 제대로 찍히는지 확인)
  console.log("Rendering PostItem for PostID:", post.postId, "Post:", post);

  // 날짜 포맷팅 (선택 사항이지만 사용자 경험 개선)
  const formattedDate = new Date(post.createdAt).toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  return (
    <div className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 mb-6">
      {/* Link href에 post.postId를 사용하도록 수정 */}
      <Link href={`/posts/${post.postId}`} className="block">
        <h2 className="text-2xl font-bold text-gray-800 mb-2 hover:text-blue-600 transition-colors duration-200">
          {post.title}
        </h2>
      </Link>
      <p className="text-gray-600 text-sm mb-3">
        작성자: {post.author} | 날짜: {formattedDate}
      </p>
      <p className="text-gray-700 text-base line-clamp-3">
        {post.content}
      </p>
      <div className="mt-4 text-right">
        {/* Link href에 post.postId를 사용하도록 수정 */}
        <Link href={`/posts/${post.postId}`} className="text-blue-500 hover:underline">
          더 보기 →
        </Link>
      </div>
    </div>
  );
}