// src/app/posts/[id]/DeleteButton.tsx
"use client"; // 이 파일은 클라이언트 컴포넌트임을 명시

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { BACKEND_BASE_URL } from '@/config/backend_config';

interface DeleteButtonProps {
    postId: string;
}

export default function DeleteButton({ postId }: DeleteButtonProps) {
    const router = useRouter();
    const [loading, setLoading] = useState(false); // 삭제 로딩 상태
    const [error, setError] = useState<string | null>(null); // 삭제 오류 상태

    const handleDelete = async () => {
        const confirmDelete = confirm(`정말로 이 게시물을 삭제하시겠습니까?`);
        if (!confirmDelete) {
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const response = await fetch(`${BACKEND_BASE_URL}/posts/${postId}`, {
                method: 'DELETE',
            });

            if (!response.ok) {
                const errorText = await response.text();
                let errorMessage = `게시물 삭제 실패: ${response.status} ${response.statusText}`;
                try {
                    // eslint-disable-next-line @typescript-eslint/no-unused-vars
                    const errorJson = JSON.parse(errorText); // _e 경고 해결: 사용하지 않는 변수에 eslint-disable 적용
                    errorMessage = errorJson.message || errorMessage;
                } catch (parseError) { // _e 대신 명확한 변수명 사용
                    // JSON 파싱 실패 시 원본 텍스트 사용
                }
                throw new Error(errorMessage);
            }

            router.push('/');
            router.refresh();
        } catch (err) {
            console.error('게시물 삭제 중 오류 발생:', err);
            setError(err instanceof Error ? err.message : '게시물 삭제 중 알 수 없는 오류가 발생했습니다.');
        } finally {
            setLoading(false);
        }
    };

    if (error) {
        return <p className="text-red-500 text-sm mt-2">{error}</p>;
    }

    return (
        <button
            onClick={handleDelete}
            className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={loading}
        >
            {loading ? '삭제 중...' : '삭제'}
        </button>
    );
}