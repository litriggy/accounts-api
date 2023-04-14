package social

import (
	"accounts/api/app/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

func GoogleLogin(token string, version int) (model.AuthInfo, int, *fiber.Map) {
	fmt.Println(token)
	var authInfo model.AuthInfo
	switch version {
	case 1:
		info, err := verifyIdToken(token)
		if err != nil {
			return authInfo, fiber.StatusBadRequest, &fiber.Map{"status": "error", "msg": "error on versioning", "data": err.Error()}
		}
		authInfo.ID = info.UserId
		authInfo.Email = info.Email
		authInfo.Picture = ""
		authInfo.Nickname = ""
	case 2:
		info, err := verifyAccessToken(token)
		fmt.Println(info)
		if err != nil {
			if err != nil {
				return authInfo, fiber.StatusBadRequest, &fiber.Map{"status": "error", "msg": "error on versioning", "data": err.Error()}
			}
		}
		authInfo.ID = info.ID
		authInfo.Email = info.Email
		authInfo.Picture = info.Picture
		authInfo.Nickname = info.Name
	default:
		return authInfo, fiber.StatusBadRequest, &fiber.Map{"status": "error", "msg": "error on versioning"}
	}
	return authInfo, 0, nil
}

// IdToken from google auth website
func verifyIdToken(idToken string) (*oauth2.Tokeninfo, error) {
	fmt.Println("asdf", idToken)
	oauth2Service, err := oauth2.New(httpClient)
	//oauth2.NewService()
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

// IdToken from chrome extension
func verifyAccessToken(accessToken string) (*model.GoogleUserInfo, error) {
	var userInfo *model.GoogleUserInfo
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		return nil, err
	}
	cleanToken := strings.TrimSpace(accessToken)

	req.Header.Add("Authorization", "Bearer "+cleanToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bytes)
	err = json.Unmarshal([]byte(body), &userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.ID == "" {
		return nil, errors.New("incorrect userinfo")
	}
	fmt.Println("user info info ", userInfo)

	return userInfo, nil

}
