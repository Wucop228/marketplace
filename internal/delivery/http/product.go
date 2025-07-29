package http

import (
	"database/sql"
	"github.com/Wucop228/marketplace/internal/repo"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ProductHandler struct {
	db *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func CheckReqCreateProduct(header string, image_url string, text string, price float64) echo.Map {
	if len(header) == 0 {
		return echo.Map{"error": "header must not be empty"}
	}
	if len(header) >= 100 {
		return echo.Map{"error": "header must be less than 100 characters"}
	}

	lowerURL := strings.ToLower(image_url)
	if !strings.HasPrefix(lowerURL, "http://") && !strings.HasPrefix(lowerURL, "https://") {
		return echo.Map{"error": "invalid URL format"}
	}

	if len(text) == 0 {
		return echo.Map{"error": "text must not be empty"}
	}
	if len(text) > 2000 {
		return echo.Map{"error": "text must be less than 2000 characters"}
	}

	if price == 0 {
		return echo.Map{"error": "price must not be empty"}
	}

	return nil
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req struct {
		Price    float64
		Header   string
		Text     string
		ImageURL string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request format"})
	}

	errMap := CheckReqCreateProduct(req.Header, req.ImageURL, req.Text, req.Price)
	if errMap != nil {
		return c.JSON(http.StatusBadRequest, errMap)
	}

	userID, ok := c.Get("sub").(float64)
	if !ok || userID == 0.0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "userID must not be empty"})
	}

	username, err := repo.GetUserById(h.db, int64(userID))

	err = repo.CreateProduct(h.db, req.Price, req.Header, req.Text, username, req.ImageURL)
	if err != nil {
		//return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error when creating product"})
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"product": req})
}

//func (h *ProductHandler) GetProducts(c *echo.Context) error {
//
//}
