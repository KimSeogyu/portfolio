import { json, redirect, type ActionFunctionArgs } from "@remix-run/node";
import { Form, useActionData } from "@remix-run/react";
import { boardClient } from "~/proto/client";
import { PostingStatus } from "~/proto/board/v1/service_pb";

export async function action({ request }: ActionFunctionArgs) {
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
    const response = await boardClient.createPosting({
      title,
      content,
      tags,
      status: PostingStatus.PUBLISHED,
    });
    
    return redirect(`/board/${response.postingId}`);
  } catch (error) {
    console.error("Failed to create post:", error);
    return json({
      errors: { form: "게시글 작성에 실패했습니다." },
      values: { title, content, tagsText }
    });
  }
}

function getFieldError(errors: any, field: string): string | undefined {
  return errors && errors[field];
}

export default function NewPost() {
  const actionData = useActionData<typeof action>();
  const errors = actionData?.errors || {};
  const values = actionData?.values || { title: "", content: "", tagsText: "" };
  
  return (
    <div>
      <h2 className="text-xl font-semibold mb-6">새 게시글 작성</h2>
      
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
            className={`w-full p-2 border rounded ${getFieldError(errors, 'title') ? 'border-red-500' : 'border-gray-300'}`}
            defaultValue={values.title}
          />
          {getFieldError(errors, 'title') && (
            <p className="mt-1 text-red-500 text-sm">{getFieldError(errors, 'title')}</p>
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
            className={`w-full p-2 border rounded ${getFieldError(errors, 'content') ? 'border-red-500' : 'border-gray-300'}`}
            defaultValue={values.content}
          />
          {getFieldError(errors, 'content') && (
            <p className="mt-1 text-red-500 text-sm">{getFieldError(errors, 'content')}</p>
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
            게시글 등록
          </button>
        </div>
      </Form>
    </div>
  );
} 