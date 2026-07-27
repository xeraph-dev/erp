package tokens

import "time"

type Pair struct {
	AccessToken           string    `json:"access_token" xml:"access_token"`
	RefreshToken          string    `json:"refresh_token" xml:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"-" xml:"-"`
	RefreshTokenExpiresAt time.Time `json:"-" xml:"-"`
}
