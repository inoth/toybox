package models

type UserInfo struct {
	Rid  string   `json:"rid"`
	Id   string   `json:"id"`
	Name string   `json:"name"`
	Icon string   `json:"icon"`
	Tags []string `json:"tags"`
}
