package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

	type Games struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"123456789"`
		EncryptedID   string            `bson:"encrypted_id,omitempty" json:"encrypted_id,omitempty"`
		Name        string             `bson:"name,omitempty" json:"name,omitempty" example:"Valorant"`
		Rating      float64            `bson:"rating,omitempty" json:"rating,omitempty" example:"123.12"`
		Release     primitive.DateTime `bson:"release_date,omitempty" json:"release_date,omitempty" swaggertype:"string" example:"2024-09-01T00:00:00Z" format:"date-time"`
		Desc        string             `bson:"desc,omitempty" json:"desc,omitempty" example:"A tactical first-person shooter game developed by Riot Games"`
		Genre       []string           `bson:"genre,omitempty" json:"genre,omitempty" example:"Adventure Games,Shooter,Action"`
		Dev_name    Developer          `bson:"dev_name,omitempty" json:"dev_name,omitempty"`
		Game_banner string             `bson:"game_banner,omitempty" json:"game_banner,omitempty" example:"https://i.ibb.co.com/k1KdV7t/genshin-main-banner.png"`
		Preview     string             `bson:"preview,omitempty" json:"preview,omitempty" example:"https://www.youtube.com/watch?v=qqnEjmnitgc"`
		Link_games  string             `bson:"link_games,omitempty" json:"link_games,omitempty" example:"https://genshin.hoyoverse.com/id/"`
		Game_logo   string             `bson:"game_logo,omitempty" json:"game_logo,omitempty" example:"https://i.ibb.co.com/Z6xFZP6/genshin-logo.png"`
	}

	type Developer struct {
		ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"123456789"`
		Name string             `bson:"name,omitempty" json:"name,omitempty" example:"HoYoverse"`
		Bio  string             `bson:"dev_bio,omitempty" json:"bio,omitempty" example:"tech otaku save the world"`
	}

type Admin struct {
	ID        primitive.ObjectID 	`bson:"_id,omitempty" json:"_id,omitempty"`
    User_name string             	`bson:"user_name,omitempty" json:"user_name,omitempty"`
	Name	  string				`bson:"name,omitempty" json:"name,omitempty"`
    Password  string             	`bson:"password,omitempty" json:"password,omitempty"`
	UpdatedAt time.Time          	`bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Token struct{
	ID			string 				`bson:"_id,omitempty" json:"_id,omitempty"`
	Token		string				`bson:"token" json:"token,omitempty"`
	AdminID		string				`bson:"admin_id" json:"admin_id,omitempty"`
	CreatedAt	time.Time			`bson:"created_at" json:"created_at"` 
}

type ReqGames struct {
	Name        string             `bson:"name,omitempty" json:"name,omitempty" example:"Valorant"`
	Rating      float64            `bson:"rating,omitempty" json:"rating,omitempty" example:"123.12"`
	Desc        string             `bson:"desc,omitempty" json:"desc,omitempty" example:"A tactical first-person shooter game developed by Riot Games"`
	Genre       []string           `bson:"genre,omitempty" json:"genre,omitempty" example:"Adventure Games,Shooter,Action"`
	Dev_name    ReqDeveloper       `bson:"dev_name,omitempty" json:"dev_name,omitempty"`
	Game_banner string             `bson:"game_banner,omitempty" json:"game_banner,omitempty" example:"https://i.ibb.co.com/k1KdV7t/genshin-main-banner.png"`
	Preview     string             `bson:"preview,omitempty" json:"preview,omitempty" example:"https://www.youtube.com/watch?v=qqnEjmnitgc"`
	Link_games  string             `bson:"link_games,omitempty" json:"link_games,omitempty" example:"https://genshin.hoyoverse.com/id/"`
	Game_logo   string             `bson:"game_logo,omitempty" json:"game_logo,omitempty" example:"https://i.ibb.co.com/Z6xFZP6/genshin-logo.png"`
}

type ReqDeveloper struct {
	Name string             `bson:"name,omitempty" json:"name,omitempty" example:"HoYoverse"`
	Bio  string             `bson:"dev_bio,omitempty" json:"bio,omitempty" example:"tech otaku save the world"`
}

type ReqAdmin struct {
	User_name	string 		       `bson:"user_name,omitempty" json:"user_name,omitempty" example:"Timmy"`
	Password 	string			   `bson:"password,omitempty" json:"password,omitempty" example:"timmy123"`
}

type GoogleUser struct {
	ID            string `bson:"_id,omitempty" json:"_id,omitempty"`
	Email         string `bson:"email,omitempty" json:"email,omitempty"`
	VerifiedEmail bool   `bson:"verified_email" json:"verified_email"`
	Picture       string `bson:"picture" json:"picture"`
	Name          string `bson:"name" json:"name"`
}