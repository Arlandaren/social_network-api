package models

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"solution/pkg/utils"
	"strconv"
	"strings"
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
	if err := MigrateTables(); err != nil {
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
            ADD CONSTRAINT unique_alpha2 UNIQUE (alpha2),
			ADD CONSTRAINT check_alpha2_format CHECK (alpha2 ~ '^[a-zA-Z]{2}$');
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
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS friendships (
			id SERIAL PRIMARY KEY,
			user_login TEXT NOT NULL,
			friend_login TEXT NOT NULL UNIQUE,
			added_at TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (user_login) REFERENCES users(login),
			FOREIGN KEY (friend_login) REFERENCES users(login)
	);
    `); err != nil {
		return err
	}
	_, err = DB.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	`)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			content TEXT NOT NULL,
			author TEXT NOT NULL REFERENCES users(login),
			tags TEXT[],
			created_at TIMESTAMP DEFAULT NOW(),
			likes_count BIGINT DEFAULT 0,
			dislikes_count BIGINT DEFAULT 0
		);
	`)
	if err != nil {
		return err
	}
	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS reactions (
			id SERIAL PRIMARY KEY,
			user_login TEXT NOT NULL UNIQUE REFERENCES users(login),
			post_id UUID NOT NULL REFERENCES posts(id),
			is_like bool 
		);
    `); err != nil {
		return err
	}

	return nil
}
func GetAllCountries(region string) ([]CountryResponse, error) {
	var countries []CountryResponse

	if region == "none" {
		if err := DB.Select(&countries, "SELECT name,alpha2,alpha3,region FROM countries ORDER BY alpha2 ASC"); err != nil {
			return nil, err
		}
	} else {
		_ = DB.Select(&countries, fmt.Sprintf("SELECT name,alpha2,alpha3,region FROM countries WHERE LOWER(region) = '%s' ORDER BY alpha2 ASC", strings.ToLower(region)))
		if countries == nil {
			return nil, errors.New("не найдено стран с таким кодом")
		}
	}
	return countries, nil
}
func GetCountryByid(alpha2 string) (*Countries, error) {
	var countries Countries
	// if err := DB.Get(&countries, fmt.Sprintf("SELECT * FROM countries WHERE LOWER(alpha2) = '%s'", strings.ToLower(alpha2))); err != nil {
	if err := DB.Get(&countries, fmt.Sprintf("SELECT * FROM countries WHERE alpha2 = '%s'", alpha2)); err != nil {
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
func CreateUser(username string, email string, password string, country string, is_public bool, phone_number string, image string) (map[string]interface{}, error) {
	if username == "" || email == "" || password == "" || country == "" {
		return nil, errors.New("неверный формат")
	}
	_, err := DB.Exec("INSERT INTO users (login, email, password, countryCode, isPublic, phone, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", username, email, password, country, is_public, phone_number, image)
	if err != nil {
		return nil, err
	}
	profile := map[string]interface{}{
		"login":       username,
		"email":       email,
		"countryCode": country,
		"isPublic":    is_public,
		"phone":       phone_number,
	}
	// profile := Profile{
	// 	Login     : username,
	// 	Email   : email,
	// 	CountryCode :country,
	// 	IsPublic : is_public,
	// 	Phone : phone_number,

	// }
	return profile, nil
}
func GetMyProfile(id uint) (*Profile, error) {
	var profile Profile
	err := DB.Get(&profile, fmt.Sprintf("SELECT login, email, countryCode, isPublic, phone, image FROM users WHERE id = '%d'", id))
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
func GetProfile(login string) (*Profile, error) {
	var profile Profile
	if err := DB.Get(&profile, fmt.Sprintf("SELECT login, email, countryCode, isPublic, phone, image FROM users WHERE login = '%s' AND isPublic = true", login)); err != nil {
		return nil, err
	}
	return &profile, nil
}
func UpdateProfile(userId uint, editParameters *EditParameters) error {
	query := "UPDATE users SET"
	updates := []string{}

	if editParameters.CountryCode != "" {
		updates = append(updates, fmt.Sprintf("countrycode = '%s'", editParameters.CountryCode))
	}
	if strconv.FormatBool(editParameters.IsPublic) != "" {
		updates = append(updates, fmt.Sprintf("ispublic = %t", editParameters.IsPublic))
	}
	if editParameters.Phone != "" {
		updates = append(updates, fmt.Sprintf("phone = '%s'", editParameters.Phone))
	}
	if editParameters.Image != "" {
		updates = append(updates, fmt.Sprintf("image = '%s'", editParameters.Image))
	}

	if len(updates) == 0 {
		return nil // Нет обновлений
	}

	query += " " + strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = %d", userId)

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func UpdatePassword(form *UpdatePasswordForm, id uint) error {
	var password string
	err := DB.Get(&password, fmt.Sprintf("SELECT password FROM users WHERE id = '%d'", id))
	if err != nil {
		return err
	}
	if !utils.CompareHashPassword(form.OldPasword, password) {
		return errors.New("пароль не совпадает")
	}
	newPassword, err := utils.GenerateHashPassword(form.NewPasword)
	if err != nil {
		return err
	}
	_, err = DB.Exec(fmt.Sprintf("UPDATE users SET password = '%s' WHERE id = '%d'", newPassword, id))
	if err != nil {
		return err
	}

	return nil
}
func DeactivateToken(token string) error {
	_, err := DB.Exec("INSERT INTO tokensblacklist (token) VALUES ($1)", token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func CheckBlackList(token string) (int, error) {
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM tokensblacklist WHERE token = $1", token)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func AddFriend(friendLogin string, login string) error {
	_, err := DB.Exec(fmt.Sprintf("INSERT INTO friendships (user_login,friend_login) VALUES ('%s','%s')", login, friendLogin))
	if err != nil {
		return err
	}
	return nil
}
func RemoveFriend(friendLogin string, login string) error {
	_, err := DB.Exec(fmt.Sprintf("DELETE FROM friendships WHERE friend_login = '%s' AND user_login = '%s'", friendLogin, login))
	if err != nil {
		return err
	}
	return nil
}
func GetFriendsList(login string, offset int, limit int) ([]Friend, error) {
	var friends []Friend
	if err := DB.Select(&friends, fmt.Sprintf("SELECT friend_login,added_at FROM friendships WHERE user_login = '%s' ORDER BY added_at DESC LIMIT %d OFFSET %d", login, limit, offset)); err != nil {
		return nil, err
	}
	return friends, nil
}
func CreatePost(post *PostRequest) (string, error) {
	query := "INSERT INTO posts (content, author, tags) VALUES ($1,$2,$3) RETURNING id"
	var id string
	err := DB.QueryRow(query, post.Content, post.Author, post.Tags).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func GetPostById(id string, viewerLogin string) (*Post, error) {
	var author string
	err := DB.Get(&author, "SELECT author FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	user, err := GetUser(author)
	if err != nil {
		return nil, err
	}
	if !user.PublicProfile {
		var isFriend bool
		if err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM friendships WHERE user_login = (SELECT author FROM posts WHERE id = $1) AND friend_login = $2)", id, viewerLogin).Scan(&isFriend); err != nil {
			return nil, err
		}

		if !isFriend {
			var isAuthor bool
			if err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1 AND author = $2)", id, viewerLogin).Scan(&isAuthor); err != nil {
				return nil, err
			}
			if !isAuthor {
				return nil, errors.New("пользователь не имеет доступа к данному посту")
			}
		}
	}
	var likesCount int
	err = DB.Get(&likesCount, "SELECT COUNT(*) FROM reactions WHERE post_id = $1 AND is_like = true", id)
	if err != nil {
		return nil,err
	}
	_, err = DB.Exec("UPDATE posts SET likes_count = $1 WHERE id = $2", likesCount, id)
	if err != nil {
		return nil,err
	}
	var dislikesCount int
	err = DB.Get(&dislikesCount, "SELECT COUNT(*) FROM reactions WHERE post_id = $1 AND is_like = false", id)
	if err != nil {
		return nil,err
	}
	_, err = DB.Exec("UPDATE posts SET dislikes_count = $1 WHERE id = $2", dislikesCount, id)
	if err != nil {
		return nil,err
	}
	var post Post
	if err := DB.Get(&post, "SELECT * FROM posts WHERE id = $1", id); err != nil {
		return nil, err
	}
	return &post, nil
}
func GetMyFeedList(login string, offset int, limit int) ([]Post, error) {
	var posts []Post
	if err := DB.Select(&posts, fmt.Sprintf("SELECT * FROM posts WHERE author = '%s' ORDER BY created_at DESC LIMIT %d OFFSET %d", login, limit, offset)); err != nil {
		return nil, err
	}
	return posts, nil
}
func GetFeedById(userLogin string, targetLogin string, offset int, limit int) ([]Post, error) {
	var posts []Post
	if userLogin == targetLogin {
		err := DB.Select(&posts, "SELECT * FROM posts WHERE author = $1", targetLogin)
		if err != nil {
			return nil, err
		}
		return posts, nil
	}
	user, err := GetUser(targetLogin)
	if err != nil {
		return nil, err
	}
	if !user.PublicProfile {
		var isFriend bool
		if err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM friendships WHERE user_login = $1 AND friend_login = $2)", targetLogin, userLogin).Scan(&isFriend); err != nil {
			return nil, err
		}

		if !isFriend {
			return nil, errors.New("пользователь не имеет доступа к данному посту")
		}
	}

	if err := DB.Select(&posts, fmt.Sprintf("SELECT * FROM posts WHERE author = '%s' ORDER BY created_at DESC LIMIT %d OFFSET %d", targetLogin, limit, offset)); err != nil {
		return nil, err
	}
	return posts, nil
}
func Like(userLogin string, post_id string) error {
	_,err := GetPostById(post_id, userLogin)
	if err != nil{
		return errors.New("пост не найден или к нему нет доступа")
	}
	_,err = DB.Exec("INSERT INTO reactions (user_login, post_id, is_like) VALUES ($1,$2,true)",userLogin,post_id)
	if err !=nil{
		_, err = DB.Exec("UPDATE reactions SET is_like = true WHERE post_id = $1 AND user_login = $2", post_id, userLogin)
		if err != nil {
			return err
		}
	}
	// var likesCount int
	// err = DB.Get(&likesCount, "SELECT COUNT(*) FROM reactions WHERE post_id = $1 AND is_like = true", post_id)
	// if err != nil {
	// 	return err
	// }
	// _, err = DB.Exec("UPDATE posts SET likes_count = $1 WHERE id = $2", likesCount, post_id)
	// if err != nil {
	// 	return err
	// }
	return nil
}
func Dislike(userLogin string, post_id string)error{
	_,err := GetPostById(post_id, userLogin)
	if err != nil{
		return errors.New("пост не найден или к нему нет доступа")
	}
	_,err = DB.Exec("INSERT INTO reactions (user_login, post_id, is_like) VALUES ($1,$2,false)",userLogin,post_id)
	if err !=nil{
		_, err = DB.Exec("UPDATE reactions SET is_like = false WHERE post_id = $1 AND user_login = $2", post_id,userLogin)
		if err != nil {
			return err
		}
	}
	
	// var dislikesCount int
	// err = DB.Get(&dislikesCount, "SELECT COUNT(*) FROM reactions WHERE post_id = $1 AND is_like = false", post_id)
	// if err != nil {
	// 	return err
	// }
	// _, err = DB.Exec("UPDATE posts SET dislikes_count = $1 WHERE id = $2", dislikesCount, post_id)
	// if err != nil {
	// 	return err
	// }
	return nil
}