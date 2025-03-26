import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node";
import { BoardService } from "./board/v1/service_pb";

export function createBoardGrpcClient(baseUrl: string) {
  return createClient(
    BoardService,
    createGrpcTransport({
      baseUrl,
    })
  )
}
