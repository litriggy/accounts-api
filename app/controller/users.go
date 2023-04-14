package controller

import (
	"accounts/api/app/model"
	"accounts/api/service"
	"accounts/api/utils"
	"accounts/api/utils/chain"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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
	addType := c.Params("type")
	switch addType {
	case "soft":
		var reqBody model.AddSoftWallet
		c.BodyParser(&reqBody)

		exists, err := service.FindWallet(reqBody.WalletAddr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on verifying signature",
				"data":   err.Error(),
			})
		}
		if !exists {
			fmt.Println("exists")
		}
		fmt.Println(reqBody)
		fmt.Println(fmt.Sprintf("Welcome to Crea\nNonce:" + reqBody.Salt + "."))
		verified, err := chain.Verify(reqBody.WalletAddr, fmt.Sprintf("Welcome to Crea\nNonce:"+reqBody.Salt+"."), reqBody.Signature, reqBody.WalletType)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on verifying signature",
				"data":   err.Error(),
			})
		}
		if verified {

			if err := service.CreateSoftWallet(reqBody.WalletAddr, int32(userID)); err != nil {
				return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
					"status": "error",
					"msg":    "error on inserting wallet",
					"data":   err,
				})
			}
			//insert into user_wallets
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "ok",
				"msg":    "user wallet successfully added",
			})
		} else {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status": "error",
				"msg":    "verified signature does not met",
			})
		}
	case "hard":
		var reqBody model.AddHardWallet
		c.BodyParser(&reqBody)

		hashed, _ := chain.Hash(reqBody.SecPw)
		passChk, err := service.CheckPassword(int32(userID), hashed)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on validating second password",
				"data":   err.Error(),
			})
		}
		if !passChk {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"msg":    "second password doesn't match",
			})
		}

		privateKey, err := chain.Decrypt(reqBody.PrivateKey, reqBody.SecPw)
		if err != nil {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on decrypting privateKey",
				"data":   err.Error(),
			})
		}
		addr, _, err := chain.GetPrivateKey(privateKey)
		if err != nil {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on validating privateKey",
				"data":   err,
			})
		}
		if err := service.CreateHardWallet(addr.String(), reqBody.PrivateKey, int32(userID)); err != nil {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"status": "error",
				"msg":    "error on inserting wallet",
				"data":   err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "created",
			"msg":    fmt.Sprintf("address %s has been added", addr),
			"data":   addr,
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "not allowed type",
		})
	}

}

func FindWallets(c *fiber.Ctx) {

}

func DeleteWallet(c *fiber.Ctx) {

}

// MyInfo method for fetching user's info
//
//	@Description	내 정보 조회 API 입니다.
//	@Summary		내 정보 조회 API 입니다.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/myinfo [get]

func MyInfo(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	//var userInfo model.UserInfo
	resp, _ := service.GetUser(int32(userID))
	fmt.Println(resp)
	return c.JSON(fiber.Map{
		"data": model.UserInfo{
			ID:       resp.ID,
			Nickname: resp.Nickname.String,
			Email:    resp.Email,
			Picture:  resp.Picture.String,
		},
	})
}

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	잔고 조회 API 입니다.
//	@Summary		오프체인 특정 서비스들의 잔고 조회 API
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			services[]	query		string				true	"조회 대상 id 쉼표로 구분합니다. ex) 1, 2, 3"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/balance [get]
func GetBalance(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	queryValue := strings.ReplaceAll(c.Query("services[]"), " ", "")
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

type RespMyServices struct {
	ServiceID    int32          `json:"serviceId"`
	Name         string         `json:"name"`
	Symbol       string         `json:"symbol"`
	Decimals     int32          `json:"decimals"`
	Image        string         `json:"img"`
	IsNative     bool           `json:"isNative"`
	ContractAddr string         `json:"contractAddr"`
	NetType      string         `json:"netType"`
	Balance      int64          `json:"balance"`
	Wallets      []WalletHolder `json:"wallets"`
}
type WalletHolder struct {
	Addr         string `json:"walletAddr"`
	IsIntegrated int32  `json:"isIntegrated"`
}

// GetUserServices method for adding wallet to the corresponding user.
//
//	@Description	액세스 토큰 기반하여 유저 검색 후 등록한 서비스 조회 API 입니다.
//	@Summary		등록한 서비스 조회 API
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/services [get]
func GetUserServices(c *fiber.Ctx) error {
	//var serviceList []db.GetUserServicesRow
	var err error
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	serviceList, err := service.GetUserServices(int32(userID))

	var respdata []RespMyServices

	//var res map[string]RespMyServices

	for _, row := range *serviceList {
		isNative := false
		exists := -1
		if row.IsNative.Int32 == 1 {
			isNative = true
		}

		for num, scan := range respdata {
			if scan.ServiceID == row.ServiceID {
				exists = num
			}
		}

		if exists == -1 {
			respdata = append(respdata, RespMyServices{
				ServiceID:    row.ServiceID,
				Name:         row.Name.String,
				Symbol:       row.Symbol.String,
				Decimals:     row.Decimals.Int32,
				Image:        row.Image.String,
				IsNative:     isNative,
				ContractAddr: row.ContractAddr.String,
				NetType:      row.NetType.String,
				Balance:      row.Amount.Int64,
				Wallets: []WalletHolder{
					{
						Addr:         row.WalletAddr.String,
						IsIntegrated: row.IsIntegrated.Int32,
					},
				},
			})

		} else {
			respdata[exists].Wallets = append(respdata[exists].Wallets, WalletHolder{
				Addr:         row.WalletAddr.String,
				IsIntegrated: row.IsIntegrated.Int32,
			})
		}

	}
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on fetching data",
			"data":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "fetched",
		"data":   respdata,
	})

}

// CreateSecondPassword method for adding wallet to the corresponding user.
//
// @Description	2차 비밀번호 (평문)을 받아 해싱 후 저장합니다.
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
	hashed, err := chain.Hash(txBody.SecPw)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    err.Error(),
		})
	}
	if err := service.CreatedSecondPassword(int32(userID), hashed); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "created",
		"msg":    "second password has added",
	})

}

// Check if second pass exists
//
// @Description	2차 비밀번호의 존재 여부를 확인합니다.
// @Summary		2차 비밀번호 존재 확인 API
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			Authorization	header	string	true	"액세스 토큰"
// @Success		200		{object}	[]string{}
// @Router			/v1/user/secondpass [get]
func CheckSecondPass(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    err.Error(),
		})
	}
	exists, err := service.CheckPasswordExists(int32(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "ok",
		"msg":    "second password exists",
		"exist":  exists,
	})

}
