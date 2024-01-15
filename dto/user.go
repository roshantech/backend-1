package dto

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreateUpdate struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilepic`
	Active     bool   `json:"active`
}


type UserInfo struct {
	About    string  `json:"About"`
	Country  string  `json:"Country"`
	Email    string  `json:"Email"`
	PhoneNo  string  `json:"PhoneNo"`
	Company  string  `json:"Company"`
	School   string  `json:"School"`
}