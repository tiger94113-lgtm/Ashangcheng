package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Order 订单结构
type Order struct {
	ID             string    `json:"id"`
	OrderNo        string    `json:"orderNo"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Status         string    `json:"status"`
	Wallet         string    `json:"wallet"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	Note           string    `json:"note"`
	Referrer       string    `json:"referrer"`
	UsdtAmount     string    `json:"usdtAmount"`
	RewardPoolAmount string  `json:"rewardPoolAmount"`
	AssetAmount    string    `json:"assetAmount"`
	MinAssetAmount string    `json:"minAssetAmount"`
	TxHash         string    `json:"txHash"`
	ErrorMessage   string    `json:"errorMessage,omitempty"`
}

var db *sql.DB

func main() {
	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	initDB()
	defer db.Close()

	// 创建路由
	r := gin.Default()

	// CORS 配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// API 路由组
	api := r.Group("/api")
	{
		// 订单相关
		api.GET("/orders", getOrders)
		api.GET("/orders/:id", getOrder)
		api.POST("/orders", createOrder)
		api.PUT("/orders/:id", updateOrder)
		api.DELETE("/orders/:id", deleteOrder)
		
		// 统计
		api.GET("/stats", getStats)
		
		// 健康检查
		api.GET("/health", healthCheck)
	}

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// 初始化数据库
func initDB() {
	var err error
	
	// 从环境变量获取数据库路径，默认使用本地文件
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./orders.db"
	}
	
	// 确保目录存在
	if dir := filepath.Dir(dbPath); dir != "." {
		os.MkdirAll(dir, 0755)
	}
	
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// 创建订单表
	createTableSQL := `CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		order_no TEXT UNIQUE,
		created_at DATETIME,
		updated_at DATETIME,
		status TEXT,
		wallet TEXT,
		name TEXT,
		phone TEXT,
		address TEXT,
		note TEXT,
		referrer TEXT,
		usdt_amount TEXT,
		reward_pool_amount TEXT,
		asset_amount TEXT,
		min_asset_amount TEXT,
		tx_hash TEXT,
		error_message TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// 创建索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_wallet ON orders(wallet);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_status ON orders(status);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_created ON orders(created_at);")

	log.Println("Database initialized")
}

// 获取订单列表
func getOrders(c *gin.Context) {
	wallet := c.Query("wallet")
	status := c.Query("status")
	search := c.Query("search")

	query := "SELECT * FROM orders WHERE 1=1"
	args := []interface{}{}

	if wallet != "" {
		query += " AND wallet = ?"
		args = append(args, wallet)
	}

	if status != "" && status != "all" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if search != "" {
		query += " AND (order_no LIKE ? OR wallet LIKE ? OR name LIKE ? OR phone LIKE ? OR tx_hash LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID, &o.OrderNo, &o.CreatedAt, &o.UpdatedAt, &o.Status,
			&o.Wallet, &o.Name, &o.Phone, &o.Address, &o.Note,
			&o.Referrer, &o.UsdtAmount, &o.RewardPoolAmount, &o.AssetAmount,
			&o.MinAssetAmount, &o.TxHash, &o.ErrorMessage,
		)
		if err != nil {
			continue
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"count": len(orders),
	})
}

// 获取单个订单
func getOrder(c *gin.Context) {
	id := c.Param("id")

	var o Order
	err := db.QueryRow(
		"SELECT * FROM orders WHERE id = ?",
		id,
	).Scan(
		&o.ID, &o.OrderNo, &o.CreatedAt, &o.UpdatedAt, &o.Status,
		&o.Wallet, &o.Name, &o.Phone, &o.Address, &o.Note,
		&o.Referrer, &o.UsdtAmount, &o.RewardPoolAmount, &o.AssetAmount,
		&o.MinAssetAmount, &o.TxHash, &o.ErrorMessage,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": o})
}

// 创建订单
func createOrder(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置时间
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now

	// 如果没有 ID，生成一个
	if o.ID == "" {
		o.ID = generateID()
	}

	_, err := db.Exec(
		`INSERT INTO orders (id, order_no, created_at, updated_at, status, wallet, name, phone, address, note, 
		referrer, usdt_amount, reward_pool_amount, asset_amount, min_asset_amount, tx_hash, error_message) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.ID, o.OrderNo, o.CreatedAt, o.UpdatedAt, o.Status, o.Wallet, o.Name, o.Phone,
		o.Address, o.Note, o.Referrer, o.UsdtAmount, o.RewardPoolAmount, o.AssetAmount,
		o.MinAssetAmount, o.TxHash, o.ErrorMessage,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": o})
}

// 更新订单
func updateOrder(c *gin.Context) {
	id := c.Param("id")

	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o.UpdatedAt = time.Now()

	_, err := db.Exec(
		`UPDATE orders SET updated_at = ?, status = ?, name = ?, phone = ?, address = ?, 
		note = ?, tx_hash = ?, error_message = ? WHERE id = ?`,
		o.UpdatedAt, o.Status, o.Name, o.Phone, o.Address, o.Note,
		o.TxHash, o.ErrorMessage, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
}

// 删除订单
func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

// 获取统计
func getStats(c *gin.Context) {
	var totalCount, confirmedCount int
	var totalUsdt string

	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalCount)
	db.QueryRow("SELECT COUNT(*) FROM orders WHERE status = 'confirmed'").Scan(&confirmedCount)
	db.QueryRow("SELECT COALESCE(SUM(CAST(usdt_amount AS REAL)), 0) FROM orders").Scan(&totalUsdt)

	c.JSON(http.StatusOK, gin.H{
		"totalOrders":     totalCount,
		"confirmedOrders": confirmedCount,
		"totalUsdt":       totalUsdt,
	})
}

// 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}

// 生成唯一 ID
func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
