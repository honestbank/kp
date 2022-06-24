package kp

import (
	"regexp"
	"strconv"
)

func UnmarshalStringMessage(message string) (string, int, error) {
	re := regexp.MustCompile("\\|(\\d+)$")
	findRetries := re.FindAllString(message, -1)
	if len(findRetries) == 0 {
		return message, 0, nil
	}
	foundRetry := string(findRetries[0])
	if foundRetry != "" {
		retries, err := strconv.Atoi(foundRetry[1:])
		if err != nil {
			return message, 0, err
		}
		return message[0 : len(message)-len(foundRetry)], retries, nil
	}

	return message, 0, nil
}

func MarshalStringMessage(message string, retries int) string {
	if retries == 0 {
		return message
	}
	return message + "|" + strconv.Itoa(retries)
}
