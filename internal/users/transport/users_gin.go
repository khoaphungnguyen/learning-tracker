package usertransport

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khoaphungnguyen/learning-tracker/internal/auth"
	usermodel "github.com/khoaphungnguyen/learning-tracker/internal/users/model"
)

// LoginPayload login body
type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse token response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func (h *UserHandler) Signup(c *gin.Context) {
	var user usermodel.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid inputs",
		})
		c.Abort()
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{
			"Error": "Error Hashing Password",
		})
		c.Abort()
		return
	}
	err = h.userHandler.CreateUser(user.Email, user.Password, user.Salt, user.Name)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Creating User",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successful Register",
	})
}

// Login is a function that handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var payload LoginPayload
	//var user usermodel.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs",
		})
		c.Abort()
		return
	}
	user, err := h.userHandler.GetUserByEmail(payload.Email)
	if err != nil {
		c.JSON(401, gin.H{
			"Error": "Invalid Username",
		})
		c.Abort()
		return
	}
	err = user.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"Error": "Invalid User Credentials",
		})
		c.Abort()
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:         h.JWTKey,
		Issuer:            "AuthService",
		ExpirationMinutes: 30,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	signedRefreshToken, err := jwtWrapper.RefreshToken(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	token := LoginResponse{
		Token:        signedToken,
		RefreshToken: signedRefreshToken,
	}
	c.JSON(200, token)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetInt("id")
	user, err := h.userHandler.GetUser(userID)
	if err != nil {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}
	c.JSON(200, user)
}

// Handle Update User
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("id")
	var user usermodel.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid inputs",
		})
		c.Abort()
		return
	}
	err = h.userHandler.UpdateUser(userID, user.Email, user.Name)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Updating User",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successful Update",
	})
}

func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID := c.GetInt("id")

	err := h.userHandler.DeleteUser(userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Deleting User",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successful Delete",
	})
}

// Renew token from the refresh token
func (h *UserHandler) RenewAccessToken(c *gin.Context) {
	token, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs",
		})
		c.Abort()
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:         h.JWTKey,
		Issuer:            "AuthService",
		ExpirationMinutes: 30,
		ExpirationHours:   12,
	}
	claims, err := jwtWrapper.ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{
			"Error": "Invalid Token",
		})
		c.Abort()
		return
	}
	if claims.ExpiresAt < time.Now().Add(time.Minute*30).Unix() {
		c.JSON(401, gin.H{
			"Error": "Token is expired",
		})
		c.Abort()
		return
	}
	// convert id to int
	userID, err := strconv.Atoi(claims.Audience)
	if err != nil {
		c.JSON(401, gin.H{
			"Error": "Invalid Token",
		})
		c.Abort()
		return
	}
	signedToken, err := jwtWrapper.GenerateToken(userID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	token = signedToken
	c.JSON(200, gin.H{
		"token": token,
	})

}
