package handlers

import (
	"auth-with-jwt/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"time"
)

type Claims struct {
    jwt.StandardClaims
    UserID uint
}

var jwtKey = []byte("shhhhh")

type AuthHandlers struct {
    db *gorm.DB
}

func NewAuthHandlers(db *gorm.DB) *AuthHandlers {
    return &AuthHandlers{db: db}
}

func (h *AuthHandlers) Signup(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
        return
    }
    user.Password = string(hashedPassword)

    if result := h.db.Create(&user); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (h *AuthHandlers) Login(c *gin.Context) {
    var loginParams struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&loginParams); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := h.db.Where("username = ?", loginParams.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginParams.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
        return
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString, "expires_at": expirationTime})
}