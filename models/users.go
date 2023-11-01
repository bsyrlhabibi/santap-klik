package models

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" form:"name"`
	Username  string         `json:"username" form:"username"`
	Password  string         `json:"password" form:"password"`
	Alamat    string         `json:"alamat" form:"alamat"`
	HP        string         `gorm:"type:varchar(13);uniqueIndex" json:"hp"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserUpdateModel struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Alamat   string `json:"alamat"`
	HP       string `json:"hp"`
	Password string `json:"password"`
}

type UserLoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db}
}

func (um *UserModel) LoginUser(username string, password string) *User {
	var data User
	if err := um.db.Where("username = ?", username).First(&data).Error; err != nil {
		logrus.Error("Model : Login data error,", err.Error())
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))
	if err != nil {
		logrus.Error("Model : Password does not match")
		return nil
	}

	return &data
}

func (um *UserModel) RegisterUser(name, username, alamat, hp, password string) error {

	existingUser := um.GetUserByUsername(um.db, username)
	if existingUser != nil {

		return errors.New("User with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Name:     name,
		Username: username,
		Alamat:   alamat,
		HP:       hp,
		Password: string(hashedPassword),
	}

	if err := um.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (um *UserModel) GetUserByUsername(db *gorm.DB, username string) *User {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func (u *User) ValidatePhoneNumber() error {

	if len(u.HP) < 10 || len(u.HP) > 13 {
		return errors.New("Phone number must be between 10 and 13 characters")
	}

	for _, char := range u.HP {
		if char < '0' || char > '9' {
			return errors.New("Phone number can only contain digits")
		}
	}

	return nil
}

func (um *UserModel) GetAllUsers() ([]User, error) {
	var users []User
	if err := um.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (um *UserModel) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := um.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) UpdateUser(userID uint, name, username, alamat, hp, password string) error {
	user, err := um.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Update user fields
	if name != "" {
		user.Name = name
	}
	if username != "" {
		user.Username = username
	}
	if alamat != "" {
		user.Alamat = alamat
	}
	if hp != "" {
		user.HP = hp
	}

	// Hash and update the password if it's not empty
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	if err := um.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) DeleteUser(userID uint) error {
	if err := um.db.Delete(&User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) GetPageUser(page, limit int) ([]User, error) {
	var userList []User
	offset := (page - 1) * limit
	if err := um.db.Offset(offset).Limit(limit).Find(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}
