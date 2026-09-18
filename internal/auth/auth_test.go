package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)



func TestValidateJWT(t *testing.T) {
	secret := "190uej19hndb2uidsn8sbnuyijhbv27TYCSVtysc2CS6"
	userID := uuid.New()
	jwt , _ := MakeJWT(userID , secret ,  time.Hour)

	resID , err := ValidateJWT(jwt , secret)

	if resID != userID || err != nil {
		t.Errorf("RESULT MISMATCH")
	}
}


// func TestMakeJWT(t *testing.T) {

// }
// 
// 

func TestGetBearerToken(t *testing.T) {

	// Normal Header
	headers := http.Header{}
	
	headers.Add("Authorization" , "Bearer ok")
	expected := "ok"

	result , _ := GetBearerToken(headers)
	if result != expected {
		t.Errorf("Couldn't extract bearer")
	}

	// Malformed Header
	headers2 := http.Header{}

	headers2.Add("Authorization" , "malformed")
	
	exp := "Malformed Auth information"

	result , err := GetBearerToken(headers2)
	if err.Error() != exp {
		t.Errorf("Incorrect token not detected")
	}
}