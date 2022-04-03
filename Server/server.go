package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	// "golang.org/x/crypto/bcrypt"
	pb "server/protos/user"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	port   = flag.Int("port", 50051, "The server port")
	db_url = "postgresql://postgres:postgres@localhost:5432/postgres"
)

type userAuthServer struct {
	pb.UnimplementedUserAuthServer
}

func (s *userAuthServer) Register(ctx context.Context, in *pb.RegReq) (*pb.ApiRes, error) {
	// if username is not in database, and everything is alright, register, and store in database
	// else, return error message to client that username is already in database or something is wrong
	log.Printf("Registered: %v", in.GetUsername())
	return &pb.ApiRes{ResCode: 200, Message: "success"}, nil
}

func (s *userAuthServer) Login(ctx context.Context, in *pb.LoginReq) (*pb.ApiRes, error) {
	// if username is in database, and password is correct, login, and return success message to client
	// else, return error message to client that username is not in database or password is wrong
	log.Printf("Logged in: %v", in.GetUsername())
	return &pb.ApiRes{ResCode: 200, Message: "success"}, nil
}

func (s *userAuthServer) Logout(ctx context.Context, in *pb.LogoutReq) (*pb.ApiRes, error) {
	return &pb.ApiRes{ResCode: 200, Message: "successfully logged out"}, nil
}

func verifyRegReq(in *pb.RegReq) error {
	if in.GetUsername() == "" {
		return fmt.Errorf("username cannot be empty")
	}
	// is username already in database? if yes, return error
	if len(in.GetPassword()) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

// func verifyLoginReq(in *pb.LoginReq) error {
// 	if in.GetUsername() == "" {
// 		return fmt.Errorf("username is empty")
// 	}
// 	if in.GetPassword() == "" {
// 		return fmt.Errorf("password is empty")
// 	}
// 	return nil
// }

// type user struct{
// 	email string
// 	username string
// 	password string
// }

func initDbConn() *sql.DB {
	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("could not open postgresql connection: %v", err)
	}
	log.Printf("Connected to database at port 5432")
	return conn
}

func main() {
	// establish database connection
	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("could not open postgresql connection: %v", err)
	}
	log.Printf("Connected to database at port 5432")
	defer conn.Close()

	// Query database
	data, err := conn.Query("SELECT * FROM user")
	if err != nil {
		log.Fatalf("could not query database: %v", err)
	}
	defer data.Close()

	// usual code
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserAuthServer(s, &userAuthServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
