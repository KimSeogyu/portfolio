package internal

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/server"
	"github.com/kimseogyu/portfolio/backend/internal/config"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	"github.com/kimseogyu/portfolio/backend/internal/grpcutils"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/kimseogyu/portfolio/backend/pkg/authenticator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run the board server",
	Long:  `Run the board server.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		configPath := cmd.Flag("config").Value.String()
		if configPath == "" {
			zap.S().Fatalf("Config file path is required")
		}

		zap.S().Infof("Config file path: %s", configPath)

		config, err := config.NewConfigFromFile(configPath)
		if err != nil {
			zap.S().Fatalf("Failed to read config file: %v", err)
		}

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCConfig.GrpcPort))
		if err != nil {
			zap.S().Fatalf("Failed to listen: %v", err)
		}

		db, err := db.NewDB(
			db.WithDBType(db.DBType(config.DBConfig.DBType)),
			db.WithPostgresOptions(
				db.WithHost(config.DBConfig.DB.Host),
				db.WithPort(config.DBConfig.DB.Port),
				db.WithUser(config.DBConfig.DB.User),
				db.WithPassword(config.DBConfig.DB.Password),
				db.WithDBName(config.DBConfig.DB.DBName),
			),
		)
		if err != nil {
			zap.S().Fatalf("Failed to create db: %v", err)
		}

		db.AutoMigrate(&postings.Posting{}, &comments.Comment{})

		postingRepo := postings.NewRepository(db)
		commentRepo := comments.NewRepository(db)
		authenticator := authenticator.NewRealAuthenticator()

		service, err := server.NewService(
			server.WithCommentRepository(commentRepo),
			server.WithPostingRepository(postingRepo),
			server.WithAuthenticator(authenticator),
		)
		if err != nil {
			zap.S().Fatalf("Failed to create service: %v", err)
		}

		grpcServer, err := grpcutils.NewGrpcServer(
			grpcutils.WithListener(listener),
			grpcutils.WithServices(
				grpcutils.NewService(&boardServer.BoardService_ServiceDesc, service),
			),
		)
		if err != nil {
			zap.S().Fatalf("Failed to create grpc server: %v", err)
		}

		grpcGatewayUtil, err := server.NewGrpcGatewayUtil(service, true, config.GRPCConfig.GrpcPort, config.GRPCConfig.GatewayPort)
		if err != nil {
			zap.S().Fatalf("Failed to create grpc gateway util: %v", err)
		}
		go grpcGatewayUtil.Start(ctx)
		zap.S().Infof("Gateway started on port %d", config.GRPCConfig.GatewayPort)
		defer grpcGatewayUtil.Stop(ctx)

		go grpcServer.Start(ctx)
		defer grpcServer.Stop(ctx)
		zap.S().Infof("Server started on port %d", config.GRPCConfig.GrpcPort)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT)

		select {
		case <-ctx.Done():
			return
		case <-signalChan:
			zap.S().Info("Received signal, shutting down...")
			cancel()
		}

	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	runserverCmd.Flags().StringP("config", "c", ".config/config.yaml", "config file path")

}
