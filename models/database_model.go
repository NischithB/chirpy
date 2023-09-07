package models

type DatabaseModel struct {
	Chirps        map[int]Chirp     `json:"chirps"`
	Users         map[int]User      `json:"users"`
	RevokedTokens map[string]string `json:"revoked-tokens"`
}
