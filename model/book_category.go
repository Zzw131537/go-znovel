package model

import "gorm.io/gorm"

type BookCategory struct {
	gorm.Model

	WorkDirection int64

	Name string

	Sort int64
}
