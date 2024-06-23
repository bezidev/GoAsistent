package GoAsistent

func SessionFromAuthToken(authToken string, refreshToken string, userId string, expiry int, username string, name string) (Session, error) {
	return &sessionImpl{
		AuthToken:       authToken,
		RefreshToken:    refreshToken,
		ChildId:         userId,
		TokenExpiration: expiry,
		Username:        username,
		Name:            name,
	}, nil
}
