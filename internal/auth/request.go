package auth

import (
	"net/http"
	"errors"
	"strings"
)


func GetBearerToken(headers http.Header) (string , error){
	authHeader , ok := headers["Authorization"]
	if !ok || len(authHeader) == 0 {
		return "" , errors.New("Auth information not present")
	}

	parts := strings.Split(authHeader[0] , " ")

	if len(parts) != 2 {
		return "" , errors.New("Malformed Auth information")
	}
	
	token := parts[1]
	
	return token , nil
}