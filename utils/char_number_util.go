package utils

import (
	"log"
	"strconv"
)

// GetFloat64sFromString 获取字符串中的数字
func GetFloat64sFromString(str string) []float64 {
	var result []float64
	var temp string
	for _, char := range str {
		if char >= '0' && char <= '9' {
			temp += string(char)
		} else {
			if temp != "" {
				last := temp[len(temp)-1]
				if char == '.' && last >= '0' && last <= '9' {
					temp += string(char)
				} else {
					result = append(result, ParseFloat64(temp))
					temp = ""
				}
			}
		}
	}
	if temp != "" {
		result = append(result, ParseFloat64(temp))
	}
	return result
}

// ParseFloat64 将字符串转换为float64
func ParseFloat64(str string) float64 {
	rel, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println("StrToFloat64 err:", err)
		return 0
	}
	return rel
}
