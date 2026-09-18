package main

import (
	"net/http"
	"time"

	"github.com/AnasK239/Chirpy/internal/auth"
)


type RefreshResponse struct {
	NewAccessToken 	string `json:"token"`
}

func (cfg *apiConfig) HandleRefreshAccessToken(w http.ResponseWriter , r *http.Request) {

	refreshToken , err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbRefreshToken , err := cfg.dbQueries.FindRefreshTokenById(r.Context() , refreshToken)
	if err != nil || dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized , "Non-existent / revoked refreshToken")
		return
	}

	newAccessToken ,err := auth.MakeJWT(dbRefreshToken.UserID , cfg.jwtSecret , JwtAccessTokenExpiration)
	if err != nil {
		respondWithError(w , 500 , "Couldn't generate a new access token")
		return
	}
	
	respondWithJSON(w, 200 , RefreshResponse{NewAccessToken: newAccessToken})
	
}


func (cfg *apiConfig) HandleRevokeRefreshToken(w http.ResponseWriter , r *http.Request){
	refreshToken , err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbRefreshToken , err := cfg.dbQueries.FindRefreshTokenById(r.Context() , refreshToken)
	if err != nil || dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized , "Non-existent / revoked refreshToken")
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context() , dbRefreshToken.Token)
	if err != nil {
		respondWithError(w , 500 , "Couldn't Revoke RefreshToken")
		return
	}

	w.WriteHeader(204)
}