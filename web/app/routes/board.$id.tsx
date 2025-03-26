import { json, type LoaderFunctionArgs, ActionFunctionArgs, redirect } from "@remix-run/node";
import { useLoaderData, Link, useActionData, Form } from "@remix-run/react";
import { useState } from "react";
import { getBoardClient } from "~/utils/board-client.server";
import { Comment as PbComment } from "~/proto/board/v1/service_pb";
// 날짜 포맷 함수
function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

export async function loader({ params }: LoaderFunctionArgs) {
  try {
    const postingId = Number(params.id);
    if (isNaN(postingId)) {
      throw new Error("Invalid posting ID");
    }
    
    const boardClient = getBoardClient();
    const post = await boardClient.getPosting({ postingId });
    return json({ post, error: null });
  } catch (error) {
    console.error("Failed to load posting:", error);
    return json({ 
      post: null, 
      error: "게시글을 불러오는데 실패했습니다." 
    });
  }
}

export async function action({ request, params }: ActionFunctionArgs) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const postingId = Number(params.id);
  const boardClient = getBoardClient();
  
  // 댓글 작성
  if (intent === "addComment") {
    try {
      const content = formData.get("content") as string;
      const parentId = formData.get("parentId") ? Number(formData.get("parentId")) : 0;
      
      if (!content) {
        return json({ success: false, error: "댓글 내용을 입력해주세요." });
      }
      
      await boardClient.createComment({
        postingId,
        content,
        parentId: parentId > 0 ? parentId : undefined,
      });
      
      return json({ success: true });
    } catch (error) {
      console.error("Failed to add comment:", error);
      return json({ success: false, error: "댓글 작성에 실패했습니다." });
    }
  }
  
  // 댓글 삭제
  if (intent === "deleteComment") {
    try {
      const commentId = Number(formData.get("commentId"));
      
      await boardClient.deleteComment({
        postingId,
        commentId,
      });
      
      return json({ success: true });
    } catch (error) {
      console.error("Failed to delete comment:", error);
      return json({ success: false, error: "댓글 삭제에 실패했습니다." });
    }
  }
  
  // 게시글 삭제
  if (intent === "deletePost") {
    try {
      await boardClient.deletePosting({ postingId });
      return redirect("/board");
    } catch (error) {
      console.error("Failed to delete post:", error);
      return json({ success: false, error: "게시글 삭제에 실패했습니다." });
    }
  }
  
  return json({ success: false, error: "잘못된 요청입니다." });
}

function hasError(data: any): data is { error: string } {
  return data && typeof data.error === 'string';
}

function Comment({ comment, level = 0 }: { comment: PbComment; level?: number }) {
  const [replyOpen, setReplyOpen] = useState(false);
  
  return (
    <div className={`mt-2 pl-${level * 4}`}>
      <div className={`p-3 rounded ${level > 0 ? 'bg-gray-50 border-l-2 border-blue-300' : 'bg-gray-100'}`}>
        <div className="flex justify-between items-start">
          <div>
            <span className="font-medium">{comment.authorName}</span>
            <span className="text-xs text-gray-500 ml-2">{formatDate(comment.createdAt?.toString() || "")}</span>
          </div>
          <div className="flex space-x-2">
            <button
              type="button"
              className="text-sm text-blue-500 hover:underline"
              onClick={() => setReplyOpen(!replyOpen)}
            >
              {replyOpen ? '취소' : '답글'}
            </button>
            
            {/* 자신의 댓글일 경우에만 표시 */}
            <Form method="post" className="inline">
              <input type="hidden" name="commentId" value={comment.commentId.toString()} />
              <button
                type="submit"
                name="intent"
                value="deleteComment"
                className="text-sm text-red-500 hover:underline"
                onClick={(e) => {
                  if (!confirm('정말로 이 댓글을 삭제하시겠습니까?')) {
                    e.preventDefault();
                  }
                }}
              >
                삭제
              </button>
            </Form>
          </div>
        </div>
        <div className="mt-1">{comment.content}</div>
        
        {replyOpen && (
          <Form method="post" className="mt-3">
            <input type="hidden" name="parentId" value={comment.commentId.toString()} />
            <textarea
              name="content"
              className="w-full p-2 border rounded"
              placeholder="답글을 입력하세요..."
              rows={2}
              required
            />
            <div className="mt-2 text-right">
              <button
                type="submit"
                name="intent"
                value="addComment"
                className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 text-sm"
              >
                답글 등록
              </button>
            </div>
          </Form>
        )}
      </div>
      
      {/* 대댓글 표시 */}
      {comment.children && comment.children.length > 0 && (
        <div className="ml-8">
          {comment.children.map((childComment) => (
            <Comment key={childComment.commentId} comment={childComment} level={level + 1} />
          ))}
        </div>
      )}
    </div>
  );
}

export default function PostDetail() {
  const { post, error } = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();
  
  if (error) {
    return <div className="text-red-500">{error}</div>;
  }
  
  if (!post) {
    return <div className="text-center py-8">게시글을 찾을 수 없습니다.</div>;
  }
  
  return (
    <div>
      <article>
        <header className="mb-6">
          <h1 className="text-2xl font-bold">{post.title}</h1>
          <div className="flex justify-between items-center mt-2 text-gray-500 text-sm">
            <div>
              <span>작성자: {post.authorName}</span>
              <span className="mx-2">|</span>
              <span>작성일: {formatDate(post.createdAt)}</span>
              <span className="mx-2">|</span>
              <span>조회수: {post.viewCount}</span>
            </div>
            <div className="flex space-x-2">
              <Link
                to={`/board/${post.postingId}/edit`}
                className="text-blue-500 hover:underline"
              >
                수정
              </Link>
              <Form method="post" className="inline">
                <button
                  type="submit"
                  name="intent"
                  value="deletePost"
                  className="text-red-500 hover:underline"
                  onClick={(e) => {
                    if (!confirm('정말로 이 게시글을 삭제하시겠습니까?')) {
                      e.preventDefault();
                    }
                  }}
                >
                  삭제
                </button>
              </Form>
            </div>
          </div>
        </header>
        
        <div className="prose max-w-none mb-8 pb-8 border-b">
          {post.content}
        </div>
        
        <section className="mt-8">
          <h3 className="text-lg font-semibold mb-4">
            댓글 ({post.commentCount})
          </h3>
          
          {hasError(actionData) && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
              {actionData.error}
            </div>
          )}
          
          <Form method="post" className="mb-6">
            <textarea
              name="content"
              className="w-full p-3 border rounded"
              placeholder="댓글을 입력하세요..."
              rows={3}
              required
            />
            <div className="mt-2 text-right">
              <button
                type="submit"
                name="intent"
                value="addComment"
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                댓글 등록
              </button>
            </div>
          </Form>
          
          <div className="space-y-4">
            {post.comments.length === 0 ? (
              <p className="text-gray-500 text-center py-4">
                아직 댓글이 없습니다. 첫 댓글을 작성해보세요!
              </p>
            ) : (
              post.comments
                .filter((comment: PbComment) => !comment.parentId) // Root comments only
                .map((comment: PbComment) => (
                  <Comment key={comment.commentId} comment={comment} />
                ))
            )}
          </div>
        </section>
      </article>
    </div>
  );
} 