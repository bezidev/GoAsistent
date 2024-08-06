package GoAsistent

import (
	"fmt"
	"time"
)

type MealsResponse struct {
	Items []any `json:"items"`
}

type MenusResponse struct {
	Items []any `json:"items"`
}

type OrderCreateRequest struct {
	Type string `json:"type"`
	Menu string `json:"menu"`
	Date string `json:"date"`
}

func (s *sessionImpl) GetMeals() (MealsResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return MealsResponse{}, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/meals", EASISTENT_URL))
	if err != nil {
		return MealsResponse{}, err
	}
	var response MealsResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return MealsResponse{}, err
	}
	return response, nil
}

func (s *sessionImpl) GetMenus() (MenusResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return MenusResponse{}, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/meals/menus", EASISTENT_URL))
	if err != nil {
		return MenusResponse{}, err
	}
	var response MenusResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return MenusResponse{}, err
	}
	return response, nil
}

func (s *sessionImpl) CreateOrder(date string, menu string) error {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	r := OrderCreateRequest{
		Type: "",
		Menu: menu,
		Date: date,
	}

	_, err := s.Client.R().SetBodyJsonMarshal(r).Post(fmt.Sprintf("%s/m/meals/meal", EASISTENT_URL))
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionImpl) DeleteOrder(mealId string) error {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	_, err := s.Client.R().Delete(fmt.Sprintf("%s/m/meals/meal/%s", EASISTENT_URL, mealId))
	if err != nil {
		return err
	}
	return nil
}
