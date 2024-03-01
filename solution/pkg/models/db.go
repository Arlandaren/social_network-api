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
	"errors"
	"fmt"
	"solution/pkg/utils"
	"strconv"
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
	if err := MigrateTables(); err != nil{
		return err
	}

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
		CREATE TABLE IF NOT EXISTS countries (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			alpha2 VARCHAR(2),
			alpha3 VARCHAR(3),
			region TEXT
		);
    `); err != nil {
		return err
	}
	var constraintExists bool
    err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE table_name = 'countries' AND constraint_name = 'unique_alpha2')").Scan(&constraintExists)
    if err != nil {
        return err
    }

    if !constraintExists {
        if _, err := DB.Exec(`
            ALTER TABLE countries 
            ADD CONSTRAINT unique_alpha2 UNIQUE (alpha2)
        `); err != nil {
            return err
        }
    }
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(30) UNIQUE NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			countryCode VARCHAR(2) NOT NULL,
			CONSTRAINT fk_country FOREIGN KEY (countryCode) REFERENCES countries(alpha2),
			isPublic BOOLEAN DEFAULT true NOT NULL,
			phone VARCHAR(20) UNIQUE,
			image VARCHAR(200)
	);
    `); err != nil {
		return err
	}
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS tokensblacklist (
			id SERIAL PRIMARY KEY,
			token TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_tokensblacklist_token ON tokensblacklist (token);
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
func CreateUser(username string, email string, password string, country string, is_public bool, phone_number string, image string) (map[string]interface{},error) { 
	if username == "" || email == "" || password == "" || country == ""{
        return nil, errors.New("неверный формат")
    }
    _, err := DB.Exec("INSERT INTO users (login, email, password, countryCode, isPublic, phone, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", username, email, password, country, is_public, phone_number, image)
    if err != nil {
        return nil,err
    }
	profile := map[string]interface{}{
		"Login": username,
		"Email": email,
		"CountryCode": country,
		"IsPublic": is_public,
		"Phone": phone_number,
	}
    return profile, nil
}
func GetMyProfile(id uint)(*Profile, error){
	var profile Profile
	err := DB.Get(&profile, fmt.Sprintf("SELECT login, email, countryCode, isPublic, phone, image FROM users WHERE id = '%d'", id))
	if err != nil{
		return nil, err
	}
	return &profile,nil
}
func GetProfile(login string)(*Profile, error){
	var profile Profile
	if err := DB.Get(&profile, fmt.Sprintf("SELECT login, email, countryCode, isPublic, phone, image FROM users WHERE login = '%s' AND isPublic = true", login)); err != nil {
		return nil, err
	}
	return &profile,nil
}
func UpdateProfile(userId uint, editParameters *EditParameters)error{
	query := "UPDATE users SET"
	if editParameters.CountryCode != "" {
		query +=  fmt.Sprintf(" countrycode = '%s',", editParameters.CountryCode)
	}
	if strconv.FormatBool(editParameters.IsPublic) != "" {
		query += fmt.Sprintf("ispublic = %s,", strconv.FormatBool(editParameters.IsPublic))
	}
	if editParameters.Phone != "" {
		query += fmt.Sprintf("phone = '%s',", editParameters.Phone)
	}
	if editParameters.Image != "" {
		query += fmt.Sprintf("image = '%s',", editParameters.Image)
	}
	query = query[:len(query)-1]
	
	query += fmt.Sprintf(" WHERE id = '%d'", userId)
	_, err := DB.Exec(query)
	if err != nil{
		return err
	}
	return nil
}
func UpdatePassword(form *UpdatePasswordForm, id uint) error{
	var password string
	err:= DB.Get(&password,fmt.Sprintf("SELECT password FROM users WHERE id = '%d'",id))
	if err !=nil{
		return err
	}
	if !utils.CompareHashPassword(form.OldPasword,password){
		return errors.New("пароль не совпадает")
	}
	newPassword,err := utils.GenerateHashPassword(form.NewPasword)
	if err !=nil{
		return err
	}
	_,err = DB.Exec(fmt.Sprintf("UPDATE users SET password = '%s' WHERE id = '%d'",newPassword,id))
	if err != nil{
		return err
	}

	return nil
}
func DeactivateToken(token string) error{
	_,err := DB.Exec("INSERT INTO tokensblacklist (token) VALUES ($1)",token)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func CheckBlackList(token string)(int,error){
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM tokensblacklist WHERE token = $1", token)
	if err != nil {
		return 0,err
	}

	return count,nil
}