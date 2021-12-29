package models

import (
	"errors"
	"fmt"
	"github.com/Creedowl/NiuwaBI/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	Password    string `json:"-"`
	Permissions uint   `json:"permissions"`
}

type WorkspaceInfo struct {
	WorkspaceName string   `json:"workspaceName"`
	ReportName    []string `json:"reportName"`
}
type UserStatistics struct {
	WorkspaceInfo []WorkspaceInfo `json:"workspaceInfo"`
}

const (
	Admin = 1
)

func GetUserByName(name string) (*User, error) {
	var user User
	err := database.GetDB().Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetAuthUser(username, password string) (*User, error) {
	var user User
	err := database.GetDB().Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}
	return &user, nil
}

func CreateUser(username, nickname, password string, permissions uint) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := User{
		Name:        username,
		Nickname:    nickname,
		Password:    string(hashed),
		Permissions: permissions,
	}
	err = database.GetDB().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) IsAdmin() bool {
	return u.Permissions&Admin == 1
}
