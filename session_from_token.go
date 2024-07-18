package GoAsistent

import (
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
)

func SessionFromAuthToken(authToken string, refreshToken string, userId string, expiry int, username string, name string, devMode bool, refreshTokenCallback func(username string, refreshToken string)) (Session, error) {
	client := req.C()
	client.Headers = make(http.Header)
	for i, v := range WEB_HEADER {
		client.Headers.Set(i, v)
	}
	client.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	client.Headers.Set("X-Child-Id", userId)
	if devMode {
		client.DevMode()
	}

	return &sessionImpl{
		AuthToken:            authToken,
		RefreshToken:         refreshToken,
		ChildId:              userId,
		TokenExpiration:      expiry,
		Username:             username,
		Name:                 name,
		DevMode:              devMode,
		Client:               client,
		RefreshTokenCallback: refreshTokenCallback,
	}, nil
}
