package controller

import (
	"accounts/api/service"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AddService method to create service
//
//	@Description	서비스 추가 API 입니다.
//	@Summary		서비스 추가 API
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"액세스 토큰"
//	@Param			serviceId	path		string				true	"추가 대상 서비스 id"
//	@Success		200		{object}	[]string{}
//	@Router			/v1/user/service/{serviceId} [post]
func AddService(c *fiber.Ctx) error {
	userID, _ := strconv.ParseInt(fmt.Sprintf("%v", c.Locals("userID")), 10, 64)
	serviceID, _ := strconv.ParseInt(c.Params("serviceId"), 10, 64)
	result, err := service.UserAddService(int32(userID), int32(serviceID))
	if !result {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error on fetching",
			"data":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "added",
		"msg":    "Service added successfully",
	})
}

func DeleteService(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "deleted",
		"msg":    "service has deleted",
	})
}

type RespServices struct {
	ServiceID    int32  `json:"serviceId"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Decimals     int32  `json:"decimals"`
	Image        string `json:"img"`
	IsNative     bool   `json:"isNative"`
	ContractAddr string `json:"contractAddr"`
	NetType      string `json:"netType"`
}

// AddService method to create service
//
//	@Description	서비스 리스트 API 입니다.
//	@Summary		서비스 리스트 API
//	@Tags			Info
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]string{}
//	@Router			/v1/info/services [get]
func GetTotalServices(c *fiber.Ctx) error {
	res, err := service.GetAllServices()
	var respdata []RespServices
	for _, value := range *res {
		isNative := false
		if value.IsNative == 1 {
			isNative = true
		}
		respdata = append(respdata, RespServices{
			ServiceID:    value.ID,
			Name:         value.Name,
			Symbol:       value.Symbol,
			Image:        value.Image.String,
			Decimals:     value.Decimals,
			ContractAddr: value.ContractAddr.String,
			NetType:      value.NetType.String,
			IsNative:     isNative,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "error while fetching",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "fetched",
		"data":   respdata,
	})
}
