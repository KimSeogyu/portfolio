import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import { BoardService } from "~/proto/board/v1/service_pb";

// 서버 사이드에서만 사용하는 클라이언트
let boardClient: ReturnType<typeof createClient<typeof BoardService>>;

// 싱글톤 패턴으로 서버 클라이언트 생성
export function getBoardClient() {
  if (!boardClient) {
    boardClient = createClient(
      BoardService,
      createGrpcTransport({
        baseUrl: process.env.GRPC_SERVER_URL || "http://localhost:8080",
      })
    );
  }
  return boardClient;
} 