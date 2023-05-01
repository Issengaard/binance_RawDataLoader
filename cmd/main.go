package main

import (
	"github.com/issengaard/binance_raw_data_loader/internal/constants"
	"github.com/issengaard/binance_raw_data_loader/internal/services/limited_file_loader"
	"github.com/issengaard/binance_raw_data_loader/internal/utils"
	"os"
)

func main() {
	//cli_cmd.InitRootCmd()

	filePath, _ := os.Getwd()
	fileLink := "https://data.binance.vision/data/spot/monthly/trades/BTCUSDT/BTCUSDT-trades-2017-12.zip"

	limiFileLoader := limited_file_loader.New(5000)
	err := limiFileLoader.Download(fileLink, filePath)
	if err != nil {
		utils.ToConsole(constants.APP_NAME, err.Error())
	}
}
