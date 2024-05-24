package main

import (
	"log/slog"
	"myHabr/internal/config"
	genUsers "myHabr/internal/users/delivery/grpc/gen"
	userR "myHabr/internal/users/repository"
	userUc "myHabr/internal/users/usecase"
	"strconv"

	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	grpcUsers "myHabr/internal/users/delivery/grpc"

	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	cfg := config.MustLoad()
	_ = godotenv.Load()

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	if err = db.Ping(); err != nil {
		slog.Error("fail ping postgres")
		err = fmt.Errorf("error happened in db.Ping: %w", err)
		slog.Error(err.Error())
	}

	UserRepo := userR.NewUserRepo(db)
	UserUsecase := userUc.NewUserUsecase(UserRepo)
	UserHandler := grpcUsers.NewUsersServerHandler(UserUsecase)
	gRPCServer := grpc.NewServer()
	genUsers.RegisterUserServer(gRPCServer, UserHandler)

	go func() {
		slog.Info(fmt.Sprintf("Start server on %s\n", ":"+strconv.FormatInt(int64(cfg.GRPC.UserPort), 10)))
		listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(cfg.GRPC.UserPort), 10))
		if err != nil {
			log.Fatal(err)
		}
		if err := gRPCServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	gRPCServer.GracefulStop()
	return nil
}
