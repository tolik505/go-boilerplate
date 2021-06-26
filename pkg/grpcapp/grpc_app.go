package grpcapp

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"goboilerplate/pkg/grpcapp/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App struct is instance of grpc application
type App struct {
	server *grpc.Server
}

// Initialize instantiates http application with configured routes
func Initialize(
	postsServer pb.PostsServiceServer,
) *App {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpcrecovery.UnaryServerInterceptor(),
			),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	pb.RegisterPostsServiceServer(server, postsServer)
	reflection.Register(server)

	return &App{
		server: server,
	}
}

// Run starts grpc server
func (a *App) Run(lis net.Listener) {
	if err := a.server.Serve(lis); err != nil {
		err = errors.Wrap(err, "failed to serve grpc app")
		log.Fatal(err)
	}
}

//WaitForInterrupt waits for interrupt signal to gracefully shutdown the server.
func (a *App) WaitForInterrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gRPC server")

	a.server.GracefulStop()
}
