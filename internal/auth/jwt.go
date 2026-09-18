package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)




func MakeJWT(userID uuid.UUID , tokenSecret string , expiresIn time.Duration) (string , error){

	issuer := "chirpy-access"
	subject := userID.String()

	issuet := time.Now().UTC()
	expirationt := issuet.Add(expiresIn)
	
	issuedAt := jwt.NewNumericDate(issuet)
	expiresAt := jwt.NewNumericDate(expirationt)
	
	claims := jwt.RegisteredClaims{
		Issuer: issuer,
		IssuedAt: issuedAt,
		ExpiresAt: expiresAt,
		Subject: subject,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken , err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "" , err
	}

	return jwtToken , nil
}



func ValidateJWT(tokenString , tokenSecret string) (uuid.UUID , error){

	claims := jwt.RegisteredClaims{}
	token , err := jwt.ParseWithClaims(tokenString, &claims ,func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

	if err != nil {
		log.Printf("Error parsing jwt: %v" ,err)
		return uuid.UUID{} , err
	}

	userID , err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting subject from jwt: %v" , err)
		return uuid.UUID{} , err
	}

	id , err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v" , err)
		return uuid.UUID{} , err
	}
	
	return  id , nil
}