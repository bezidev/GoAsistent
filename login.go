package GoAsistent

import (
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	SupportedUserTypes []string `json:"supported_user_types"`
}

type LoginResponse struct {
	AccessToken struct {
		Token          string `json:"token"`
		ExpirationDate string `json:"expiration_date"`
	} `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID            int    `json:"id"`
		Language      string `json:"language"`
		Username      string `json:"username"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		FreshPassword any    `json:"freshPassword"`
	} `json:"user"`
	Redirect any `json:"redirect"`
}

func Login(username, password string, devMode bool, refreshTokenCallback func(username string, refreshToken string)) (Session, error) {
	client := req.C()
	if devMode {
		client.DevMode()
	}

	r := LoginRequest{
		Username:           username,
		Password:           password,
		SupportedUserTypes: []string{"child"},
	}
	res, err := client.R().SetBodyJsonMarshal(r).SetHeaders(MOBILE_HEADER).Post(fmt.Sprintf("%s/m/login", EASISTENT_URL))
	if err != nil {
		return nil, err
	}
	var response LoginResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return nil, err
	}
	parse, err := time.Parse("2006-01-02T15:04:05-0700", response.AccessToken.ExpirationDate)
	if err != nil {
		return nil, err
	}

	c := req.C()
	c.Headers = make(http.Header) // idfk, zakaj tega ne naredi req.C()
	for i, v := range WEB_HEADER {
		c.Headers.Set(i, v)
	}
	c.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", response.AccessToken.Token))
	c.Headers.Set("X-Child-Id", fmt.Sprint(response.User.ID))
	c.Cookies = res.Cookies()
	c.Cookies = append(c.Cookies, &http.Cookie{
		Name:  "easistent_cookie",
		Value: "zapri",
	})
	if devMode {
		c.DevMode()
	}

	return &sessionImpl{
		AuthToken:            response.AccessToken.Token,
		RefreshToken:         response.RefreshToken,
		ChildId:              fmt.Sprint(response.User.ID),
		TokenExpiration:      int(parse.Unix()),
		Username:             response.User.Username,
		Name:                 response.User.Name,
		DevMode:              devMode,
		Client:               c,
		RefreshTokenCallback: refreshTokenCallback,
	}, nil
}
