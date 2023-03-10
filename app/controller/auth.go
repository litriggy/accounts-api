package controller

import (
	"database/sql"
	"strconv"

	"accounts/api/app/model"
	"accounts/api/pkg/utils/social"
	"accounts/api/platform/memcached"
	"accounts/api/service"

	"github.com/gofiber/fiber/v2"
)

// SignIn method to create a new user.
//
//	@Description	소셜 로그인으로 로그인 혹은 회원가입을 진행합니다. 계정 생성이면 201, 로그인이면 200을 리턴합니다.
//	@Summary		로그인/회원가입
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			type		path		string			true	"소셜 로그인 종류 (ex: google, apple)"
//	@Param			version		path		string			true	"소셜 로그인 버전 (ex: 1, 2, ...)"
//	@Param			accessToken	body		model.SignIn	true	"소셜 로그인 측에서 제공한 access token"
//	@Success		200			{object}	[]string{}
//	@Router			/v1/auth/signin/{type}/{version} [post]
func SignIn(c *fiber.Ctx) error {
	params := c.AllParams()
	body := &model.SignIn{}
	var status int
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on retrieving access token",
			"data":   err.Error(),
		})
	}
	//validate google access token with google oauth
	//var userdata db.User
	version, _ := strconv.Atoi(params["version"])
	var userID string
	var authInfo *model.AuthInfo
	var content *fiber.Map
	switch params["type"] {
	case "google":
		authInfo, status, content = social.GoogleLogin(body.AccessToken, version)
		//|| info.Audience != "997490977948-j1krge4vtu5rn7llojk5s3m75eg3c3ol.apps.googleusercontent.com"
		if status != 0 {
			return c.Status(status).JSON(content)
		}
		userdata, err := service.FindUser(authInfo.ID, params["type"])
		if err != nil {
			if sql.ErrNoRows == err {
				//create Account
				newUserID, err := service.CreateUser(authInfo.ID, params["type"], int32(version), "", authInfo.Email, "USER")
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": "error",
						"msg":    "error on creating user",
						"data":   err.Error(),
					})
				}
				userID = newUserID
				status = fiber.StatusCreated
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": "error",
					"msg":    "error on finding user",
					"data":   err.Error(),
				})
			}
		} else {
			userID = strconv.Itoa(int(userdata.ID))
			status = fiber.StatusOK
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "invalid type value",
			"data":   params["type"],
		})
	}
	sessionKey, err := memcached.CreateSession(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on generating session key",
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"status":     "ok",
		"type":       params["type"],
		"sessionKey": sessionKey,
		"userID":     userID,
		"email":      authInfo.Email,
	})
}

func SignUp(oauthId string, oauthType string, name string, email string) {

}

// CheckSession func for check userID and revoke new session key from old-session key
//
//	@Description	세션키를 확인하여 유저 정보를 반환 및 세션키를 재생성합니다.
//	@Summary		세션키 확인
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//
//	@Param			Authorization	header	string	true	"액세스 토큰"
//
//	@Router			/v1/auth/check [get]
func CheckSession(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	newSessionKey := c.Locals("newSessionKey")
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"userID":        userID,
		"newSessionKey": newSessionKey,
	})
}
