package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		Id         string    `json:"id"`
		MaxPlayers int       `json:"max"`
		MinPlayers int       `json:"min"`
		Location   string    `json:"location"`
		StartTime  time.Time `json:"start"`
		CutoffTime time.Time `json:"cutoffTime"`
	}

	router := gin.Default()
	// Add API handlers here

	// GET a cronjob
	router.GET("/pod/:id", func(c *gin.Context) {
		var (
			pod    Pod
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, name, crondef, command, description, active, logtime from jobs where id = ?;", id)
		err = row.Scan(&cronjob.Id, &cronjob.Name, &cronjob.Cron_def, &cronjob.Command, &cronjob.Description, &cronjob.Active, &cronjob.Logtime)
		if err != nil {
			// if no results, send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": cronjob,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	// GET all cronjobs
	router.GET("/pods", func(c *gin.Context) {
		var (
			cronjob  Cronjob
			cronjobs []Cronjob
		)

		rows, err := db.Query("SELECT id, name, crondef, command, description, active, logtime from jobs;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&cronjob.Id, &cronjob.Name, &cronjob.Cron_def, &cronjob.Command, &cronjob.Description, &cronjob.Active, &cronjob.Logtime)
			cronjobs = append(cronjobs, cronjob)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.JSON(http.StatusOK, gin.H{
			"result": cronjobs,
			"count":  len(cronjobs),
		})
	})

	// POST new cronjob
	router.POST("/pod", func(c *gin.Context) {
		var cronjob Cronjob
		c.BindJSON(&cronjob)

		stmt, err := db.Prepare("insert into jobs (name, crondef, command, description, active) values(?,?,?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(cronjob.Name, cronjob.Cron_def, cronjob.Command, cronjob.Description, cronjob.Active)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Append strings
		defer stmt.Close()
		c.JSON(http.StatusOK, cronjob)
	})

	router.PUT("/pod", func(c *gin.Context) {
		var cronjob Cronjob
		c.BindJSON(&cronjob)

		stmt, err := db.Prepare("update jobs set name = ?, crondef = ?, command = ?, description = ?, active = ? where id = ?;")

		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(cronjob.Name, cronjob.Cron_def, cronjob.Command, cronjob.Description, cronjob.Active, cronjob.Id)
		if err != nil {
			fmt.Print(err.Error())
		}

		defer stmt.Close()
		c.JSON(http.StatusOK, cronjob)
	})

	router.DELETE("/pod/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("delete from jobs where id=?;")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted job with id: %s", id),
		})

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
