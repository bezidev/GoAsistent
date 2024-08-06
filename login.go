package GoAsistent

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var CaptchaRequired = errors.New("captcha required")

type EAsistentJWTClaims struct {
	jwt.MapClaims
	ConsumerKey  string `json:"consumerKey"`
	UserId       string `json:"userId"`
	UserType     string `json:"userType"`
	SchoolId     any    `json:"schoolId"`
	SessionId    string `json:"sessionId"`
	PasswordHash any    `json:"password_hash"`
	IssuedAt     string `json:"issuedAt"`
	AppName      string `json:"appName"`
	ExpiresAt    string `json:"exp"`
}

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

type WebLoginResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrFields []any  `json:"errfields"`
	Data      struct {
		ResetForm       bool   `json:"reset_form"`
		RequireCaptcha  bool   `json:"require_captcha"`
		PrijavaRedirect string `json:"prijava_redirect,omitempty"`
	} `json:"data"`
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

func WebLogin(username, password string, devMode bool, refreshTokenCallback func(username string, refreshToken string)) (Session, error) {
	client := req.C()
	if devMode {
		client.DevMode()
	}

	fd := map[string]string{
		"uporabnik": username,
		"geslo":     password,
		"pin":       "",
		"captcha":   "",
		"koda":      "",
	}
	res, err := client.R().SetFormData(fd).SetHeaders(WEB_HEADER).Post(fmt.Sprintf("%s/p/ajax_prijava", EASISTENT_URL))
	if err != nil {
		return nil, err
	}
	var response WebLoginResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return nil, err
	}

	if response.Data.RequireCaptcha {
		return nil, CaptchaRequired
	}

	client.Cookies = res.Cookies()

	result, err := client.R().Get(fmt.Sprintf("%s/webapp", EASISTENT_URL))
	if err != nil {
		return nil, err
	}

	responseBody := result.String()

	cidr := regexp.MustCompile(`<meta name="x-child-id" content="(?P<ChildId>.*)">`)
	submatches := cidr.FindStringSubmatch(responseBody)
	if len(submatches) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid child id: %s", responseBody))
	}
	childId := submatches[1]

	cidr = regexp.MustCompile(`<meta name="access-token" content="(?P<AccessToken>.*)">`)
	submatches = cidr.FindStringSubmatch(responseBody)
	if len(submatches) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid access token: %s", responseBody))
	}
	accessToken := strings.ReplaceAll(submatches[1], "Bearer ", "")

	cidr = regexp.MustCompile(`<meta name="refresh-token" content="(?P<RefreshToken>.*)">`)
	submatches = cidr.FindStringSubmatch(responseBody)
	if len(submatches) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid refresh token: %s", responseBody))
	}
	refreshToken := submatches[1]

	p := jwt.NewParser()
	token, _, err := p.ParseUnverified(accessToken, &EAsistentJWTClaims{})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*EAsistentJWTClaims)
	if claims == nil {
		return nil, errors.New(fmt.Sprintf("invalid access token (claims): %s", accessToken))
	}
	expires := claims.ExpiresAt
	expirationTime, err := time.Parse("2006-01-02T15:04:05-0700", expires)
	if err != nil {
		return nil, err
	}

	client.Headers = make(http.Header)
	for i, v := range WEB_HEADER {
		client.Headers.Set(i, v)
	}
	client.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client.Headers.Set("X-Child-Id", childId)

	return &sessionImpl{
		AuthToken:            accessToken,
		RefreshToken:         refreshToken,
		ChildId:              childId,
		TokenExpiration:      int(expirationTime.Unix()),
		Username:             username,
		Name:                 "",
		DevMode:              devMode,
		Client:               client,
		RefreshTokenCallback: refreshTokenCallback,
	}, nil
}
