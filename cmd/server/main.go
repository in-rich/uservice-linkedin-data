package main

import (
	"database/sql"
	"fmt"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/config"
	"github.com/in-rich/uservice-linkedin-data/migrations"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.App.Postgres.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())

	defer func() {
		_ = db.Close()
		_ = sqldb.Close()
	}()

	err := db.Ping()
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	profilePicturesBucket, err := config.StorageClient.Bucket(config.Firebase.Buckets.ProfilePictures)
	if err != nil {
		log.Fatalf("failed to connect to profile pictures bucket: %v", err)
	}
	companyLogosBucket, err := config.StorageClient.Bucket(config.Firebase.Buckets.CompanyLogos)
	if err != nil {
		log.Fatalf("failed to connect to company logos bucket: %v", err)
	}

	getUsersDAO := dao.NewGetUserRepository(db)
	listUsersDAO := dao.NewListUsersRepository(db)
	createUserDAO := dao.NewCreateUserRepository(db)
	updateUserDAO := dao.NewUpdateUserRepository(db)
	getUserLastUpdateDAO := dao.NewGetUserLastUpdateRepository(db)

	getCompaniesDAO := dao.NewGetCompanyRepository(db)
	listCompaniesDAO := dao.NewListCompaniesRepository(db)
	createCompanyDAO := dao.NewCreateCompanyRepository(db)
	updateCompanyDAO := dao.NewUpdateCompanyRepository(db)
	getCompanyLastUpdateDAO := dao.NewGetCompanyLastUpdateRepository(db)

	getProfilePicturesDAO := dao.NewGetProfilePictureRepository(profilePicturesBucket)
	upsertProfilePictureDAO := dao.NewUpsertProfilePictureRepository(profilePicturesBucket)

	getCompanyLogoDAO := dao.NewGetProfilePictureRepository(companyLogosBucket)
	upsertCompanyLogoDAO := dao.NewUpsertProfilePictureRepository(companyLogosBucket)

	getUserService := services.NewGetUserService(getUsersDAO, getProfilePicturesDAO)
	listUsersService := services.NewListUsersService(listUsersDAO, getProfilePicturesDAO)
	upsertUserService := services.NewUpsertUserService(createUserDAO, updateUserDAO, upsertProfilePictureDAO)
	getUserLastUpdateService := services.NewGetUserLastUpdateService(getUserLastUpdateDAO)

	getCompanyService := services.NewGetCompanyService(getCompaniesDAO, getCompanyLogoDAO)
	listCompaniesService := services.NewListCompaniesService(listCompaniesDAO, getCompanyLogoDAO)
	upsertCompanyService := services.NewUpsertCompanyService(createCompanyDAO, updateCompanyDAO, upsertCompanyLogoDAO)
	getCompanyLastUpdateService := services.NewGetCompanyLastUpdateService(getCompanyLastUpdateDAO)

	getUserHandler := handlers.NewGetUser(getUserService)
	listUsersHandler := handlers.NewListUsers(listUsersService)
	upsertUserHandler := handlers.NewUpsertUser(upsertUserService)
	getUserLastUpdateHandler := handlers.NewGetUserLastUpdate(getUserLastUpdateService)

	getCompanyHandler := handlers.NewGetCompany(getCompanyService)
	listCompaniesHandler := handlers.NewListCompanies(listCompaniesService)
	upsertCompanyHandler := handlers.NewUpsertCompany(upsertCompanyService)
	getCompanyLastUpdateHandler := handlers.NewGetCompanyLastUpdate(getCompanyLastUpdateService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.App.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	defer func() {
		server.GracefulStop()
		_ = listener.Close()
	}()

	linkedin_data_pb.RegisterGetUserServer(server, getUserHandler)
	linkedin_data_pb.RegisterListUsersServer(server, listUsersHandler)
	linkedin_data_pb.RegisterUpsertUserServer(server, upsertUserHandler)
	linkedin_data_pb.RegisterGetUserLastUpdateServer(server, getUserLastUpdateHandler)

	linkedin_data_pb.RegisterGetCompanyServer(server, getCompanyHandler)
	linkedin_data_pb.RegisterListCompaniesServer(server, listCompaniesHandler)
	linkedin_data_pb.RegisterUpsertCompanyServer(server, upsertCompanyHandler)
	linkedin_data_pb.RegisterGetCompanyLastUpdateServer(server, getCompanyLastUpdateHandler)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
