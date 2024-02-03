package xfu

import (
	"time"
)

type File struct {
	ID        string    `gorm:"size:256;not null;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Size      int64     `gorm:"not null" json:"size"`
	Type      string    `gorm:"not null" json:"type"`
	Data      []byte    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"<-:create;not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type FileResult struct {
	File *File `json:"file"`
}

type FilesResult struct {
	Files []*File `json:"files"`
}
