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
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			country VARCHAR(255),
			public_profile BOOLEAN DEFAULT FALSE,
			phone_number VARCHAR(20),
			image VARCHAR(255)
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
