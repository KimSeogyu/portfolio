import { json, redirect, type ActionFunctionArgs, LoaderFunctionArgs } from "@remix-run/node";
import { Form, useActionData, useLoaderData } from "@remix-run/react";
import { boardClient } from "~/proto/client";
import { PostingStatus } from "~/proto/board/v1/service_pb";


export async function loader({ params }: LoaderFunctionArgs) {
  try {
    const postingId = Number(params.id);
    if (isNaN(postingId)) {
      throw new Error("Invalid posting ID");
    }
    
    const post = await boardClient.getPosting({ postingId });
    
    // 게시글이 있는지 확인
    if (!post) {
      throw new Error("Posting not found");
    }
    
    return json({
      post,
      error: null
    });
  } catch (error) {
    console.error("Failed to load posting:", error);
    return json({ 
      post: null, 
      error: "게시글을 불러오는데 실패했습니다." 
    });
  }
}

export async function action({ request, params }: ActionFunctionArgs) {
  const postingId = Number(params.id);
  if (isNaN(postingId)) {
    return json({ errors: { form: "잘못된 게시글 ID입니다." } });
  }
  
  const formData = await request.formData();
  const title = formData.get("title") as string;
  const content = formData.get("content") as string;
  const tagsText = formData.get("tags") as string;
  
  // Validation
  const errors: Record<string, string> = {};
  if (!title || title.trim() === "") {
    errors.title = "제목을 입력해주세요.";
  }
  if (!content || content.trim() === "") {
    errors.content = "내용을 입력해주세요.";
  }
  
  if (Object.keys(errors).length > 0) {
    return json({ errors, values: { title, content, tagsText } });
  }
  
  // Parse tags
  const tags = tagsText
    ? tagsText.split(",").map(tag => tag.trim()).filter(Boolean)
    : [];
  
  try {
    await boardClient.updatePosting({
      postingId,
      title,
      content,
      tags,
      status: PostingStatus.PUBLISHED,
    });
    
    return redirect(`/board/${postingId}`);
  } catch (error) {
    console.error("Failed to update post:", error);
    return json({
      errors: { form: "게시글 수정에 실패했습니다." },
      values: { title, content, tagsText }
    });
  }
}

// Type guard to check if actionData has values property
function hasValues(data: any): data is { values: { title: string; content: string; tagsText: string } } {
  return data && 'values' in data;
}

// Add this helper function to check if a field has an error
function hasFieldError(errors: any, field: string): boolean {
  return errors && typeof errors === 'object' && field in errors;
}

export default function EditPost() {
  const { post, error } = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();
  
  if (error) {
    return <div className="text-red-500">{error}</div>;
  }
  
  if (!post) {
    return <div className="text-center py-8">게시글을 찾을 수 없습니다.</div>;
  }

  const errors = actionData?.errors || {};
  
  // Fixed version with type check
  const values = hasValues(actionData) ? actionData.values : { 
    title: post.title || "", 
    content: post.content || "", 
    tagsText: Array.isArray(post.tags) ? post.tags.join(", ") : ""
  };
  
  return (
    <div>
      <h2 className="text-xl font-semibold mb-6">게시글 수정</h2>
      
      {errors.form && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {errors.form}
        </div>
      )}
      
      <Form method="post" className="space-y-4">
        <div>
          <label htmlFor="title" className="block mb-1 font-medium">
            제목
          </label>
          <input
            type="text"
            id="title"
            name="title"
            className={`w-full p-2 border rounded ${hasFieldError(errors, 'title') ? 'border-red-500' : 'border-gray-300'}`}
            defaultValue={values.title}
          />
          {hasFieldError(errors, 'title') && (
            <p className="mt-1 text-red-500 text-sm">{(errors as any).title}</p>
          )}
        </div>
        
        <div>
          <label htmlFor="content" className="block mb-1 font-medium">
            내용
          </label>
          <textarea
            id="content"
            name="content"
            rows={10}
            className={`w-full p-2 border rounded ${hasFieldError(errors, 'content') ? 'border-red-500' : 'border-gray-300'}`}
            defaultValue={values.content}
          />
          {hasFieldError(errors, 'content') && (
            <p className="mt-1 text-red-500 text-sm">{(errors as any).content}</p>
          )}
        </div>
        
        <div>
          <label htmlFor="tags" className="block mb-1 font-medium">
            태그 (쉼표로 구분)
          </label>
          <input
            type="text"
            id="tags"
            name="tags"
            className="w-full p-2 border border-gray-300 rounded"
            placeholder="예: 일상, 질문, 정보"
            defaultValue={values.tagsText}
          />
          <p className="mt-1 text-gray-500 text-sm">
            태그는 쉼표(,)로 구분하여 입력해주세요.
          </p>
        </div>
        
        <div className="text-right">
          <button
            type="submit"
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            게시글 수정
          </button>
        </div>
      </Form>
    </div>
  );
} 