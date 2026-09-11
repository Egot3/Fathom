package models

import "github.com/uptrace/bun"

func RegisterAll(db *bun.DB) {
	db.RegisterModel((*GroupsUsers)(nil))
	db.RegisterModel((*UserGroupsTests)(nil))
	db.RegisterModel((*TestsQuizzes)(nil))
	db.RegisterModel((*Answer)(nil))
}
