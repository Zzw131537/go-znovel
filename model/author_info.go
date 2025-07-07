package model

import "gorm.io/gorm"

type AuthorInfo struct {
	gorm.Model

	UserId int64

	PenName string

	TelPhone string

	InviteCode string

	ChatAccount string

	Email string

	Status int64

	WorkDirection int64
}
