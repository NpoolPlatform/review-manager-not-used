package api

import (
	review "github.com/NpoolPlatform/message/npool/review/mgr/v2"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	review.UnimplementedManagerServer
}

func Register(server grpc.ServiceRegistrar) {
	review.RegisterManagerServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}
