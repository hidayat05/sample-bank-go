package startup

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	"sample-bank/app/service"
	"sample-bank/config"
	"sample-bank/migration"
	pb "sample-bank/proto"
)

type App struct {
	Server *grpc.Server
	DB     *gorm.DB
}

func (a *App) Initialize(dbConfig *config.Config) *grpc.Server {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		dbConfig.DB.Username,
		dbConfig.DB.Password,
		dbConfig.DB.Host,
		dbConfig.DB.Port,
		dbConfig.DB.Name,
		dbConfig.DB.Charset)

	db, err := gorm.Open(dbConfig.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("could not connect database", err)
	} else {
		fmt.Printf("Database connected successfully\n")
	}

	a.DB = migration.DBMigrate(db)
	a.Server = grpc.NewServer()

	// register service handler here
	pb.RegisterBankServiceServer(a.Server, &service.BankService{DB: a.DB})
	return a.Server
}

func (a *App) Run(server *grpc.Server, port string) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	errServe := server.Serve(listen)
	if errServe != nil {
		log.Fatalf("failed to serve: %v", errServe)
	}
}
