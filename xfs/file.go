package xfs

import (
	"time"
)

type File struct {
	ID   string    `gorm:"size:255;not null;primaryKey" json:"id"`
	Name string    `gorm:"not null;" json:"name"`
	Ext  string    `gorm:"not null;" json:"ext"`
	Size int64     `gorm:"not null;" json:"size"`
	Data []byte    `gorm:"not null" json:"-"`
	Time time.Time `gorm:"not null" json:"time,omitempty"`
}

type FileResult struct {
	File *File `json:"file"`
}

type FilesResult struct {
	Files []*File `json:"files"`
}
