package database

import (
	"fmt"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		utilities.GoDotEnvVariable("DB_HOST"),
		utilities.GoDotEnvVariable("DB_USER"),
		utilities.GoDotEnvVariable("DB_PASSWORD"),
		utilities.GoDotEnvVariable("DB_NAME"),
		utilities.GoDotEnvVariable("DB_PORT"),
		utilities.GoDotEnvVariable("SSL_MODE"),
		utilities.GoDotEnvVariable("TIME_ZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Connection to database could not be established")
	}
	log.Println("Database connection established successfull")
	DB = db

	cmdArgs := os.Args[1:]

	fresh := false

	migrate := false

	seed := false
	for _, cmdArgs := range cmdArgs {

		if cmdArgs == "--fresh" {
			fresh = true
		}

		if cmdArgs == "--migrate" {
			migrate = true
		}

		if cmdArgs == "--seed" {
			seed = true
		}

	}
	//Check application environment

	if utilities.GoDotEnvVariable("APP_ENV") == "local" {
		//Drop Tables
		if fresh {
			DB.Migrator().DropTable(

				models.User{},
				models.Gender{},
				models.Role{},
				models.Relationship{},
				models.AccountStatus{},
				models.County{},
				models.SubCounty{},
				models.Specialty{},
				models.Doctor{},
				models.Appointment{},
				models.AppointmentStatus{},
				models.Patient{},
				models.AppointmentPaymentMethod{},
				models.InsuranceProviderStatus{},
				models.InsuranceProvider{},
				models.Language{},
				models.PatientFeedback{},
				models.Donation{},
			)
		}

		//Run Auto Migration
		if migrate {
			result := db.AutoMigrate(

				models.User{},
				models.Gender{},
				models.Role{},
				models.Relationship{},
				models.AccountStatus{},
				models.County{},
				models.SubCounty{},
				models.Specialty{},
				models.Doctor{},
				models.Appointment{},
				models.AppointmentStatus{},
				models.Patient{},
				models.AppointmentPaymentMethod{},
				models.InsuranceProviderStatus{},
				models.InsuranceProvider{},
				models.Language{},
				models.PatientFeedback{},
				models.Donation{},
			)
			if result != nil {
				log.Print(result.Error())
			}
		}
		if seed {
			CreateGenderSeeder()
			CreateRelationshipSeeder()
			CreateRolesSeeder()
			CreateAccountStatusSeeder()
			CreateAppointmentStatusSeeder()
			CreateCountiesSeeder()
			CreateSubCountiesSeeder()
			CreateAppointmentPaymentMethodsSeeder()
			CreateLanguagesSeeder()
			CreateInsuranceProviderStatusSeeder()
			CreateInsuranceProvidersSeeder()

		}

	} else {
		log.Fatalf("The command is not allowed for application in production")
	}
}
