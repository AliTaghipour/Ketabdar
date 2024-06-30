package core

import (
	"Ketab/model"
	"Ketab/repository"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"time"
)

type Claims struct {
	UserId int32
	jwt.StandardClaims
}

type Handler struct {
	usersRepo     repository.UsersRepository
	booksRepo     repository.BookRepository
	userBooksRepo repository.UserBooksRepository
	jwtKey        []byte
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

func (h *Handler) AddBook(w http.ResponseWriter, r *http.Request) {
	_, err := h.authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req := &model.RequestAddBook{}
	h.getRequestBody(r.Body, req)

	err = h.booksRepo.AddNewBook(req.Book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("afarin"))
	return

}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	_, err := h.authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	books, err := h.booksRepo.GetBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
	return
}

func (h *Handler) AddBookToLibrary(w http.ResponseWriter, r *http.Request) {
	userId, err := h.authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var req model.RequestAddBookLibrary
	h.getRequestBody(r.Body, &req)

	err = h.userBooksRepo.AddNewUserBook(&req.UserBook, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("afarin"))
	return
}

func (h *Handler) GetLibraryBooks(w http.ResponseWriter, r *http.Request) {
	userId, err := h.authorize(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	books, err := h.userBooksRepo.GetUserBooks(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshal)
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

func (h *Handler) authorize(request *http.Request) (int32, error) {
	token, err := getTokenFromRequest(request)
	if err != nil {
		log.Println("fail to get cookie")
		return 0, err
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil
	})
	if err != nil {
		log.Println("bad credentials")
		return 0, err
	}
	err = tkn.Claims.Valid()
	if err != nil {
		log.Println("invalid claim")
		return 0, err
	}
	if !tkn.Valid {
		log.Println("invalid token")
		return 0, nil
	}
	return claims.UserId, nil
}

func getTokenFromRequest(request *http.Request) (string, error) {
	c, err := request.Cookie("token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
