package main

import (
	"github.com/jinzhu/gorm"
)

type Restaurant struct {
	Model
	ImageURL      string
	Name          string
	Description   string `json:"Descr"`
	WebSite       string `json:"Web"`
	Location      string `json:"Address""`
	ContactName   string `json:"-"`
	City          string
	Coordinates   string
	Email         string   `json:"-"`
	Phone         string   `json:"-"`
	Moderated     bool     `json:"-"`
	OwnerApproved bool     `json:"-"`
	Tags          []string `gorm:"-"`
}

func (r *Restaurant) BeforeCreate(scope *gorm.Scope) error {
	var tagsStr string
	for _, tag := range r.Tags {
		tagsStr += tag + ";"
	}
	if len(tagsStr) > 1 {
		tagsStr = tagsStr[:len(tagsStr)-1]
	}
	scope.SetColumn("Tags", tagsStr)
	return nil
}
