package common

type Account struct {
	Id          int    `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	Coin        int    `json:"coin" db:"coin"`
	CreatedDate string `json:"created_date" db:"created_date"`
	UpdateDate  string `json:"updated_date" db:"updated_date"`
}
