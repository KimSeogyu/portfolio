package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kimseogyu/portfolio/backend/internal/docs"
	boardv1 "github.com/kimseogyu/portfolio/backend/internal/proto/grpcgateway/board/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcGatewayUtil struct {
	httpMux     *http.ServeMux
	exposeSpec  bool
	gatewayPort int
	grpcPort    int
}

func NewGrpcGatewayUtil(
	apiService *Service,
	exposeSpec bool,
	grpcPort int,
	gatewayPort int,
) (*GrpcGatewayUtil, error) {
	mux := runtime.NewServeMux()
	err := boardv1.RegisterBoardServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf(":%d", grpcPort), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		return nil, err
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	return &GrpcGatewayUtil{httpMux: httpMux, exposeSpec: exposeSpec, gatewayPort: gatewayPort, grpcPort: grpcPort}, nil
}
func (g *GrpcGatewayUtil) Start(ctx context.Context) error {
	serveFileFS := func(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
		fs := http.FileServer(http.FS(fsys))
		r.URL.Path = name
		fs.ServeHTTP(w, r)
	}

	if g.exposeSpec {
		g.httpMux.HandleFunc("/v1/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			serveFileFS(w, r, docs.StaticFiles, "openapi.yaml")
		})
		g.httpMux.HandleFunc("/v1/spec", func(w http.ResponseWriter, r *http.Request) {
			serveFileFS(w, r, docs.StaticFiles, "openapi.html")
		})
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("[::]:%d", g.gatewayPort), g.httpMux); err != nil {
			fmt.Printf("failed to serve gateway: %v", err)
		}
	}()
	return nil
}

func (g *GrpcGatewayUtil) Stop(ctx context.Context) error {
	return nil
}
