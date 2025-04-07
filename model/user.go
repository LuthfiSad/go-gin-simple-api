package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()" db:"id"`
	Name      string       `gorm:"size:100;not null" db:"name"`
	Email     string       `gorm:"size:100;uniqueIndex;not null" db:"email"`
	Password  string       `gorm:"size:255;not null" db:"password"`
	Role      string       `gorm:"size:50;not null;default:user" db:"role"`
	CreatedAt sql.NullTime `gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt sql.NullTime `gorm:"autoUpdateTime" db:"updated_at"`
}
