package controller

import (
	"accounts/api/app/model"
	"accounts/api/service"
	"accounts/api/utils"
	"accounts/api/utils/chain"
	"fmt"
	"strconv"
	"strings"

	db "accounts/api/platform/database/models"

	"github.com/gofiber/fiber/v2"
)

func GetUserByEmail(c *fiber.Ctx) {

}

func GetUserByNickname(c *fiber.Ctx) {

}

// AddWallet method for adding wallet to the corresponding user.
//
// @Description	지갑 등록 API 입니다.
// @Summary		지갑 등록 API
// @Tags		User
// @Accept		json
// @Produce		json
// @Param		Authorization	header	string	true	"액세스 토큰"
// @Param		type		path		string				true	"지갑 등록 타입 soft, hard"
// @Param		softWallet	body		model.AddSoftWallet	false	"soft wallet 등록 시 필요한 body 값"
// @Param		hardWallet	body		model.AddHardWallet	false	"hard wallet 등록 시 필요한 body 값"
// @Success		200			{object}	[]string{}
// @Router			/v1/user/wallet/{type} [post]
func AddWallet(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	newSessionKey := c.Locals("newSessionKey")
	addType := c.Params("type")
	fmt.Println(addType)
	switch addType {
	case "soft":
		var reqBody *model.AddSoftWallet
		c.BodyParser(&reqBody)
		verified, err := chain.Verify(reqBody.WalletAddr, reqBody.Salt, reqBody.Signature, reqBody.WalletType)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":        "error",
				"msg":           "error on verifying signature",
				"data":          err,
				"newSessionKey": newSessionKey,
			})
		}
		if verified {
			//insert into user_wallets
			if err := service.CreateSoftWallet(reqBody.WalletAddr, int32(userID)); err != nil {
				return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
					"status":        "error",
					"msg":           "error on inserting wallet",
					"data":          err,
					"newSessionKey": newSessionKey,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":        "ok",
				"msg":           "user wallet successfully added",
				"newSessionKey": newSessionKey,
			})
		} else {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status":        "error",
				"msg":           "verified signature does not met",
				"newSessionKey": newSessionKey,
			})
		}
	case "hard":
		var reqBody model.AddHardWallet
		c.BodyParser(&reqBody)
		addr, _, err := chain.GetPrivateKey(reqBody.PrivateKey)
		if err != nil {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status":        "error",
				"msg":           "error on validating privateKey",
				"data":          err,
				"newSessionKey": newSessionKey,
			})
		}
		if err := service.CreateHardWallet(addr.String(), reqBody.PrivateKey, int32(userID)); err != nil {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status":        "error",
				"msg":           "error on inserting wallet",
				"data":          err,
				"newSessionKey": newSessionKey,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":        "created",
			"msg":           fmt.Sprintf("address %s has been added", addr),
			"data":          addr,
			"newSessionKey": newSessionKey,
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":        "error",
			"msg":           "not allowed type",
			"newSessionKey": newSessionKey,
		})
	}

}

func FindWallets(c *fiber.Ctx) {

}

func DeleteWallet(c *fiber.Ctx) {

}

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	잔고 조회 API 입니다.
//	@Summary		잔고 조회 API
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			items[]	query		string				true	"조회 대상 id 쉼표로 구분합니다. ex) 1, 2, 3"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/balance [get]
func GetBalance(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	queryValue := strings.ReplaceAll(c.Query("items[]"), " ", "")
	selServices := utils.Filter(strings.Split(queryValue, ","), func(v string) bool {
		return v != ""
	})
	fmt.Println("('" + strings.Join(selServices, "','") + "')")
	res, err := service.GetBalances(int32(userID), fmt.Sprintln("('"+strings.Join(selServices, "','")+"')"))
	if err != nil {
		fmt.Println("Error()")
	}
	//fmt.Sprintln("('" + strings.Join(s, "','") + "')")
	return c.JSON(fiber.Map{
		"resp": res,
	})
}

// GetUserServices method for adding wallet to the corresponding user.
//
//	@Description	등록한 서비스 조회 API 입니다.
//	@Summary		등록한 서비스 조회 API
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/services [get]
func GetUserServices(c *fiber.Ctx) error {
	var serviceList *[]db.GetUserServicesRow
	var err error
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	serviceList, err = service.GetUserServices(int32(userID))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on fetching data",
			"data":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "fetched",
		"data":   serviceList,
	})

}

// CreateSecondPassword method for adding wallet to the corresponding user.
//
// @Description	2차 비밀번호 등록 API 입니다.
// @Summary		2차 비밀번호 등록 API
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			Authorization	header	string	true	"액세스 토큰"
// @Param		secondPass	body 	model.AddSecondPass	true	"2차 비밀번호"
// @Success		200		{object}	[]string{}
// @Router			/v1/user/secondpass [post]
func CreateSecondPassword(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	var txBody *model.AddSecondPass
	c.BodyParser(&txBody)
	if err := service.CreatedSecondPassword(int32(userID), txBody.SecPw); err != nil {
		return err
	}
	return nil

}
