package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	MenuID   string `json:"menuID"`
	FoodName string `json:"foodName"`
	Price    int    `json:"price"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/public", "./public")

	// Serve frontend
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// ✅ Create Menu
	router.POST("/api/menu", func(ctx *gin.Context) {
		var req Menu
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
			return
		}

		fmt.Printf("Creating Menu: %+v\n", req)

		res := submitTxnFn("org1", "autochannel", "foodorder", "SmartContract", "invoke", nil,
			"CreateMenu", req.MenuID, req.FoodName, fmt.Sprintf("%d", req.Price))

		ctx.JSON(http.StatusOK, gin.H{"message": "Menu added", "result": res})
	})

	// ✅ Get All Menus
	router.GET("/api/menu/all", func(ctx *gin.Context) {
		res := submitTxnFn("org1", "autochannel", "foodorder", "SmartContract", "query", nil, "GetAllMenus")

		var menus []Menu
		if err := json.Unmarshal([]byte(res), &menus); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse menus"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": menus})
	})

	// ✅ Get Menu by ID
	router.GET("/api/menu/:id", func(ctx *gin.Context) {
		menuID := ctx.Param("id")
		res := submitTxnFn("org1", "autochannel", "foodorder", "SmartContract", "query", nil, "GetMenuByID", menuID)

		var menu Menu
		if err := json.Unmarshal([]byte(res), &menu); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": menu})
	})

	// ✅ Search Menu by Food Name
	router.GET("/api/menu/search", func(ctx *gin.Context) {
		name := ctx.Query("name")
		res := submitTxnFn("org1", "autochannel", "foodorder", "SmartContract", "query", nil, "SearchMenuByFoodName", name)

		var menus []Menu
		if err := json.Unmarshal([]byte(res), &menus); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse results"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": menus})
	})

	router.Run(":3001")
}
