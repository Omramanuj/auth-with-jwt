package handlers

import (
    "auth-with-jwt/models"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)


func (h *AuthHandlers) GetAllProducts(c *gin.Context) {
    var products []models.Product
    if result := h.db.Find(&products); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, products)
}
func (h *AuthHandlers) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if result := h.db.Create(&product); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *AuthHandlers) GetProduct(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var product models.Product
    if result := h.db.First(&product, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

func (h *AuthHandlers) UpdateProduct(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var product models.Product
    if result := h.db.First(&product, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    var updateData map[string]interface{}
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if result := h.db.Model(&product).Updates(updateData); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, product)
}

func (h *AuthHandlers) DeleteProduct(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    var product models.Product
    if result := h.db.Delete(&product, id); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
