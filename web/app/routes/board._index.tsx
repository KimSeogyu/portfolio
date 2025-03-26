import { Timestamp } from "@bufbuild/protobuf/wkt";
import { json, type LoaderFunctionArgs } from "@remix-run/node";
import { useLoaderData, Link } from "@remix-run/react";
import { useState } from "react";
import { boardClient } from "~/proto/client";
import { Posting } from "~/proto/board/v1/service_pb";

// 날짜 포맷 함수
function formatDate(dateString?: Timestamp): string {
  if (!dateString) {
    return "";
  }
  const nanos = dateString.nanos;
  const seconds = dateString.seconds;
  const date = new Date(Number(seconds) * 1000 + Number(nanos) / 1000000);
  return date.toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export async function loader({ request }: LoaderFunctionArgs) {
  try {
    const url = new URL(request.url);
    const pageSize = 10;
    const pageToken = url.searchParams.get("pageToken") || "";
    
    const response = await boardClient.listPostings({
      pageSize,
      pageToken,
    });
    
    return json({
      postings: response.postings,
      nextPageToken: response.nextPageToken,
      totalCount: response.totalCount,
    });
  } catch (error) {
    console.error("Failed to load postings:", error);
    return json({ postings: [], nextPageToken: "", totalCount: 0, error: "게시글을 불러오는데 실패했습니다." });
  }
}

export default function BoardIndex() {
  const { postings, nextPageToken, totalCount, error } = useLoaderData<typeof loader>();
  const [currentPageToken, setCurrentPageToken] = useState("");
  
  if (error) {
    return <div className="text-red-500">{error}</div>;
  }

  return (
    <div>
      <h2 className="text-xl font-semibold mb-6">게시글 목록 (총 {totalCount}개)</h2>
      
      {postings.length === 0 ? (
        <p className="text-gray-500 text-center py-8">게시글이 없습니다.</p>
      ) : (
        <>
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr className="bg-gray-100">
                  <th className="px-4 py-2 text-left">번호</th>
                  <th className="px-4 py-2 text-left">제목</th>
                  <th className="px-4 py-2 text-left">작성자</th>
                  <th className="px-4 py-2 text-left">작성일</th>
                  <th className="px-4 py-2 text-left">조회수</th>
                </tr>
              </thead>
              <tbody>
                {postings.map((post: Posting) => (
                  <tr key={post.postingId} className="border-b hover:bg-gray-50">
                    <td className="px-4 py-3">{post.postingId.toString()}</td>
                    <td className="px-4 py-3">
                      <Link to={`/board/${post.postingId}`} className="text-blue-600 hover:underline">
                        {post.title} 
                        {post.commentCount > 0 && (
                          <span className="ml-2 text-gray-500">({post.commentCount})</span>
                        )}
                      </Link>
                    </td>
                    <td className="px-4 py-3">{post.authorName}</td>
                    <td className="px-4 py-3">{formatDate(post.createdAt)}</td>
                    <td className="px-4 py-3">{post.viewCount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          
          {nextPageToken && (
            <div className="mt-6 text-center">
              <Link
                to={`/board?pageToken=${nextPageToken}`}
                className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                onClick={() => setCurrentPageToken(nextPageToken)}
              >
                더보기
              </Link>
            </div>
          )}
        </>
      )}
    </div>
  );
} 