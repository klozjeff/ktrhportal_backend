package database

import (
	"encoding/json"
	"fmt"
	"io"
	"ktrhportal/models"
	"log"
	"os"

	"github.com/gosimple/slug"
)

// CreateGenderSeeder creates Gender Seeder
func CreateGenderSeeder() {
	genderList := [3]string{"Male", "Female", "Other"}
	var genders []models.Gender
	if result := DB.Find(&genders); result.RowsAffected == 0 {
		for _, gender := range genderList {
			genders = append(genders, models.Gender{
				Title: gender,
				Slug:  slug.Make(gender),
			})
		}
		DB.Create(&genders)
	}

}

/*
* CreateRelationshipSeeder Seeded
 */
func CreateRelationshipSeeder() {
	relationshipList := [7]string{"Parent", "Guardian", "Spouse", "Daughter", "Son", "Sibling", "Relative"}
	var relationships []models.Relationship
	if result := DB.Find(&relationships); result.RowsAffected == 0 {
		for _, relationship := range relationshipList {
			relationships = append(relationships, models.Relationship{
				Title: relationship,
			})
		}
		DB.Create(&relationships)
	}
}

/*
* CreateRolesSeeder Seeded
 */
func CreateRolesSeeder() {
	roleList := [5]string{"Superuser", "Admin", "Employee"}
	var roles []models.Role
	if result := DB.Find(&roles); result.RowsAffected == 0 {
		for _, role := range roleList {
			roles = append(roles, models.Role{
				Title: role,
			})
		}
		DB.Create(&roles)
	}
}

// CreateAccountStatusSeeder creates statuses Seeder
func CreateAccountStatusSeeder() {
	statusList := [3]string{"Active", "Inactive", "Suspended"}
	var statuses []models.AccountStatus
	if result := DB.Find(&statuses); result.RowsAffected == 0 {
		for _, status := range statusList {
			statuses = append(statuses, models.AccountStatus{
				Title: status,
			})
		}
		DB.Create(&statuses)
	}

}

// CreateAppointmentStatusSeeder creates statuses Seeder
func CreateAppointmentStatusSeeder() {
	statusList := [4]string{"New", "Rescheduled", "Completed", "Cancelled"}
	var statuses []models.AppointmentStatus
	if result := DB.Find(&statuses); result.RowsAffected == 0 {
		for _, status := range statusList {
			statuses = append(statuses, models.AppointmentStatus{
				Title: status,
			})
		}
		DB.Create(&statuses)
	}

}

// CreateLanguagesSeeder creates statuses Seeder
func CreateLanguagesSeeder() {
	langsList := [2]string{"English", "Swahili"}
	var langs []models.Language
	if result := DB.Find(&langs); result.RowsAffected == 0 {
		for _, lang := range langsList {
			langs = append(langs, models.Language{
				Title: lang,
			})
		}
		DB.Create(&langs)
	}

}

// CreateAppointmentStatusSeeder creates Payment Methods Seeder
func CreateAppointmentPaymentMethodsSeeder() {
	methodsList := [2]string{"Cash", "Insurance"}
	var methods []models.AppointmentPaymentMethod
	if result := DB.Find(&methods); result.RowsAffected == 0 {
		for _, method := range methodsList {
			methods = append(methods, models.AppointmentPaymentMethod{
				Title: method,
			})
		}
		DB.Create(&methods)
	}

}

// CreateInsuranceProvidersSeeder creates Providers Seeder
func CreateInsuranceProvidersSeeder() {
	providersList := [6]string{"NHIF", "Jubilee Insurance", "CIC Group", "AAR Insurance", "Old Mutual", "Minet"}
	var providers []models.InsuranceProvider
	if result := DB.Find(&providers); result.RowsAffected == 0 {
		for _, provider := range providersList {
			providers = append(providers, models.InsuranceProvider{
				Title: provider,
			})
		}
		DB.Create(&providers)
	}

}

// CreateInsuranceProvidersSeeder creates Providers Seeder
func CreateInsuranceProviderStatusSeeder() {
	statusList := [4]string{"Active", "Suspended", "Inactive"}
	var statuses []models.InsuranceProviderStatus
	if result := DB.Find(&statuses); result.RowsAffected == 0 {
		for _, status := range statusList {
			statuses = append(statuses, models.InsuranceProviderStatus{
				Title: status,
			})
		}
		DB.Create(&statuses)
	}

}

// CreateInsuranceProvidersSeeder creates Providers Seeder
func CreateAppointmentTypeSeeder() {
	typeList := [3]string{"Facility visit", "Televisit", "Emergency"}
	var types []models.AppointmentType
	if result := DB.Find(&types); result.RowsAffected == 0 {
		for _, status := range typeList {
			types = append(types, models.AppointmentType{
				Title: status,
			})
		}
		DB.Create(&types)
	}

}

// CreateEncounterStatusSeeder creates Encounter Status Seeder
func CreateEncounterStatusSeeder() {
	statusList := [4]string{"New", "Inprogress", "Onhold", "Completed"}
	var statuses []models.EncounterStatus
	if result := DB.Find(&statuses); result.RowsAffected == 0 {
		for _, status := range statusList {
			statuses = append(statuses, models.EncounterStatus{
				Title: status,
			})
		}
		DB.Create(&statuses)
	}
}

// CreateCountrySeeder countries from json
func CreateCountrySeeder() {
	// Open our jsonFile
	jsonFile, err := os.Open("data/countries.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened countries.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var countries []models.Country
	if jsonError := json.Unmarshal(byteValue, &countries); jsonError != nil {
		log.Printf("Error Unmarshaling Json: %s", jsonError.Error())
	}

	//insert countries into database
	var countriesIndb []models.Country
	if result := DB.Find(&countriesIndb); result.RowsAffected == 0 {
		DB.Create(&countries)
	}

}

/*
* CreateCountiesSeeder Seeded
 */
func CreateCountiesSeeder() {
	// Open our jsonFile
	jsonFile, err := os.Open("data/kenyan_counties.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened kenya-counties.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var counties []models.County
	if jsonError := json.Unmarshal(byteValue, &counties); jsonError != nil {
		log.Printf("Error Unmarshaling Json: %s", jsonError.Error())
	}
	//insert states into database
	DB.Create(&counties)
}

/*
* CreateSubCountiesSeeder Seeded
 */
func CreateSubCountiesSeeder() {
	// Open our jsonFile
	jsonFile, err := os.Open("data/kenyan_subcounties.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened kenya-subcounties.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var subcounties []models.SubCounty
	if jsonError := json.Unmarshal(byteValue, &subcounties); jsonError != nil {
		log.Printf("Error Unmarshaling Json: %s", jsonError.Error())
	}
	//insert states into database
	DB.Create(&subcounties)
}
