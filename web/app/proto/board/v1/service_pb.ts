// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file board/v1/service.proto (package board.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { EmptySchema, Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_empty, file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_api_annotations } from "../../google/api/annotations_pb";
import { file_gnostic_openapi_v3_annotations } from "../../gnostic/openapi/v3/annotations_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file board/v1/service.proto.
 */
export const file_board_v1_service: GenFile = /*@__PURE__*/
  fileDesc("ChZib2FyZC92MS9zZXJ2aWNlLnByb3RvEghib2FyZC52MSL8AgoHUG9zdGluZxISCgpwb3N0aW5nX2lkGAEgASgDEg0KBXRpdGxlGAIgASgJEg8KB2NvbnRlbnQYAyABKAkSEQoJYXV0aG9yX2lkGAQgASgJEhMKC2F1dGhvcl9uYW1lGAUgASgJEi4KCmNyZWF0ZWRfYXQYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEi4KCnVwZGF0ZWRfYXQYByABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEi4KCmRlbGV0ZWRfYXQYCCABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEhIKCnZpZXdfY291bnQYCSABKAUSFQoNY29tbWVudF9jb3VudBgKIAEoBRIMCgR0YWdzGAsgAygJEicKBnN0YXR1cxgMIAEoDjIXLmJvYXJkLnYxLlBvc3RpbmdTdGF0dXMSIwoIY29tbWVudHMYDSADKAsyES5ib2FyZC52MS5Db21tZW50ItsCCgdDb21tZW50EhIKCmNvbW1lbnRfaWQYASABKAMSEgoKcG9zdGluZ19pZBgCIAEoAxIPCgdjb250ZW50GAMgASgJEhEKCWF1dGhvcl9pZBgEIAEoCRITCgthdXRob3JfbmFtZRgFIAEoCRIuCgpjcmVhdGVkX2F0GAYgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIuCgp1cGRhdGVkX2F0GAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIuCgpkZWxldGVkX2F0GAggASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIRCglwYXJlbnRfaWQYCSABKAMSJwoGc3RhdHVzGAogASgOMhcuYm9hcmQudjEuQ29tbWVudFN0YXR1cxIjCghjaGlsZHJlbhgLIAMoCzIRLmJvYXJkLnYxLkNvbW1lbnQigAEKFENyZWF0ZVBvc3RpbmdSZXF1ZXN0Eg0KBXRpdGxlGAEgASgJEg8KB2NvbnRlbnQYAiABKAkSEQoJYXV0aG9yX2lkGAMgASgJEgwKBHRhZ3MYBCADKAkSJwoGc3RhdHVzGAUgASgOMhcuYm9hcmQudjEuUG9zdGluZ1N0YXR1cyInChFHZXRQb3N0aW5nUmVxdWVzdBISCgpwb3N0aW5nX2lkGAEgASgDIqcBChNMaXN0UG9zdGluZ3NSZXF1ZXN0EhEKCXBhZ2Vfc2l6ZRgBIAEoBRISCgpwYWdlX3Rva2VuGAIgASgJEiEKB3NvcnRfYnkYAyABKA4yEC5ib2FyZC52MS5Tb3J0QnkSLwoOc29ydF9kaXJlY3Rpb24YBCABKA4yFy5ib2FyZC52MS5Tb3J0RGlyZWN0aW9uEhUKDWZpbHRlcl9ieV90YWcYBSABKAkiaQoUTGlzdFBvc3RpbmdzUmVzcG9uc2USIwoIcG9zdGluZ3MYASADKAsyES5ib2FyZC52MS5Qb3N0aW5nEhcKD25leHRfcGFnZV90b2tlbhgCIAEoCRITCgt0b3RhbF9jb3VudBgDIAEoBSKBAQoUVXBkYXRlUG9zdGluZ1JlcXVlc3QSEgoKcG9zdGluZ19pZBgBIAEoAxINCgV0aXRsZRgCIAEoCRIPCgdjb250ZW50GAMgASgJEgwKBHRhZ3MYBCADKAkSJwoGc3RhdHVzGAUgASgOMhcuYm9hcmQudjEuUG9zdGluZ1N0YXR1cyIqChREZWxldGVQb3N0aW5nUmVxdWVzdBISCgpwb3N0aW5nX2lkGAEgASgDIrcCChVTZWFyY2hQb3N0aW5nc1JlcXVlc3QSDAoEdGV4dBgBIAEoCRIWCg50aXRsZV9jb250YWlucxgCIAEoCRIYChBjb250ZW50X2NvbnRhaW5zGAMgASgJEgwKBHRhZ3MYBCADKAkSEQoJYXV0aG9yX2lkGAUgASgJEjEKDWNyZWF0ZWRfYWZ0ZXIYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEjIKDmNyZWF0ZWRfYmVmb3JlGAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIRCglwYWdlX3NpemUYCCABKAUSEgoKcGFnZV90b2tlbhgJIAEoCRIvCg5zb3J0X2RpcmVjdGlvbhgKIAEoDjIXLmJvYXJkLnYxLlNvcnREaXJlY3Rpb24iYQoUQ3JlYXRlQ29tbWVudFJlcXVlc3QSEgoKcG9zdGluZ19pZBgBIAEoAxIPCgdjb250ZW50GAIgASgJEhEKCWF1dGhvcl9pZBgDIAEoCRIRCglwYXJlbnRfaWQYBCABKAMiigEKHExpc3RDb21tZW50c0J5UG9zdGluZ1JlcXVlc3QSEgoKcG9zdGluZ19pZBgBIAEoAxIRCglwYWdlX3NpemUYAiABKAUSEgoKcGFnZV90b2tlbhgDIAEoCRIvCg5zb3J0X2RpcmVjdGlvbhgEIAEoDjIXLmJvYXJkLnYxLlNvcnREaXJlY3Rpb24iaQoUTGlzdENvbW1lbnRzUmVzcG9uc2USIwoIY29tbWVudHMYASADKAsyES5ib2FyZC52MS5Db21tZW50EhcKD25leHRfcGFnZV90b2tlbhgCIAEoCRITCgt0b3RhbF9jb3VudBgDIAEoBSJPChRVcGRhdGVDb21tZW50UmVxdWVzdBISCgpjb21tZW50X2lkGAEgASgDEhIKCnBvc3RpbmdfaWQYAiABKAMSDwoHY29udGVudBgDIAEoCSI+ChREZWxldGVDb21tZW50UmVxdWVzdBISCgpjb21tZW50X2lkGAEgASgDEhIKCnBvc3RpbmdfaWQYAiABKAMqgwEKDVBvc3RpbmdTdGF0dXMSHgoaUE9TVElOR19TVEFUVVNfVU5TUEVDSUZJRUQQABIcChhQT1NUSU5HX1NUQVRVU19QVUJMSVNIRUQQARIYChRQT1NUSU5HX1NUQVRVU19EUkFGVBACEhoKFlBPU1RJTkdfU1RBVFVTX0RFTEVURUQQAyppCg1Db21tZW50U3RhdHVzEh4KGkNPTU1FTlRfU1RBVFVTX1VOU1BFQ0lGSUVEEAASHAoYQ09NTUVOVF9TVEFUVVNfUFVCTElTSEVEEAESGgoWQ09NTUVOVF9TVEFUVVNfREVMRVRFRBACKmkKBlNvcnRCeRIXChNTT1JUX0JZX1VOU1BFQ0lGSUVEEAASFgoSU09SVF9CWV9DUkVBVEVEX0FUEAESFgoSU09SVF9CWV9VUERBVEVEX0FUEAISFgoSU09SVF9CWV9WSUVXX0NPVU5UEAMqYAoNU29ydERpcmVjdGlvbhIeChpTT1JUX0RJUkVDVElPTl9VTlNQRUNJRklFRBAAEhYKElNPUlRfRElSRUNUSU9OX0FTQxABEhcKE1NPUlRfRElSRUNUSU9OX0RFU0MQAjL7CAoMQm9hcmRTZXJ2aWNlElsKDUNyZWF0ZVBvc3RpbmcSHi5ib2FyZC52MS5DcmVhdGVQb3N0aW5nUmVxdWVzdBoRLmJvYXJkLnYxLlBvc3RpbmciF4LT5JMCEToBKiIML3YxL3Bvc3RpbmdzEl8KCkdldFBvc3RpbmcSGy5ib2FyZC52MS5HZXRQb3N0aW5nUmVxdWVzdBoRLmJvYXJkLnYxLlBvc3RpbmciIYLT5JMCGxIZL3YxL3Bvc3RpbmdzL3twb3N0aW5nX2lkfRJjCgxMaXN0UG9zdGluZ3MSHS5ib2FyZC52MS5MaXN0UG9zdGluZ3NSZXF1ZXN0Gh4uYm9hcmQudjEuTGlzdFBvc3RpbmdzUmVzcG9uc2UiFILT5JMCDhIML3YxL3Bvc3RpbmdzEosBChVMaXN0Q29tbWVudHNCeVBvc3RpbmcSJi5ib2FyZC52MS5MaXN0Q29tbWVudHNCeVBvc3RpbmdSZXF1ZXN0Gh4uYm9hcmQudjEuTGlzdENvbW1lbnRzUmVzcG9uc2UiKoLT5JMCJBIiL3YxL3Bvc3RpbmdzL3twb3N0aW5nX2lkfS9jb21tZW50cxJxCg1DcmVhdGVDb21tZW50Eh4uYm9hcmQudjEuQ3JlYXRlQ29tbWVudFJlcXVlc3QaES5ib2FyZC52MS5Db21tZW50Ii2C0+STAic6ASoiIi92MS9wb3N0aW5ncy97cG9zdGluZ19pZH0vY29tbWVudHMSfgoNVXBkYXRlQ29tbWVudBIeLmJvYXJkLnYxLlVwZGF0ZUNvbW1lbnRSZXF1ZXN0GhEuYm9hcmQudjEuQ29tbWVudCI6gtPkkwI0OgEqGi8vdjEvcG9zdGluZ3Mve3Bvc3RpbmdfaWR9L2NvbW1lbnRzL3tjb21tZW50X2lkfRKAAQoNRGVsZXRlQ29tbWVudBIeLmJvYXJkLnYxLkRlbGV0ZUNvbW1lbnRSZXF1ZXN0GhYuZ29vZ2xlLnByb3RvYnVmLkVtcHR5IjeC0+STAjEqLy92MS9wb3N0aW5ncy97cG9zdGluZ19pZH0vY29tbWVudHMve2NvbW1lbnRfaWR9EmgKDVVwZGF0ZVBvc3RpbmcSHi5ib2FyZC52MS5VcGRhdGVQb3N0aW5nUmVxdWVzdBoRLmJvYXJkLnYxLlBvc3RpbmciJILT5JMCHjoBKhoZL3YxL3Bvc3RpbmdzL3twb3N0aW5nX2lkfRJqCg1EZWxldGVQb3N0aW5nEh4uYm9hcmQudjEuRGVsZXRlUG9zdGluZ1JlcXVlc3QaFi5nb29nbGUucHJvdG9idWYuRW1wdHkiIYLT5JMCGyoZL3YxL3Bvc3RpbmdzL3twb3N0aW5nX2lkfRJuCg5TZWFyY2hQb3N0aW5ncxIfLmJvYXJkLnYxLlNlYXJjaFBvc3RpbmdzUmVxdWVzdBoeLmJvYXJkLnYxLkxpc3RQb3N0aW5nc1Jlc3BvbnNlIhuC0+STAhUSEy92MS9wb3N0aW5ncy9zZWFyY2hCpQEKDGNvbS5ib2FyZC52MUIMU2VydmljZVByb3RvUAFaRmdpdGh1Yi5jb20va2ltc2VvZ3l1L3BvcnRmb2xpby9iYWNrZW5kL2ludGVybmFsL3Byb3RvL2JvYXJkL3YxO2JvYXJkdjGiAgNCWFiqAghCb2FyZC5WMcoCCEJvYXJkXFYx4gIUQm9hcmRcVjFcR1BCTWV0YWRhdGHqAglCb2FyZDo6VjFiBnByb3RvMw", [file_google_protobuf_timestamp, file_google_protobuf_empty, file_google_api_annotations, file_gnostic_openapi_v3_annotations]);

/**
 * Posting은 게시판의 게시물 정보를 나타냅니다.
 *
 * @generated from message board.v1.Posting
 */
export type Posting = Message<"board.v1.Posting"> & {
  /**
   * 게시물의 고유 식별자
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;

  /**
   * 게시물 제목 (최소 1자, 최대 200자)
   *
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * 게시물 본문 내용 (최대 50000자)
   *
   * @generated from field: string content = 3;
   */
  content: string;

  /**
   * 작성자 ID (인증 시스템의 사용자 ID)
   *
   * @generated from field: string author_id = 4;
   */
  authorId: string;

  /**
   * 작성자 표시 이름
   *
   * @generated from field: string author_name = 5;
   */
  authorName: string;

  /**
   * 게시물 생성 시간 (서버 시간 기준)
   *
   * @generated from field: google.protobuf.Timestamp created_at = 6;
   */
  createdAt?: Timestamp;

  /**
   * 게시물 최종 수정 시간
   *
   * @generated from field: google.protobuf.Timestamp updated_at = 7;
   */
  updatedAt?: Timestamp;

  /**
   * 게시물 삭제 시간
   *
   * @generated from field: google.protobuf.Timestamp deleted_at = 8;
   */
  deletedAt?: Timestamp;

  /**
   * 게시물 조회수
   *
   * @generated from field: int32 view_count = 9;
   */
  viewCount: number;

  /**
   * 게시물에 달린 댓글 수
   *
   * @generated from field: int32 comment_count = 10;
   */
  commentCount: number;

  /**
   * 게시물에 적용된 태그 목록 (각 태그 최대 50자)
   *
   * @generated from field: repeated string tags = 11;
   */
  tags: string[];

  /**
   * 게시물의 현재 상태
   *
   * @generated from field: board.v1.PostingStatus status = 12;
   */
  status: PostingStatus;

  /**
   * 게시물의 댓글 목록
   *
   * @generated from field: repeated board.v1.Comment comments = 13;
   */
  comments: Comment[];
};

/**
 * Describes the message board.v1.Posting.
 * Use `create(PostingSchema)` to create a new message.
 */
export const PostingSchema: GenMessage<Posting> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 0);

/**
 * Comment는 게시물에 달린 댓글 정보를 나타냅니다.
 *
 * @generated from message board.v1.Comment
 */
export type Comment = Message<"board.v1.Comment"> & {
  /**
   * 댓글의 고유 식별자
   *
   * @generated from field: int64 comment_id = 1;
   */
  commentId: bigint;

  /**
   * 댓글이 달린 게시물의 ID
   *
   * @generated from field: int64 posting_id = 2;
   */
  postingId: bigint;

  /**
   * 댓글 내용 (최소 1자, 최대 5000자)
   *
   * @generated from field: string content = 3;
   */
  content: string;

  /**
   * 작성자 ID (인증 시스템의 사용자 ID)
   *
   * @generated from field: string author_id = 4;
   */
  authorId: string;

  /**
   * 작성자 표시 이름
   *
   * @generated from field: string author_name = 5;
   */
  authorName: string;

  /**
   * 댓글 생성 시간 (서버 시간 기준)
   *
   * @generated from field: google.protobuf.Timestamp created_at = 6;
   */
  createdAt?: Timestamp;

  /**
   * 댓글 최종 수정 시간
   *
   * @generated from field: google.protobuf.Timestamp updated_at = 7;
   */
  updatedAt?: Timestamp;

  /**
   * 댓글 삭제 시간
   *
   * @generated from field: google.protobuf.Timestamp deleted_at = 8;
   */
  deletedAt?: Timestamp;

  /**
   * 대댓글인 경우 부모 댓글의 ID (최상위 댓글인 경우 0)
   *
   * @generated from field: int64 parent_id = 9;
   */
  parentId: bigint;

  /**
   * 댓글의 현재 상태
   *
   * @generated from field: board.v1.CommentStatus status = 10;
   */
  status: CommentStatus;

  /**
   * 대댓글 목록
   *
   * @generated from field: repeated board.v1.Comment children = 11;
   */
  children: Comment[];
};

/**
 * Describes the message board.v1.Comment.
 * Use `create(CommentSchema)` to create a new message.
 */
export const CommentSchema: GenMessage<Comment> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 1);

/**
 * CreatePostingRequest는 새 게시물 생성 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.CreatePostingRequest
 */
export type CreatePostingRequest = Message<"board.v1.CreatePostingRequest"> & {
  /**
   * 게시물 제목 (필수, 최소 1자, 최대 200자)
   *
   * @generated from field: string title = 1;
   */
  title: string;

  /**
   * 게시물 본문 내용 (필수, 최대 50000자)
   *
   * @generated from field: string content = 2;
   */
  content: string;

  /**
   * 작성자 ID (필수, 인증 시스템의 사용자 ID)
   *
   * @generated from field: string author_id = 3;
   */
  authorId: string;

  /**
   * 게시물에 적용할 태그 목록 (선택, 각 태그 최대 50자, 최대 10개)
   *
   * @generated from field: repeated string tags = 4;
   */
  tags: string[];

  /**
   * 게시물 상태 (선택, 기본값: PUBLISHED)
   *
   * @generated from field: board.v1.PostingStatus status = 5;
   */
  status: PostingStatus;
};

/**
 * Describes the message board.v1.CreatePostingRequest.
 * Use `create(CreatePostingRequestSchema)` to create a new message.
 */
export const CreatePostingRequestSchema: GenMessage<CreatePostingRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 2);

/**
 * GetPostingRequest는 게시물 조회 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.GetPostingRequest
 */
export type GetPostingRequest = Message<"board.v1.GetPostingRequest"> & {
  /**
   * 조회할 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;
};

/**
 * Describes the message board.v1.GetPostingRequest.
 * Use `create(GetPostingRequestSchema)` to create a new message.
 */
export const GetPostingRequestSchema: GenMessage<GetPostingRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 3);

/**
 * ListPostingsRequest는 게시물 목록 조회 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.ListPostingsRequest
 */
export type ListPostingsRequest = Message<"board.v1.ListPostingsRequest"> & {
  /**
   * 한 페이지에 반환할 최대 게시물 수 (기본값: 20, 최대: 100)
   *
   * @generated from field: int32 page_size = 1;
   */
  pageSize: number;

  /**
   * 페이지네이션 토큰 (첫 페이지는 빈 문자열)
   *
   * @generated from field: string page_token = 2;
   */
  pageToken: string;

  /**
   * 정렬 기준 (기본값: CREATED_AT)
   *
   * @generated from field: board.v1.SortBy sort_by = 3;
   */
  sortBy: SortBy;

  /**
   * 정렬 방향 (기본값: DESC)
   *
   * @generated from field: board.v1.SortDirection sort_direction = 4;
   */
  sortDirection: SortDirection;

  /**
   * 태그로 필터링 (선택, 특정 태그가 있는 게시물만 표시)
   *
   * @generated from field: string filter_by_tag = 5;
   */
  filterByTag: string;
};

/**
 * Describes the message board.v1.ListPostingsRequest.
 * Use `create(ListPostingsRequestSchema)` to create a new message.
 */
export const ListPostingsRequestSchema: GenMessage<ListPostingsRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 4);

/**
 * ListPostingsResponse는 게시물 목록 조회 응답 정보를 포함합니다.
 *
 * @generated from message board.v1.ListPostingsResponse
 */
export type ListPostingsResponse = Message<"board.v1.ListPostingsResponse"> & {
  /**
   * 조회된 게시물 목록
   *
   * @generated from field: repeated board.v1.Posting postings = 1;
   */
  postings: Posting[];

  /**
   * 다음 페이지를 요청할 때 사용할 토큰 (마지막 페이지면 빈 문자열)
   *
   * @generated from field: string next_page_token = 2;
   */
  nextPageToken: string;

  /**
   * 필터링 조건을 만족하는 전체 게시물 수
   *
   * @generated from field: int32 total_count = 3;
   */
  totalCount: number;
};

/**
 * Describes the message board.v1.ListPostingsResponse.
 * Use `create(ListPostingsResponseSchema)` to create a new message.
 */
export const ListPostingsResponseSchema: GenMessage<ListPostingsResponse> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 5);

/**
 * UpdatePostingRequest는 게시물 수정 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.UpdatePostingRequest
 */
export type UpdatePostingRequest = Message<"board.v1.UpdatePostingRequest"> & {
  /**
   * 수정할 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;

  /**
   * 수정할 제목 (선택, 최소 1자, 최대 200자)
   *
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * 수정할 내용 (선택, 최대 50000자)
   *
   * @generated from field: string content = 3;
   */
  content: string;

  /**
   * 수정할 태그 목록 (선택, 각 태그 최대 50자, 최대 10개)
   *
   * @generated from field: repeated string tags = 4;
   */
  tags: string[];

  /**
   * 수정할 게시물 상태 (선택)
   *
   * @generated from field: board.v1.PostingStatus status = 5;
   */
  status: PostingStatus;
};

/**
 * Describes the message board.v1.UpdatePostingRequest.
 * Use `create(UpdatePostingRequestSchema)` to create a new message.
 */
export const UpdatePostingRequestSchema: GenMessage<UpdatePostingRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 6);

/**
 * DeletePostingRequest는 게시물 삭제 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.DeletePostingRequest
 */
export type DeletePostingRequest = Message<"board.v1.DeletePostingRequest"> & {
  /**
   * 삭제할 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;
};

/**
 * Describes the message board.v1.DeletePostingRequest.
 * Use `create(DeletePostingRequestSchema)` to create a new message.
 */
export const DeletePostingRequestSchema: GenMessage<DeletePostingRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 7);

/**
 * SearchPostingsRequest는 게시물 검색 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.SearchPostingsRequest
 */
export type SearchPostingsRequest = Message<"board.v1.SearchPostingsRequest"> & {
  /**
   * 기본 검색어
   *
   * @generated from field: string text = 1;
   */
  text: string;

  /**
   * 제목 내 검색
   *
   * @generated from field: string title_contains = 2;
   */
  titleContains: string;

  /**
   * 내용 내 검색
   *
   * @generated from field: string content_contains = 3;
   */
  contentContains: string;

  /**
   * 특정 태그 검색
   *
   * @generated from field: repeated string tags = 4;
   */
  tags: string[];

  /**
   * 작성자 ID 필터
   *
   * @generated from field: string author_id = 5;
   */
  authorId: string;

  /**
   * 작성 기간 필터
   *
   * @generated from field: google.protobuf.Timestamp created_after = 6;
   */
  createdAfter?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp created_before = 7;
   */
  createdBefore?: Timestamp;

  /**
   * 페이지네이션 정보
   *
   * @generated from field: int32 page_size = 8;
   */
  pageSize: number;

  /**
   * @generated from field: string page_token = 9;
   */
  pageToken: string;

  /**
   * @generated from field: board.v1.SortDirection sort_direction = 10;
   */
  sortDirection: SortDirection;
};

/**
 * Describes the message board.v1.SearchPostingsRequest.
 * Use `create(SearchPostingsRequestSchema)` to create a new message.
 */
export const SearchPostingsRequestSchema: GenMessage<SearchPostingsRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 8);

/**
 * CreateCommentRequest는 댓글 생성 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.CreateCommentRequest
 */
export type CreateCommentRequest = Message<"board.v1.CreateCommentRequest"> & {
  /**
   * 댓글을 작성할 게시물의 ID (필수)
   *
   * string에서 int64로 변경
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;

  /**
   * 댓글 내용 (필수, 최소 1자, 최대 5000자)
   *
   * @generated from field: string content = 2;
   */
  content: string;

  /**
   * 작성자 ID (필수, 인증 시스템의 사용자 ID)
   *
   * @generated from field: string author_id = 3;
   */
  authorId: string;

  /**
   * 대댓글인 경우 부모 댓글의 ID (최상위 댓글인 경우 0)
   *
   * string에서 int64로 변경
   *
   * @generated from field: int64 parent_id = 4;
   */
  parentId: bigint;
};

/**
 * Describes the message board.v1.CreateCommentRequest.
 * Use `create(CreateCommentRequestSchema)` to create a new message.
 */
export const CreateCommentRequestSchema: GenMessage<CreateCommentRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 9);

/**
 * ListCommentsByPostingRequest는 게시물별 댓글 목록 조회 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.ListCommentsByPostingRequest
 */
export type ListCommentsByPostingRequest = Message<"board.v1.ListCommentsByPostingRequest"> & {
  /**
   * 댓글을 조회할 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 1;
   */
  postingId: bigint;

  /**
   * 한 페이지에 반환할 최대 댓글 수 (기본값: 50, 최대: 200)
   *
   * @generated from field: int32 page_size = 2;
   */
  pageSize: number;

  /**
   * 페이지네이션 토큰 (첫 페이지는 빈 문자열)
   *
   * @generated from field: string page_token = 3;
   */
  pageToken: string;

  /**
   * 정렬 방향 (기본값: ASC - 오래된 댓글부터)
   *
   * @generated from field: board.v1.SortDirection sort_direction = 4;
   */
  sortDirection: SortDirection;
};

/**
 * Describes the message board.v1.ListCommentsByPostingRequest.
 * Use `create(ListCommentsByPostingRequestSchema)` to create a new message.
 */
export const ListCommentsByPostingRequestSchema: GenMessage<ListCommentsByPostingRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 10);

/**
 * ListCommentsResponse는 댓글 목록 조회 응답 정보를 포함합니다.
 *
 * @generated from message board.v1.ListCommentsResponse
 */
export type ListCommentsResponse = Message<"board.v1.ListCommentsResponse"> & {
  /**
   * 조회된 댓글 목록
   *
   * @generated from field: repeated board.v1.Comment comments = 1;
   */
  comments: Comment[];

  /**
   * 다음 페이지를 요청할 때 사용할 토큰 (마지막 페이지면 빈 문자열)
   *
   * @generated from field: string next_page_token = 2;
   */
  nextPageToken: string;

  /**
   * 해당 게시물의 전체 댓글 수
   *
   * @generated from field: int32 total_count = 3;
   */
  totalCount: number;
};

/**
 * Describes the message board.v1.ListCommentsResponse.
 * Use `create(ListCommentsResponseSchema)` to create a new message.
 */
export const ListCommentsResponseSchema: GenMessage<ListCommentsResponse> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 11);

/**
 * UpdateCommentRequest는 댓글 수정 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.UpdateCommentRequest
 */
export type UpdateCommentRequest = Message<"board.v1.UpdateCommentRequest"> & {
  /**
   * 수정할 댓글의 ID (필수)
   *
   * @generated from field: int64 comment_id = 1;
   */
  commentId: bigint;

  /**
   * 댓글이 속한 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 2;
   */
  postingId: bigint;

  /**
   * 수정할 댓글 내용 (필수, 최소 1자, 최대 5000자)
   *
   * @generated from field: string content = 3;
   */
  content: string;
};

/**
 * Describes the message board.v1.UpdateCommentRequest.
 * Use `create(UpdateCommentRequestSchema)` to create a new message.
 */
export const UpdateCommentRequestSchema: GenMessage<UpdateCommentRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 12);

/**
 * DeleteCommentRequest는 댓글 삭제 요청 정보를 포함합니다.
 *
 * @generated from message board.v1.DeleteCommentRequest
 */
export type DeleteCommentRequest = Message<"board.v1.DeleteCommentRequest"> & {
  /**
   * 삭제할 댓글의 ID (필수)
   *
   * @generated from field: int64 comment_id = 1;
   */
  commentId: bigint;

  /**
   * 댓글이 속한 게시물의 ID (필수)
   *
   * @generated from field: int64 posting_id = 2;
   */
  postingId: bigint;
};

/**
 * Describes the message board.v1.DeleteCommentRequest.
 * Use `create(DeleteCommentRequestSchema)` to create a new message.
 */
export const DeleteCommentRequestSchema: GenMessage<DeleteCommentRequest> = /*@__PURE__*/
  messageDesc(file_board_v1_service, 13);

/**
 * PostingStatus는 게시물의 현재 상태를 나타냅니다.
 *
 * @generated from enum board.v1.PostingStatus
 */
export enum PostingStatus {
  /**
   * 상태가 명시되지 않음(기본값)
   *
   * @generated from enum value: POSTING_STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 공개 상태의 게시물
   *
   * @generated from enum value: POSTING_STATUS_PUBLISHED = 1;
   */
  PUBLISHED = 1,

  /**
   * 작성 중인 초안 상태
   *
   * @generated from enum value: POSTING_STATUS_DRAFT = 2;
   */
  DRAFT = 2,

  /**
   * 삭제된 게시물(논리적 삭제)
   *
   * @generated from enum value: POSTING_STATUS_DELETED = 3;
   */
  DELETED = 3,
}

/**
 * Describes the enum board.v1.PostingStatus.
 */
export const PostingStatusSchema: GenEnum<PostingStatus> = /*@__PURE__*/
  enumDesc(file_board_v1_service, 0);

/**
 * CommentStatus는 댓글의 현재 상태를 나타냅니다.
 *
 * @generated from enum board.v1.CommentStatus
 */
export enum CommentStatus {
  /**
   * 상태가 명시되지 않음(기본값)
   *
   * @generated from enum value: COMMENT_STATUS_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 공개 상태의 댓글
   *
   * @generated from enum value: COMMENT_STATUS_PUBLISHED = 1;
   */
  PUBLISHED = 1,

  /**
   * 삭제된 댓글(논리적 삭제)
   *
   * @generated from enum value: COMMENT_STATUS_DELETED = 2;
   */
  DELETED = 2,
}

/**
 * Describes the enum board.v1.CommentStatus.
 */
export const CommentStatusSchema: GenEnum<CommentStatus> = /*@__PURE__*/
  enumDesc(file_board_v1_service, 1);

/**
 * SortBy는 게시물 목록 정렬 기준을 정의합니다.
 *
 * @generated from enum board.v1.SortBy
 */
export enum SortBy {
  /**
   * 정렬 기준이 명시되지 않음(기본값: CREATED_AT)
   *
   * @generated from enum value: SORT_BY_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 생성 시간 기준 정렬
   *
   * @generated from enum value: SORT_BY_CREATED_AT = 1;
   */
  CREATED_AT = 1,

  /**
   * 수정 시간 기준 정렬
   *
   * @generated from enum value: SORT_BY_UPDATED_AT = 2;
   */
  UPDATED_AT = 2,

  /**
   * 조회수 기준 정렬
   *
   * @generated from enum value: SORT_BY_VIEW_COUNT = 3;
   */
  VIEW_COUNT = 3,
}

/**
 * Describes the enum board.v1.SortBy.
 */
export const SortBySchema: GenEnum<SortBy> = /*@__PURE__*/
  enumDesc(file_board_v1_service, 2);

/**
 * SortDirection은 정렬 방향을 정의합니다.
 *
 * @generated from enum board.v1.SortDirection
 */
export enum SortDirection {
  /**
   * 정렬 방향이 명시되지 않음(기본값: DESC)
   *
   * @generated from enum value: SORT_DIRECTION_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * 오름차순 정렬 (과거→현재, 적은→많은)
   *
   * @generated from enum value: SORT_DIRECTION_ASC = 1;
   */
  ASC = 1,

  /**
   * 내림차순 정렬 (현재→과거, 많은→적은)
   *
   * @generated from enum value: SORT_DIRECTION_DESC = 2;
   */
  DESC = 2,
}

/**
 * Describes the enum board.v1.SortDirection.
 */
export const SortDirectionSchema: GenEnum<SortDirection> = /*@__PURE__*/
  enumDesc(file_board_v1_service, 3);

/**
 * BoardService는 게시물과 댓글의 생성, 조회, 수정, 삭제 기능을 제공합니다.
 * 이 서비스는 RESTful API와 gRPC 인터페이스를 모두 지원합니다.
 *
 * 인증 요구사항:
 * - 모든 API 호출은 gRPC 메타데이터에 인증 정보를 포함해야 합니다.
 * - 필수 메타데이터 키:
 *   * "authorization": Bearer 토큰 형식으로 "Bearer {token}" 형태로 제공
 *   * "x-api-key": API 키 (선택적, 서비스 계정 인증 시 사용)
 *
 * 인증 오류 코드:
 * - UNAUTHENTICATED: 인증 정보가 없거나 유효하지 않은 경우
 * - PERMISSION_DENIED: 인증은 됐지만 요청된 작업에 대한 권한이 없는 경우
 * - RESOURCE_EXHAUSTED: 속도 제한 초과 시
 *
 * @generated from service board.v1.BoardService
 */
export const BoardService: GenService<{
  /**
   * CreatePosting은 새로운 게시물을 생성합니다.
   * 성공 시 생성된 게시물의 전체 정보를 반환합니다.
   * 필수 필드: title, content, author_id
   *
   * 인증 요구사항:
   * - 유효한 사용자 토큰이 필요합니다.
   * - author_id는 토큰의 사용자 ID와 일치해야 합니다(관리자 제외).
   *
   * @generated from rpc board.v1.BoardService.CreatePosting
   */
  createPosting: {
    methodKind: "unary";
    input: typeof CreatePostingRequestSchema;
    output: typeof PostingSchema;
  },
  /**
   * GetPosting은 지정된 ID의 게시물을 조회합니다.
   * 존재하지 않는 ID를 요청하면 NOT_FOUND 오류가 반환됩니다.
   * 요청 시 조회수(view_count)가 자동으로 증가합니다.
   *
   * 인증 요구사항:
   * - 공개 게시물(PUBLISHED)은 인증 없이 접근 가능합니다.
   * - 초안(DRAFT)은 작성자나 관리자만 조회 가능합니다.
   * - 삭제된 게시물(DELETED)은 관리자만 조회 가능합니다.
   *
   * @generated from rpc board.v1.BoardService.GetPosting
   */
  getPosting: {
    methodKind: "unary";
    input: typeof GetPostingRequestSchema;
    output: typeof PostingSchema;
  },
  /**
   * ListPostings는 페이지네이션을 지원하는 게시물 목록 조회 API입니다.
   * 정렬 및 태그 기반 필터링을 지원합니다.
   * page_token을 사용한 커서 기반 페이지네이션을 구현합니다.
   *
   * 인증 요구사항:
   * - 인증 없이 공개 게시물만 조회 가능합니다.
   * - 인증된 사용자는 자신의 초안 게시물도 결과에 포함됩니다.
   * - 관리자는 모든 게시물을 조회할 수 있습니다.
   *
   * @generated from rpc board.v1.BoardService.ListPostings
   */
  listPostings: {
    methodKind: "unary";
    input: typeof ListPostingsRequestSchema;
    output: typeof ListPostingsResponseSchema;
  },
  /**
   * ListCommentsByPosting은 특정 게시물에 달린 모든 댓글을 조회합니다.
   * 페이지네이션을 지원하며, 정렬 방향을 지정할 수 있습니다.
   *
   * @generated from rpc board.v1.BoardService.ListCommentsByPosting
   */
  listCommentsByPosting: {
    methodKind: "unary";
    input: typeof ListCommentsByPostingRequestSchema;
    output: typeof ListCommentsResponseSchema;
  },
  /**
   * CreateComment는 특정 게시물에 새 댓글을 작성합니다.
   * parent_id 필드를 통해 대댓글 작성이 가능합니다.
   * 필수 필드: posting_id, content, author_id
   *
   * @generated from rpc board.v1.BoardService.CreateComment
   */
  createComment: {
    methodKind: "unary";
    input: typeof CreateCommentRequestSchema;
    output: typeof CommentSchema;
  },
  /**
   * UpdateComment는 기존 댓글의 내용을 수정합니다.
   * 댓글 작성자만 수정할 수 있으며, 권한이 없으면 PERMISSION_DENIED 오류가 발생합니다.
   * 필수 필드: comment_id, posting_id, content
   *
   * @generated from rpc board.v1.BoardService.UpdateComment
   */
  updateComment: {
    methodKind: "unary";
    input: typeof UpdateCommentRequestSchema;
    output: typeof CommentSchema;
  },
  /**
   * DeleteComment는 댓글을 삭제 상태로 변경합니다(실제 삭제는 아님).
   * 댓글 작성자나 게시물 작성자, 또는 관리자만 삭제할 수 있습니다.
   * 권한이 없으면 PERMISSION_DENIED 오류가 발생합니다.
   *
   * @generated from rpc board.v1.BoardService.DeleteComment
   */
  deleteComment: {
    methodKind: "unary";
    input: typeof DeleteCommentRequestSchema;
    output: typeof EmptySchema;
  },
  /**
   * UpdatePosting은 기존 게시물의 내용을 수정합니다.
   * 게시물 작성자만 수정할 수 있으며, 권한이 없으면 PERMISSION_DENIED 오류가 발생합니다.
   * 필수 필드: posting_id
   * 선택적으로 title, content, tags, status 필드를 수정할 수 있습니다.
   *
   * @generated from rpc board.v1.BoardService.UpdatePosting
   */
  updatePosting: {
    methodKind: "unary";
    input: typeof UpdatePostingRequestSchema;
    output: typeof PostingSchema;
  },
  /**
   * DeletePosting은 게시물을 삭제 상태로 변경합니다(실제 삭제는 아님).
   * 게시물 작성자나 관리자만 삭제할 수 있습니다.
   * 권한이 없으면 PERMISSION_DENIED 오류가 발생합니다.
   *
   * @generated from rpc board.v1.BoardService.DeletePosting
   */
  deletePosting: {
    methodKind: "unary";
    input: typeof DeletePostingRequestSchema;
    output: typeof EmptySchema;
  },
  /**
   * SearchPostings는 게시물 제목과 내용에서 검색어를 포함하는 게시물을 찾습니다.
   * 페이지네이션과 정렬 방향을 지원합니다.
   * query가 비어있으면 모든 게시물을 반환합니다.
   *
   * @generated from rpc board.v1.BoardService.SearchPostings
   */
  searchPostings: {
    methodKind: "unary";
    input: typeof SearchPostingsRequestSchema;
    output: typeof ListPostingsResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_board_v1_service, 0);

