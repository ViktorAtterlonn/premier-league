package utils

import (
	"strconv"
	"strings"
)

func GetTeamFileName(teamName string) string {
	return strings.Replace(strings.ToLower(teamName), " ", "-", -1)
}

func Includes(str string, arr []string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}

	return false
}

func ToInt(str string) int {
	intVar, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}

	return intVar
}

func GetTeamFileExtension(teamName string) string {
	teamFileName := GetTeamFileName(teamName)
	exception := []string{"manchester-united", "newcastle-united"}

	if contains(exception, GetTeamFileName(teamFileName)) {
		return ".png"
	}

	return ".svg"
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}
