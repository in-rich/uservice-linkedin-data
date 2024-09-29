package main

import (
	"context"
	"fmt"
	"github.com/in-rich/lib-go/deploy"
	"github.com/in-rich/lib-go/monitor"
	linkedin_data_pb "github.com/in-rich/proto/proto-go/linkedin-data"
	"github.com/in-rich/uservice-linkedin-data/config"
	"github.com/in-rich/uservice-linkedin-data/migrations"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/in-rich/uservice-linkedin-data/pkg/handlers"
	"github.com/in-rich/uservice-linkedin-data/pkg/services"
	"github.com/rs/zerolog"
	"os"
)

func getLogger() monitor.GRPCLogger {
	if deploy.IsReleaseEnv() {
		return monitor.NewGCPGRPCLogger(zerolog.New(os.Stdout), "uservice-linkedin-data")
	}

	return monitor.NewConsoleGRPCLogger()
}

func main() {
	logger := getLogger()

	logger.Info("Starting server")
	db, closeDB, err := deploy.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Fatal(err, "failed to connect to database")
	}
	defer closeDB()

	if err := migrations.Migrate(db); err != nil {
		logger.Fatal(err, "failed to migrate")
	}

	profilePicturesBucket, err := config.StorageClient.Bucket(config.Firebase.Buckets.ProfilePictures)
	if err != nil {
		logger.Fatal(err, "failed to connect to profile pictures bucket")
	}
	companyLogosBucket, err := config.StorageClient.Bucket(config.Firebase.Buckets.CompanyLogos)
	if err != nil {
		logger.Fatal(err, "failed to connect to company logos bucket")
	}

	depCheck := deploy.DepsCheck{
		Dependencies: func() map[string]error {
			_, errProfileBucket := profilePicturesBucket.Attrs(context.Background())
			_, errCompanyBucket := companyLogosBucket.Attrs(context.Background())

			return map[string]error{
				"Postgres":                     db.Ping(),
				"CloudBucket(ProfilePictures)": errProfileBucket,
				"CloudBucket(CompanyLogos)":    errCompanyBucket,
			}
		},
		Services: deploy.DepCheckServices{
			"GetUser":              {"Postgres", "CloudBucket(ProfilePictures)"},
			"ListUsers":            {"Postgres", "CloudBucket(ProfilePictures)"},
			"UpsertUser":           {"Postgres", "CloudBucket(ProfilePictures)"},
			"GetUserLastUpdate":    {"Postgres", "CloudBucket(ProfilePictures)"},
			"GetCompany":           {"Postgres", "CloudBucket(CompanyLogos)"},
			"ListCompanies":        {"Postgres", "CloudBucket(CompanyLogos)"},
			"UpsertCompany":        {"Postgres", "CloudBucket(CompanyLogos)"},
			"GetCompanyLastUpdate": {"Postgres", "CloudBucket(CompanyLogos)"},
		},
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

	getUserHandler := handlers.NewGetUser(getUserService, logger)
	listUsersHandler := handlers.NewListUsers(listUsersService, logger)
	upsertUserHandler := handlers.NewUpsertUser(upsertUserService, logger)
	getUserLastUpdateHandler := handlers.NewGetUserLastUpdate(getUserLastUpdateService, logger)

	getCompanyHandler := handlers.NewGetCompany(getCompanyService, logger)
	listCompaniesHandler := handlers.NewListCompanies(listCompaniesService, logger)
	upsertCompanyHandler := handlers.NewUpsertCompany(upsertCompanyService, logger)
	getCompanyLastUpdateHandler := handlers.NewGetCompanyLastUpdate(getCompanyLastUpdateService, logger)

	logger.Info(fmt.Sprintf("Starting to listen on port %v", config.App.Server.Port))
	listener, server, health := deploy.StartGRPCServer(logger, config.App.Server.Port, depCheck)
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

	logger.Info("Server started")
	if err := server.Serve(listener); err != nil {
		logger.Fatal(err, "failed to serve")
	}
}
