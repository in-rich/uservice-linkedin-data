package main

import (
	"github.com/in-rich/lib-go/deploy"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/config"
	"github.com/in-rich/uservice-linkedin-data/migrations"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"log"
)

func main() {
	log.Println("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer closeDB()

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

	log.Println("Starting to listen on port", config.App.Server.Port)
	listener, server, health := deploy.StartGRPCServer(config.App.Server.Port)
	defer deploy.CloseGRPCServer(listener, server)
	go health()

	linkedin_data_pb.RegisterGetUserServer(server, getUserHandler)
	linkedin_data_pb.RegisterListUsersServer(server, listUsersHandler)
	linkedin_data_pb.RegisterUpsertUserServer(server, upsertUserHandler)
	linkedin_data_pb.RegisterGetUserLastUpdateServer(server, getUserLastUpdateHandler)

	linkedin_data_pb.RegisterGetCompanyServer(server, getCompanyHandler)
	linkedin_data_pb.RegisterListCompaniesServer(server, listCompaniesHandler)
	linkedin_data_pb.RegisterUpsertCompanyServer(server, upsertCompanyHandler)
	linkedin_data_pb.RegisterGetCompanyLastUpdateServer(server, getCompanyLastUpdateHandler)

	log.Println("Server started")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
