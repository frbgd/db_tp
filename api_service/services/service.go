package services

import (
	"db_tp/db"
	"db_tp/models"
)

var ServiceSrv ServiceService

type ServiceService struct {
	db *db.PostgresDbEngine
}

func (forum *ServiceService) GetStatus() (*models.Status, error) {

}

func (forum *ServiceService) ClearDatabase() error {

}
