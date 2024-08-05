package models

type Results struct {
	Id             int `json:"id"`
	Vowels         int `json:"vowels"`
	Spaces         int `json:"spaces"`
	Capitalletters int `json:"capitalletters"`
	Smallletters   int `json:"smallleters"`
	Words          int `json:"words"`
}

type DBResults struct {
	Id             int `db:"id"`
	Vowels         int `db:"vowels"`
	Spaces         int `db:"spaces"`
	Capitalletters int `db:"capitalletters"`
	Smallletters   int `db:"smallleters"`
	Words          int `db:"words"`
}

type UpdateField struct {
	Field string `json:"field"`
	Value int    `json:"value"`
	Id    int    `json:"id" db:"id"`
}

type Identity struct {
	User     string `json: user`
	Email    string `json: email`
	Password string `json: password`
}
type Identify struct {
	User     string `json: user`
	Password string `json: password`
}
