package GoAsistent

import (
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

type RefreshSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshSessionResponse struct {
	AccessToken struct {
		Token          string `json:"token"`
		ExpirationDate string `json:"expiration_date"`
	} `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Redirect     any    `json:"redirect"`
}

func (s *sessionImpl) RefreshSession() error {
	client := req.C()
	if s.DevMode {
		client.DevMode()
	}
	r := RefreshSessionRequest{
		RefreshToken: s.RefreshToken,
	}
	res, err := client.R().SetBodyJsonMarshal(r).SetHeaders(MOBILE_HEADER).Post(fmt.Sprintf("%s/m/refresh_token", EASISTENT_URL))
	if err != nil {
		return err
	}
	var response RefreshSessionResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return err
	}
	s.AuthToken = response.AccessToken.Token
	s.RefreshToken = response.RefreshToken
	parse, err := time.Parse("2006-01-02T15:04:05-0700", response.AccessToken.ExpirationDate)
	if err != nil {
		return err
	}
	s.TokenExpiration = int(parse.Unix())
	s.Client.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", response.AccessToken.Token))
	s.Client.Cookies = res.Cookies()
	return nil
}

func (s *sessionImpl) RefreshWebSession() error {
	client := req.C()
	if s.DevMode {
		client.DevMode()
	}
	client.Cookies = s.Client.Cookies
	res, err := client.R().SetHeaders(WEB_HEADER).Get(fmt.Sprintf("%s/webapp", EASISTENT_URL))
	if err != nil {
		return err
	}
	s.Client.Cookies = res.Cookies()
	return nil
}
