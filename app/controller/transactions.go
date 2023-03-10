package controller

import (
	"accounts/api/app/model"
	"accounts/api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	오프체인 송금 API 입니다.
//	@Summary		오프체인 송금 API
//	@Tags			Transaction
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			wallet	body		model.Transfer		true	"송금 시도 시 필요한 body 값"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/tx/transfer [post]
func TransferBalance(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Locals("userID").(string))
	//c.BodyParser()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "unidentified userId",
			"data":   err.Error(),
		})
	}
	var txBody *model.Transfer

	if err := c.BodyParser(&txBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "unable to parse body",
			"data":   err.Error(),
		})
	}
	if err := service.TransferBalance(int32(userID), int32(txBody.Target), int32(txBody.ServiceID), txBody.OnchainEvent, txBody.OffchainEvent, txBody.SecPw); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "unable to transfer",
			"data":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "created",
		"msg":    "Service created successfully",
	})
}

// func TransferOnchain(c *fiber.Ctx) error {
// 	var reqBody *model.TransferOnchain
// 	if len(reqBody.Amount) != len(reqBody.Sender) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status": "error",
// 			"msg":    "length of sending info doesnt match",
// 		})

// 	}
// 	for i, el := range reqBody.Sender {
// 		//get PrivateKey of such Address from el
// 		var privateKey string
// 		// PK part should be implemented
// 		privateKey = el
// 		if err := chain.Transfer(privateKey, reqBody.Target, reqBody.Amount[i], int32(reqBody.ServiceID)); err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"status": "error",
// 				"msg":    "error on sending transaction",
// 				"data":   err,
// 			})
// 		}
// 	}
// 	if err := c.BodyParser(&reqBody); err != nil {
// 		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 			"status": "error",
// 			"msg":    "unidentified userId",
// 			"data":   err.Error(),
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status": "OK",
// 	})
// }

func TransferFromBalance(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "created",
		"msg":    "Service created successfully",
	})
}

func BalanceOf(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "created",
		"msg":    "Service created successfully",
	})
}

func TransferHistory(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Locals("userID").(string))
	sessionKey := c.Locals("newSessionKey")
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "error",
			"msg":    "unidentified userId",
			"data":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "OK",
		"sessionKey": sessionKey,
		"userId":     userID,
	})
}
