package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-api/internal/core/session"
	"go-api/internal/core/user"
	"go-api/pkg/apierrors"
	"go-api/pkg/config"
	"go-api/pkg/utils"
)

type userHandler struct {
	cfg       *config.Config
	userUC    user.UseCase
	sessionUC session.UseCase
}

func NewUserHandler(cfg *config.Config, userUC user.UseCase, sessionUC session.UseCase) user.Handlers {
	return &userHandler{
		cfg:       cfg,
		userUC:    userUC,
		sessionUC: sessionUC,
	}
}

func (h *userHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		usr := &user.Model{}
		err := c.Bind(usr)
		if err != nil {
			c.JSON(apierrors.BadRequest().JSON())
			return
		}

		token, err := h.userUC.Register(c, usr)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		sess, err := h.sessionUC.CreateSession(c, token.User.ID)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		utils.CreateSessionCookie(h.cfg, c, sess)

		c.JSON(http.StatusCreated, token)
	}
}

func (h *userHandler) Login() gin.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	return func(c *gin.Context) {
		login := &Login{}
		err := c.Bind(login)
		if err != nil {
			c.JSON(apierrors.BadRequest().JSON())
			return
		}

		token, err := h.userUC.Login(c, login.Email, login.Password)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		sess, err := h.sessionUC.CreateSession(c, token.User.ID)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		utils.CreateSessionCookie(h.cfg, c, sess)

		c.JSON(http.StatusOK, token)
	}
}

func (h *userHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session-id")
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		if sessionID == "" {
			c.JSON(apierrors.Unauthorized().JSON())
			return
		}

		err = h.sessionUC.DeleteByID(c, sessionID)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		utils.DeleteSessionCookie(h.cfg, c, h.cfg.Session.Name)

		c.Status(http.StatusNoContent)
	}
}

func (h *userHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		user := &user.Model{}
		err = c.Bind(user)
		if err != nil {
			c.JSON(apierrors.BadRequest().JSON())
			return
		}
		user.ID = id

		updatedUser, err := h.userUC.Update(c, user)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *userHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		err = h.userUC.Delete(c, id)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		err = h.sessionUC.DeleteByID(c, id.String())
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		c.Writer.WriteHeader(http.StatusNoContent)
	}
}

func (h *userHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		user, err := h.userUC.GetByID(c, id)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (h *userHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		size, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			c.JSON(apierrors.BadRequest("invalid size").JSON())
			return
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			c.JSON(apierrors.BadRequest("invalid page").JSON())
			return
		}

		pagination := &utils.PaginationQuery{
			Size:    size,
			Page:    page,
			OrderBy: c.Query("order_by"),
		}

		users, err := h.userUC.GetUsers(c, pagination)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (h *userHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session-id")
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		session, err := h.sessionUC.GetSessionByID(c, sessionID)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		user, err := h.userUC.GetByID(c, session.UserID)
		if err != nil {
			c.JSON(apierrors.Parse(err).JSON())
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
