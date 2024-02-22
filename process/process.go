package process

import (
	"fmt"
	"net/http"

	"phase1/pkg/database"

	"github.com/gin-gonic/gin"
)

type Company struct {
	Name   string  `json:"Name"`
	IShare float64 `json:"IShare"`
	CShare float64 `json:"CShare"`
}

var offset = 0

var (
	dbName   = "accounts_database"
	host     = "sample-account-sarthak13102232-a003.a.aivencloud.com"
	password = "AVNS_JtiZAi9JGargoHQ2eOs"
	port     = 25885
	user     = "avnadmin"
	timeout  = "30s"
	idle     = 1
	open     = 1
	life     = 1
	retry    = 5
	delay    = 5

	table = "invested_company_share_price"
)

var mysqlObj *database.Object

func Run() {

	mysqlObj = &database.Object{}
	mysqlObj.LoadConfigurations(host, user, dbName, password, "", "", "", port, idle, open, life, retry, delay, false, false, false, timeout)
	mysqlObj.Connect()

	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Serve static files (including the chart image)
	r.Static("/static", "./static")

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "shares.html", gin.H{})
	})
	// New API endpoint for company shares data
	r.GET("/api/company-share", func(c *gin.Context) {
		cmp := getShareDetails()
		if len(cmp) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No company shares data found"})
			return
		}
		c.JSON(http.StatusOK, cmp)
	})

	r.GET("/next", func(c *gin.Context) {
		offset += 5
		// fetch the next batch of data from the database
		handleIndex(c)
	})

	r.GET("/prev", func(c *gin.Context) {
		offset -= 5
		if offset < 0 {
			offset = 0
		}
		// fetch the next batch of data from the database
		handleIndex(c)
	})
	return r
}

func handleIndex(c *gin.Context) {
	cmp := getShareDetails()
	if len(cmp) == 0 {
		c.HTML(http.StatusInternalServerError, "shares.html", gin.H{"error": "empty"})
		return
	}

	c.HTML(http.StatusOK, "shares.html", gin.H{})
}

func getShareDetails() []Company {
	query := fmt.Sprintf("SELECT vcInvestedCompany, dEntrySharePrice, dCurrentSharePrice FROM %s.%s ORDER BY dEntrySharePrice, dCurrentSharePrice LIMIT 5 OFFSET ?", mysqlObj.Database, table)
	shares, flag := mysqlLoadShares(query, offset)
	if !flag {
		return []Company{}
	}
	return shares
}
