package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)



func HashPassword(password string) (hashedPass string ,e error){

	hashedPass , e = argon2id.CreateHash(password , argon2id.DefaultParams)
	if e != nil {
		return "" , e
	}
	
	return hashedPass , nil
}

func CheckPassword(password , hash string) (bool , error){

	match ,err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("Error comparing password and hash: %v" , err)
		return false , err
	}

	if match {
		return true , nil
	}else{
		return false , nil
	}
	
}