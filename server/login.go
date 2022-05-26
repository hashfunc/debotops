package server

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/hashfunc/debotops/pkg/auth"
	"github.com/hashfunc/debotops/pkg/core"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		Error405(writer, nil)
	}

	body := request.Body
	defer func() {
		_ = body.Close()
	}()

	var payload LoginRequest
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		Error400(writer, err)
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

	err = bcrypt.CompareHashAndPassword(
		[]byte(root.PasswordHash),
		[]byte(payload.Password+root.SecretKey),
	)
	if err != nil {
		Error400(writer, err)
		return
	}

	refreshToken, err := auth.NewToken(root.SecretKey, auth.KindRefreshToken)
	if err != nil {
		Error500(writer, err)
		return
	}

	cookie := fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken)
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Set-Cookie", cookie)

	token, err := auth.NewToken(root.SecretKey, auth.KindRefreshToken)
	if err != nil {
		Error500(writer, err)
		return
	}

	response := map[string]string{
		"token": token,
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

	writer.WriteHeader(http.StatusOK)
}
