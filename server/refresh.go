package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashfunc/debotops/pkg/auth"
	"github.com/hashfunc/debotops/pkg/core"
)

func (server *Server) refresh(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		Error405(writer, nil)
		return
	}

	refreshTokenCookie, err := request.Cookie("refresh-token")
	if err != nil {
		Error400(writer, nil)
		return
	}

	secret, err := server.getDeBotOpsSecret()
	if err != nil {
		Error500(writer, err)
		return
	}

	secretData := secret.Data

	var root core.Root
	if err := json.Unmarshal(secretData["root"], &root); err != nil {
		Error500(writer, err)
		return
	}

	token, err := jwt.ParseWithClaims(refreshTokenCookie.Value, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(root.SecretKey), nil
	})
	if err != nil {
		Error500(writer, err)
		return
	}

	claims := token.Claims.(*auth.Claims)

	if err := claims.Valid(); err != nil {
		Error400(writer, err)
		return
	}

	refresh, err := auth.NewAuth(root.SecretKey, auth.KindRefresh)
	if err != nil {
		Error500(writer, err)
		return
	}

	expires := refresh.Expires.UTC().Format(http.TimeFormat)
	cookie := fmt.Sprintf("refresh-token=%s; Expires=%s; HttpOnly", refresh, expires)
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Set-Cookie", cookie)

	access, err := auth.NewAuth(root.SecretKey, auth.KindAccess)
	if err != nil {
		Error500(writer, err)
		return
	}

	response := map[string]string{
		"token": access.String(),
	}

	responseBody, err := json.Marshal(&response)
	if err != nil {
		Error500(writer, err)
		return
	}

	if _, err := writer.Write(responseBody); err != nil {
		Error500(writer, err)
		return
	}
}
