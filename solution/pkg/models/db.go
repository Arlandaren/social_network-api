// package models

// import (
// 	"fmt"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func InitDB(conn_string string){
// 	con, err := gorm.Open(postgres.Open(conn_string), &gorm.Config{})
// 	if err != nil{
// 		fmt.Println(err.Error())
// 	}
// 	if err := con.AutoMigrate(&User{}, &Countries{}); err!=nil{
// 		fmt.Println(err)
// 	}
// 	DB = con
// 	defer func() {
// 		dbInstance, _ := DB.DB()
// 		_ = dbInstance.Close()
// 	}()
// }

// func DropAllTables() error {
// 	if err := DB.Migrator().DropTable(&User{}, &Countries{}); err != nil {
//         return err
//     }
// 	return nil
// }

//	func GetCountries(){
//		Countries =
//	}
package models

import (
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(cfg string) error {
	db, err := sqlx.Connect("postgres", cfg)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {

		return err
	}
	DB = db
	MigrateTables()

	return nil
}

func DropAllTables() error {
	queries := []string{
		"DROP TABLE IF EXISTS users CASCADE;",
		"DROP TABLE IF EXISTS countries CASCADE;",
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func MigrateTables() error {
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			countryCode TEXT,
			isPublic BOOLEAN DEFAULT FALSE,
			phone TEXT,
			image TEXT
	);
    `); err != nil {
		return err
	}

	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS countries (
			id SERIAL PRIMARY KEY,
			name TEXT,
			alpha2 TEXT,
			alpha3 TEXT,
			region TEXT
		);
    `); err != nil {
		return err
	}

	return nil
}

func GetAllCountries(region string) ([]CountryResponse, error) {
	var countries []CountryResponse
	if region == "" {
		if err := DB.Select(&countries, "SELECT name,alpha2,alpha3,region FROM countries"); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	} else {
		if err := DB.Select(&countries, fmt.Sprintf("SELECT name,alpha2,alpha3,region FROM countries WHERE LOWER(region) = '%s'", strings.ToLower(region))); err != nil {
			return nil, err
		}
	}

	return countries, nil
}
func GetCountryByid(alpha2 string) (*Countries, error) {
	var countries Countries
	if err := DB.Get(&countries, fmt.Sprintf("SELECT * FROM countries WHERE LOWER(alpha2) = '%s'", strings.ToLower(alpha2))); err != nil {
		return nil, err
	}
	return &countries, nil
}
func GetUser(username string) (*User, error) {
	var user User
	err := DB.Get(&user, fmt.Sprintf("SELECT * FROM users WHERE LOWER(login) = '%s'", strings.ToLower(username)))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CreateUser(username string, email string, password string, country string, is_public bool, phone_number string, image string) (*Profile,error) { 
    _, err := DB.Exec("INSERT INTO users (login, email, password, countryCode, isPublic, phone, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", username, email, password, country, is_public, phone_number, image)
    if err != nil {
        return nil,err
    }
	profile := Profile{
		Login: username,
		Email: email,
		CountryCode: country,
		IsPublic: is_public,
		Phone: phone_number,
	}
    return &profile, nil
}
func GetProfile(id uint)(*Profile, error){
	var profile Profile
	if err := DB.Get(&profile, fmt.Sprintf("SELECT login, email, countryCode, isPublic, phone FROM users WHERE id = '%d'", id)); err != nil {
		return nil, err
	}
	return &profile,nil
}