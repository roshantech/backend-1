package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string      `json:"username"`
	Email       string      `gorm:"not null;unique" json:"email"`
	Password    string      `gorm:"not null" json:"-"`
	Education   string      `json:"education"`
	ProfilePic  string      `json:"profilepic`
	Profession  string      `json:"profession"`
	Address     Address     `json:"address" gorm:"foreignKey:UserID"`
	Followers   []Follower  `json:"followers" gorm:"foreignKey:UserID"`
	Following   []Following `json:"followers" gorm:"foreignKey:UserID"`
	Designation string      `json:"designation" gorm:"foreignKey:UserID"`
	Role        string      `json:"role"`
	Active      bool        `json:"active"`
	PhoneNumber string      `json:"phone_number"`
	Name        string      `json:"name"`
	About       string      `json:"about"`
	IsVerified  bool        `json:"is_verified"`
	Company     string      `json:"company"`
	CreatedAt   string      `json:"createdat`
	UpdatedAt   string      `json:"updatedat`
}

type Address struct {
	UserID     uint   `gorm:"uniqueIndex"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type Follower struct {
	gorm.Model
	UserID     uint `json:"user_id"`
	FollowerID uint `json:"follower_id"`
}

type Following struct {
	gorm.Model
	UserID     uint `json:"user_id"`
	FollowerID uint `json:"follower_id"`
}

type Post struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Caption   string `json:"caption"`
	MediaURL  string `json:"media_url"`
	MediaType string `json:"media_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt   string `json:"updatedat`
	Comments  []Comment
	Likes     []Like
}

type Comment struct {
	gorm.Model
	PostID    uint   `json:"post_id"`
	UserID    uint   `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Like struct {
	gorm.Model
	PostID uint `json:"post_id"`
	UserID uint `json:"user_id"`
}

func MigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Address{}, &Follower{}, &Post{}, &Comment{}, &Like{})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
