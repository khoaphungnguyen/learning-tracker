package usertransport

import (
	"log"

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
	err = h.userHandler.CreateUser(user.Email, user.Password, user.Salt, user.Fullname)
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
		ExpirationMinutes: 1,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	signedRefreshToken, err := jwtWrapper.RefreshToken(user.Email)
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
	email := c.GetString("email")
	user, err := h.userHandler.GetProfileByEmail(email)
	if err != nil {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}
	c.JSON(200, user)
}


func (h *UserHandler) UpdateProfile(c *gin.Context) {
	
}

// // checkAuthorization checks if the user is authorized to perform the operation
// func checkAuthorization(r *http.Request) (int, bool) {
// 	userIDFromContext := r.Context().Value("userID").(string)
// 	queryValues := r.URL.Query()
// 	userIDParam := queryValues.Get("userID")
// 	if userIDFromContext != userIDParam {
// 		return 0, false
// 	}
// 	userID, err := strconv.Atoi(userIDParam)
// 	if err != nil {
// 		return 0, false
// 	}
// 	return userID, true

// }





// handleUsers manages CRUD operations for the User model
// func (h *NetHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
// 	// Check if the user is authorized to perform the operation
// 	userID, isAuthorized := checkAuthorization(r)
// 	if !isAuthorized {
// 		http.Error(w, "Invalid user", http.StatusUnauthorized)
// 		return
// 	}


// 	case http.MethodPut:
// 		// Assuming the user ID is passed in the URL for update
// 		var updatedUser models.User
// 		err := json.NewDecoder(r.Body).Decode(&updatedUser)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		err = h.netHandler.UpdateUser(userID, updatedUser.Username, updatedUser.FirstName, updatedUser.LastName)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(fmt.Sprintf("Updated User#%d successfully", userID)))
// 		return

// 	case http.MethodDelete:
// 		err := h.netHandler.DeleteUser(userID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusNoContent)
// 		return
// 	}
// }
