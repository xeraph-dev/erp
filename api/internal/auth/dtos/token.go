package dtos

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" xml:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token" xml:"access_token"`
	RefreshToken string `json:"refresh_token" xml:"refresh_token"`
}
