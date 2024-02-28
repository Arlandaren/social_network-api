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
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sqlx.DB

func InitDB(conn_string string) error {
	db, err := sqlx.Connect("postgres", conn_string)
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
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			country TEXT,
			public_profile BOOLEAN DEFAULT FALSE,
			phone_number TEXT,
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

func GetAllCountries(region string) ([]Countries, error) {
	var countries []Countries
	if region == "" {
		if err := DB.Select(&countries, "SELECT * FROM countries"); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	} else {
		if err := DB.Select(&countries, fmt.Sprintf("SELECT * FROM countries WHERE LOWER(region) = '%s'", strings.ToLower(region))); err != nil {
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
	err := DB.Get(&user, fmt.Sprintf("SELECT * FROM users WHERE LOWER(username) = '%s'", strings.ToLower(username)))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CreateUser(username string, email string, password string, country string, is_public bool, phone_number string, image string) error { 
    _, err := DB.Exec("INSERT INTO users (username, email, password, created_at, country, public_profile, phone_number, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", strings.ToLower(username), email, password, time.Now(), country, is_public, phone_number, image)
    if err != nil {
        return err
    }
    return nil
}