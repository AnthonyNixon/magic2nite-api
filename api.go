package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newPodCode() string {
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	DB_USER := os.Getenv("M2NDBUSER")
	DB_PASS := os.Getenv("M2NDBPASS")
	DB_HOST := os.Getenv("DBHOST")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/magic2nite?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	type Pod struct {
		Id         int       `json:"id"`
		ShortCode  string    `json:"short_code"`
		MaxPlayers int       `json:"max_players"`
		MinPlayers int       `json:"min_players"`
		Location   string    `json:"location"`
		StartTime  time.Time `json:"start_time"`
		CutoffTime time.Time `json:"cutoff_time"`
		Private    bool      `json:"private"`
		Password   string    `json:"password"`
		Format     string    `json:"format"`
	}

	router := gin.Default()
	// Add API handlers here

	// GET a pod
	router.GET("/pod/:shortCode", func(c *gin.Context) {
		var pod Pod

		shortCode := c.Param("shortCode")
		row := db.QueryRow("select id, short_code, max_players, min_players, private, password, format, location, start_time, cutoff_time from pods where short_code = ?;", shortCode)
		err = row.Scan(&pod.Id, &pod.ShortCode, &pod.MaxPlayers, &pod.MinPlayers, &pod.Private, &pod.Password, &pod.Format, &pod.Location, &pod.StartTime, &pod.CutoffTime)
		if err != nil {
			// if no results, send null
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusNotFound, nil)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusOK, pod)
		}
	})

	// GET all pods
	router.GET("/pods", func(c *gin.Context) {
		var (
			pod  Pod
			pods []Pod
		)

		rows, err := db.Query("SELECT id, short_code, max_players, min_players, private, password, format, location, start_time, cutoff_time from pods;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&pod.Id, &pod.ShortCode, &pod.MaxPlayers, &pod.MinPlayers, &pod.Private, &pod.Password, &pod.Format, &pod.Location, &pod.StartTime, &pod.CutoffTime)
			pods = append(pods, pod)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, gin.H{
			"result": pods,
			"count":  len(pods),
		})
	})

	// POST new pod
	router.POST("/pod", func(c *gin.Context) {
		var pod Pod
		c.BindJSON(&pod)

		var shortCode = newPodCode()
		pod.ShortCode = shortCode

		stmt, err := db.Prepare("insert into pods (short_code, max_players, min_players, private, password, format, location, start_time, cutoff_time) values(?,?,?,?,?,?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(shortCode, pod.MaxPlayers, pod.MinPlayers, pod.Private, pod.Password, pod.Format, pod.Location, pod.StartTime, pod.CutoffTime)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Append strings
		defer stmt.Close()
		c.JSON(http.StatusOK, pod)
	})

	router.OPTIONS("/pod", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.OPTIONS("/pods", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, struct{}{})
	})

	router.Use(cors.Default())

	router.Run(":3000")
}