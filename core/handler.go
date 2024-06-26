package core

import (
	"Ketab/model"
	"Ketab/repository"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"time"
)

type Claims struct {
	UserId int32
	jwt.StandardClaims
}

type Handler struct {
	usersRepo repository.UsersRepository
	jwtKey    []byte
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.RequestRegister
	h.getRequestBody(r.Body, &req)

	user, err := h.usersRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return
	}
	if user != nil {
		return
	}

	err = h.usersRepo.AddUser(context.Background(), &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("afarin"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.RequestLogin
	h.getRequestBody(r.Body, &req)

	user, err := h.usersRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return
	}

	if user.Password != req.Password {
		return
	}

	expirationTime := time.Now().UTC().Add(3 * time.Hour)
	claims := &Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    tokenString,
		Expires:  expirationTime,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) getRequestBody(reader io.ReadCloser, req any) {
	all, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	err = json.Unmarshal(all, &req)
	if err != nil {
		return
	}
}
