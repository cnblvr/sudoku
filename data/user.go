package data

type UserInfo struct {
	Username string `json:"-"`
	Name     string `json:"name"`
}
