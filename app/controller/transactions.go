package controller

import (
	"accounts/api/app/model"
	"accounts/api/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	온/오프체인 송금 API 입니다. 온/오프체인 기록을 전송하고 싶을 경우 body에 추가.
//	@Summary		온/오프체인 송금 API
//	@Tags			Transaction
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			wallet	body		model.RECVTransfer		true	"송금 시도 시 필요한 body 값"
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
	var txBody *model.RECVTransfer

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

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	거래내역 조회 API 입니다. 상세 거래내역 조회를 통해 상세 내역을 확인할 수 있습니다.
//	@Summary		거래내역 조회 API
//	@Tags			Transaction
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			lim	query	string	false	"조회 limit"
//	@Param			off	query	string	false	"조회 offset"
//	@Success		200	{object}	[]string{}
//	@Router			/v1/tx/history [get]
func TransferHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	lim := strings.ReplaceAll(c.Query("lim"), " ", "")
	off := strings.ReplaceAll(c.Query("off"), " ", "")
	var respdata []model.TransactionsList
	if lim == "" {
		lim = "10"
	}
	if off == "" {
		off = "0"
	}

	if _, ok := strconv.Atoi(lim); ok != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "error",
			"msg":    "given limit value is incorrect",
			"data":   ok.Error(),
		})
	}

	if _, ok := strconv.Atoi(off); ok != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "error",
			"msg":    "given offset value is incorrect",
			"data":   ok.Error(),
		})
	}

	result, err := service.TransactionHistory(userID, lim, off)
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on fetching history",
			"data":   err.Error(),
		})
	}

	for _, row := range result {
		isNative := false
		if row.IsNative.Int32 == 1 {
			isNative = true
		}
		respdata = append(respdata, model.TransactionsList{
			TransactionID: row.ID,
			FromID:        row.FromID,
			ToID:          row.ToID,
			Memo:          row.Memo.String,
			TotalAmount:   row.TotalAmount,
			TokenInfo: model.TokenInfo{
				Name:       row.Name.String,
				Symbol:     row.Symbol.String,
				Decimals:   row.Decimals.Int32,
				Image:      row.Image.String,
				IsNative:   isNative,
				NetType:    row.NetType.String,
				WalletType: row.WalletType.String,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
		"data":   respdata,
	})
}

// AddWallet method for adding wallet to the corresponding user.
//
//	@Description	거래내역 상세 조회 API 입니다. 상세 거래내역 조회를 통해 상세 내역을 확인할 수 있습니다.
//	@Summary		거래내역 상세 조회 API
//	@Tags			Transaction
//	@Accept			json
//	@Produce		json
//	@Param			id				path	string	true	"트랜잭션 ID 값"
//	@Success		200	{object}	[]string{}
//	@Router			/v1/tx/history/detail/{id} [get]
func TransactionHistoryDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on transaction id conversion",
		})
	}
	result, err := service.TransactionDetail(int32(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on fetching details",
		})
	}
	var respdata []model.TransactionDetail
	for _, row := range result {
		isOnChain := false
		if row.IsOnchain == 1 {
			isOnChain = true
		}
		respdata = append(respdata, model.TransactionDetail{
			From:      row.From,
			To:        row.To,
			Amount:    row.Amount,
			IsOnchain: isOnChain,
			Txhash:    row.Txhash.String,
			Status:    row.Status,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
		"data":   respdata,
	})
}
