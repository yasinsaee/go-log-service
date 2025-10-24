package app

import (
	"log"
	"net"

	"github.com/yasinsaee/go-log-service/internal/app/config"
	"google.golang.org/grpc"
)

func StartGRPCServer() {
	port := config.GetEnv("PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// s := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor()))
	s := grpc.NewServer()

	//repos
	// permissionRepo := repository_permission.NewMongoPermissionRepository(mongo.DB.Database, "permission")
	// roleRepo := repository_role.NewMongoRoleRepository(mongo.DB.Database, "role")
	// userRepo := repository_user.NewMongoUserRepository(mongo.DB.Database, "user")

	// //services
	// permissionService := permission.NewPermissionService(permissionRepo)
	// roleService := role.NewRoleService(roleRepo)
	// userService := user.NewUserService(userRepo)

	// //handlers
	// permissionHandler := permissiongrpc.New(permissionService)
	// roleHandler := rolegrpc.New(roleService, permissionService)
	// userHandler := usergrpc.New(userService, roleService, permissionService)

	// //register grpc services
	// permissionpb.RegisterPermissionServiceServer(s, permissionHandler)
	// rolepb.RegisterRoleServiceServer(s, roleHandler)
	// userpb.RegisterUserServiceServer(s, userHandler)

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
