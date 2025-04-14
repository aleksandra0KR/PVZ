package controller

import (
	"final/internal/domain"
	"final/internal/metrics"
	auth "final/internal/midleware"
	"final/internal/usecase"
	"github.com/gin-gonic/gin"
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

	router.Use(metrics.PrometheusMiddleware())

	router.POST("/pvz", auth.Middleware(), auth.RequireRole("moderator"), h.CreatePVZ)
	router.GET("/pvz", auth.Middleware(), auth.RequireRole("employee", "moderator"), h.GetPVZList)
	router.POST("/pvz/:pvzId/close_last_reception", auth.Middleware(), auth.RequireRole("employee"), h.CloseLastReception)
	router.POST("/pvz/:pvzId/delete_last_product", auth.Middleware(), auth.RequireRole("employee"), h.DeleteLastProductForPVZ)

	router.POST("/receptions", auth.Middleware(), auth.RequireRole("employee"), h.CreateReception)

	router.POST("/products", auth.Middleware(), auth.RequireRole("employee"), h.CreateProduct)

	router.POST("/dummyLogin", h.DummyLogin)
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented,
			domain.ErrorResponse{Message: domain.ErrMethodNotImplemented.Error()})
	})
	return router
}

func (h *Handler) CreatePVZ(c *gin.Context) {
	var pvz *domain.PVZ
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&pvz); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	pvz, err := h.service.PVZUsecase.CreatePVZ(pvz)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pvz)
	metrics.PvzCreated.Inc()
}

func (h *Handler) CreateReception(c *gin.Context) {
	var reception *domain.Reception
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&reception); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	reception, err := h.service.ReceptionUsecase.CreateReception(reception)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reception)
	metrics.ReceptionsCreated.Inc()
}

func (h *Handler) CloseLastReception(c *gin.Context) {
	pvzID := c.Param("pvzId")

	c.Header("Content-Type", "application/json")
	reception, err := h.service.ReceptionUsecase.CloseReception(&pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, reception)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var inputProduct *domain.InputProduct
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputProduct); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	product, err := h.service.ProductUsecase.CreateProduct(inputProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
	metrics.ProductsAdded.Inc()
}

func (h *Handler) DeleteLastProductForPVZ(c *gin.Context) {
	pvzID := c.Param("pvzId")

	err := h.service.ProductUsecase.DeleteLastProductForPVZ(&pvzID)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetPVZList(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	pvzList, err := h.service.PVZUsecase.GetPVZInfo(startDateStr, endDateStr, pageStr, limitStr)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, []domain.PVZWithReceptions{})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, pvzList)
}

func (h *Handler) DummyLogin(c *gin.Context) {
	var user domain.User
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	if auth.CheckRole(user.Role) == false {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidRole.Error()})
		return
	}

	token, err := auth.GenerateJWT("0", *user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *Handler) Register(c *gin.Context) {
	var inputUser domain.InputUser
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	if auth.CheckRole(inputUser.Role) == false {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: domain.ErrInvalidRole.Error()})
		return
	}

	user, err := h.service.UserUsecase.Register(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var inputUser domain.InputUser
	c.Header("Content-Type", "application/json")
	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: domain.ErrInvalidInputData.Error()})
		return
	}

	user, err := h.service.UserUsecase.Login(&inputUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := auth.GenerateJWT(*user.ID, *user.Role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}
