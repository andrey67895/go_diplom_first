package services

import (
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
)

func OrdersStatusJob() {
	for {
		_, err := database.DBStorage.GetOrdersByNotFinalStatus()
		if err != nil {
			helpers.TLog.Error(err.Error())
			return
		}

		//status.
	}

}
