package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmcvetta/napping"
)

type CustomizedClaims struct {
        Username string `json:"username"`
        Password string `json:"password"`
        jwt.StandardClaims
}

var mySigningKey = []byte("ETCD_ACCESS")

func newToken() string {
	claims := CustomizedClaims {
		"etcd",
		"etcd123",
		jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println("%v %v", ss, err)
	return ss
}

func main() {
	token := newToken()
	auth := "bearer "+ token
	fmt.Println(auth);
	url := "http://127.0.0.1:3000"

	s := napping.Session{}
	h := &http.Header{}
	h.Set("Authorization",auth)
	s.Header = h

	var jsonStr = []byte(`{"title": "Test JWT"}`)

	var data map[string]json.RawMessage
    err := json.Unmarshal(jsonStr, &data)
    if err != nil {
        fmt.Println(err)
    }

    resp, err := s.Post(url, &data, nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("response Status:", resp.Status())
    fmt.Println("response Headers:", resp.HttpResponse().Header)
    fmt.Println("response Body:", resp.RawText())
}