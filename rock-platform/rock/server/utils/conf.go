package utils

import "strings"

func GetConfigName(configFile string) (fileName string) {
	fileNameList := strings.Split(configFile, ".")
	fileName = fileNameList[0]
	return
}
