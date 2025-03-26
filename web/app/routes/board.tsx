import { Outlet, Link } from "@remix-run/react";

export default function BoardLayout() {
  return (
    <div className="max-w-4xl mx-auto p-4">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-4">게시판</h1>
        <div className="flex items-center space-x-4">
          <Link
            to="/board"
            className="text-blue-500 hover:underline"
          >
            목록으로
          </Link>
          <Link
            to="/board/new"
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            새 글 작성
          </Link>
        </div>
      </div>
      <div className="bg-white shadow-md rounded-lg p-6">
        <Outlet />
      </div>
    </div>
  );
} 