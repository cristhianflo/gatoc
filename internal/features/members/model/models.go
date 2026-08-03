package model

import "gorm.io/gorm"

type WelcomeRole struct {
	gorm.Model
	GuildID string
	RoleID  string
	UserID  *string
}

type ResponseMessage struct {
	gorm.Model
	GuildID  string
	Message  string
	Response string
	UserID   *string
}
