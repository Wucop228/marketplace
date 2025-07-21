package http

import (
	"database/sql"
	"github.com/Wucop228/marketplace/internal/config"
	"github.com/Wucop228/marketplace/internal/models"
	"github.com/Wucop228/marketplace/internal/repo"
	"github.com/Wucop228/marketplace/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"time"
)

type AuthHandler struct {
	db         *sql.DB
	AuthConfig *config.AuthConfig
}

func NewAuthHandler(db *sql.DB, authConfig *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		db:         db,
		AuthConfig: authConfig,
	}
}

func GenerateAccessToken(h *AuthHandler, c echo.Context, user models.User) (string, error) {
	accessExp := time.Now().Add(h.AuthConfig.AccessTokenTTL).Unix()
	accessClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": accessExp,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenSigned, err := accessToken.SignedString([]byte(h.AuthConfig.JWTSecret))
	if err != nil {
		return "", c.JSON(http.StatusInternalServerError, echo.Map{"error": "error while generating access token"})
	}

	return accessTokenSigned, nil
}

func CheckPasswordAndUsername(username string, password string) echo.Map {
	if len(username) < 3 || len(username) > 20 {
		return echo.Map{"error": "username must be 3-20 characters"}
	}
	if len(password) < 8 || len(password) > 20 {
		return echo.Map{"error": "password must be 8-20 characters"}
	}

	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", password); !matched {
		return echo.Map{"error": "username contains invalid characters"}
	}
	if matched, _ := regexp.MatchString("[0-9]", password); !matched {
		return echo.Map{"error": "password must contain at least one digit"}
	}

	return nil
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Username string
		Password string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	if err := CheckPasswordAndUsername(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	hashPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "password hashing failed"})
	}

	user := models.User{
		ID:           -1,
		Username:     req.Username,
		PasswordHash: string(hashPassword),
	}
	if err := repo.CreateUser(h.db, &user); err != nil || user.ID == -1 {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "user creation failed"})
	}

	exists, err := repo.UserExists(h.db, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "user creation failed"})
	}
	if exists {
		return c.JSON(http.StatusConflict, echo.Map{"error": "username already exists"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Username string
		Password string
	}
	if err := c.Bind(&req); nil != err {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	var user models.User
	user, err := repo.LoginUser(h.db, req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error while logging in"})
	}
	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "username or password is incorrect"})
	}
	if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "username or password is incorrect"})
	}

	accessTokenSigned, err := GenerateAccessToken(h, c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "token generation failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessTokenSigned,
	})
}
