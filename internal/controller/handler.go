package controller

import (
	"final/internal/domain"
	auth "final/internal/midleware"
	"final/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	service usecase.Usecase
}

func NewHandler(service usecase.Usecase) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle() http.Handler {
	router := gin.Default()

	router.POST("/pvz", auth.AuthMiddleware(), auth.RequireRole("moderator"), h.createPVZ)
	router.GET("/pvz", auth.AuthMiddleware(), auth.RequireRole("employee", "moderator"), h.getPVZList)
	router.POST("/pvz/:pvzId/close_last_reception", auth.AuthMiddleware(), auth.RequireRole("employee"), h.closeLastReception)
	router.POST("/pvz/:pvzId/delete_last_product", auth.AuthMiddleware(), auth.RequireRole("employee"), h.deleteLastProductForPVZ)

	router.POST("/receptions", auth.AuthMiddleware(), auth.RequireRole("employee"), h.createReception)

	router.POST("/products", auth.AuthMiddleware(), auth.RequireRole("employee"), h.createProduct)

	router.POST("/dummyLogin", h.dummyLogin)
	router.POST("/register", h.register)
	router.POST("/login", h.login)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented,
			domain.ErrorResponse{Message: "not implemented"})
	})
	return router
}

func (h *Handler) createPVZ(c *gin.Context) {
	var pvz *domain.PVZ
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&pvz); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}
	pvz, err := h.service.PVZUsecase.CreatePVZ(pvz)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error creating pvz: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pvz)
}

func (h *Handler) createReception(c *gin.Context) {
	var reception *domain.Reception
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&reception); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}
	reception, err := h.service.ReceptionUsecase.CreateReception(reception)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error creating reception: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reception)
}

func (h *Handler) closeLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")

	c.Header("Content-Type", "application/json")
	reception, err := h.service.ReceptionUsecase.CloseReception(&pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error closing reception: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, reception)
}

func (h *Handler) createProduct(c *gin.Context) {
	var inputProduct *domain.InputProduct
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputProduct); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}
	product, err := h.service.ProductUsecase.CreateProduct(inputProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error creating product: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) deleteLastProductForPVZ(c *gin.Context) {
	pvzID := c.Param("pvzId")

	err := h.service.ProductUsecase.DeleteLastProductForPVZ(&pvzID)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error deleting product: " + err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) getPVZList(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	pvzList, err := h.service.PVZUsecase.GetPVZInfo(startDateStr, endDateStr, pageStr, limitStr)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error getting pvz info: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, pvzList)
}

func (h *Handler) dummyLogin(c *gin.Context) {
	var user domain.User
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}
	if auth.CheckRole(*user.Role) == false {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect role"})
		return
	}

	Id := uuid.New().String()
	token, err := auth.GenerateJWT(Id, *user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error generating token: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *Handler) register(c *gin.Context) {
	var inputUser domain.InputUser
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}
	if auth.CheckRole(*inputUser.Role) == false {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "incorrect role"})
		return
	}

	user, err := h.service.UserUsecase.Register(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error creating user: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) login(c *gin.Context) {
	var inputUser domain.InputUser
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "incorrect json format: " + err.Error()})
		return
	}

	user, err := h.service.UserUsecase.Login(&inputUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "error login user: " + err.Error()})
		return
	}
	token, err := auth.GenerateJWT(*user.ID, *user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, token)
}
