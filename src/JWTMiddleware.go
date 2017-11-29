package main

import (
	//"net/http"
	"strings"
	"fmt"
	"errors"
	"log"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTMiddleware struct {
	username string
	password string
}

type CustomizedClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var DEBUG = false
var mySigningKey = []byte("ETCD_ACCESS")

func OnError(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

func (m *JWTMiddleware) logf(format string, args ...interface{}) {
	if DEBUG {
		log.Printf(format, args...)
	}
}

func JWTMiddlewareNew () *JWTMiddleware{
	return &JWTMiddleware{
		"etcd",
		"etcd123",
	}
}

func (m *JWTMiddleware) parse (signedstring string, w http.ResponseWriter, r *http.Request) (*jwt.Token, error){
	token, err := jwt.Parse(signedstring, func(token *jwt.Token) (interface{}, error) {
    // Don't forget to validate the alg is what you expect:
	    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	    }

	    return mySigningKey, nil
	})

	if err != nil {
		m.logf("Error when parsing")
		OnError(w, r, "Error when parsing")
		return nil, errors.New("Error when parsing")
	}

	if !token.Valid {
		m.logf("Token is invalid")
		OnError(w, r, "The token isn't valid")
		return nil, errors.New("Token is invalid")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["username"] != "etcd" && 
			claims["password"] != "etcd123" {
				m.logf("username or password not correct")
				OnError(w, r, "username or password not correct")
				return nil, errors.New("username or password not correct")
			}

    	fmt.Println(claims["username"], claims["password"])
    	return token, nil
	} else {
	    m.logf("can not get claims")
		OnError(w, r, "can not get claims")
		return nil, errors.New("can not get claims")
	}
}

// FromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func (m *JWTMiddleware) ServeHTTP (rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tokenString, err := FromAuthHeader(r)
	// If the tokenString is empty...
	if tokenString == "" {
		// If we get here, the required tokenString is missing
		errorMsg := "Required authorization tokenString not found"
		OnError(rw, r, errorMsg)
		m.logf("  Error: No credentials found (CredentialsOptional=false)")
	}

	token, err := m.parse(tokenString, rw, r)
	// Check if there was an error in parsing...
	if err != nil || token == nil {
		m.logf("Error parsing tokenString: %v", err)
		OnError(rw, r, err.Error())
	}
	
	// If there was an error, do not call next.
	if err == nil && next != nil {
		next(rw, r)
	}
}