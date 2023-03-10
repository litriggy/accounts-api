package service

import (
	"accounts/api/pkg/config"
	db "accounts/api/platform/database/models"
	"context"
)

func GetService(id int32) (*db.GetServiceDataRow, error) {
	ctx := context.Background()
	query, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	serviceInfo, err := query.GetServiceData(ctx, id)
	if err != nil {
		return nil, err
	}

	return &serviceInfo, nil
}

func GetAllServices() (*[]db.Service, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	servicesList, err := query.GetAllServices(context.Background())
	if err != nil {
		return nil, err
	}
	return &servicesList, nil
}

func GetUserServices(userID int32) (*[]db.GetUserServicesRow, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	serviceList, err := query.GetUserServices(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return &serviceList, nil
}

func UserAddService(userID int32, serviceID int32) (bool, error) {
	query, err := config.ConnectDB()
	if err != nil {
		return false, err
	}
	if err := query.UserAddService(context.Background(), db.UserAddServiceParams{UserID: userID, ServiceID: serviceID}); err != nil {
		return false, err
	}
	return true, nil

}
