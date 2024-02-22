package process

import (
	"fmt"
	"net/http"
	"strconv"

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

	// r.GET("/dashboard", handleIndex)

	r.GET("/dashboard", func(c *gin.Context) {
		// Default offset to 0 if not specified
		offsetStr := c.DefaultQuery("offset", "0")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			// handle error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			return
		}

		// Use offset to fetch data
		cmp := getShareDetails(offset) // Implement this function to use offset
		if len(cmp) == 0 {
			c.HTML(http.StatusInternalServerError, "shares.html", gin.H{"error": "empty"})
			return
		}

		c.HTML(http.StatusOK, "shares.html", gin.H{"company": cmp})
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
	cmp := getShareDetailsWithoutOffset()
	if len(cmp) == 0 {
		c.HTML(http.StatusInternalServerError, "shares.html", gin.H{"error": "empty"})
		return
	}

	c.HTML(http.StatusOK, "shares.html", gin.H{"company": cmp})
}

func getShareDetailsWithoutOffset() []Company {
	query := fmt.Sprintf("SELECT vcInvestedCompany, dEntrySharePrice, dCurrentSharePrice FROM %s.%s ORDER BY iCompanyID DESC LIMIT 5 OFFSET ?", mysqlObj.Database, table)
	shares, flag := mysqlLoadShares(query, 0)
	if !flag {
		return []Company{}
	}
	return shares
}

func getShareDetails(offset int) []Company {
	query := fmt.Sprintf("SELECT vcInvestedCompany, dEntrySharePrice, dCurrentSharePrice FROM %s.%s ORDER BY iCompanyID DESC LIMIT 5 OFFSET ?", mysqlObj.Database, table)
	shares, flag := mysqlLoadShares(query, offset)
	if !flag {
		return []Company{}
	}
	return shares
}
