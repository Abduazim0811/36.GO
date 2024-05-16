package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   int    `json:"year"`
}

var db *sql.DB

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "Abdu0811", "project")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	defer db.Close()
	name :="../internal/DB/createalbums.sql"
	filename,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	_,err=db.Exec(string(filename))
	if err!=nil{
		log.Fatal(err)
	}


	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", createAlbum)
	router.PUT("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)

	if err := router.Run(":7777"); err != nil {
		log.Fatal("Server run error: ", err)
	}
}

func getAlbums(c *gin.Context) {
	name:="../internal/DB/selectalbums.sql"
	filename,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	rows, err := db.Query(string(filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	albums := []Album{}
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Year); err != nil {
			log.Println(err)
			continue
		}
		albums = append(albums, album)
	}
	c.JSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	name:="../internal/DB/insertinto.sql"
	sqlfile,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	row := db.QueryRow(string(sqlfile), id)

	var album Album
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Year); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, album)
}

func createAlbum(c *gin.Context) {
	var album Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name:="../internal/DB/insertinto.sql"
	sqlfile,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	_, err = db.Exec(string(sqlfile),
		album.Title, album.Artist, album.Year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateAlbum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var album Album
	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name:="../internal/DB/deletealbums.sql"
	sqlfile,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	_, err = db.Exec(string(sqlfile),album.Title, album.Artist, album.Year, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteAlbum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	name:="../internal/DB/deletealbums.sql"
	sqlfile,err:=os.ReadFile(name)
	if err!=nil{
		log.Fatal(err)
	}
	_, err = db.Exec(string(sqlfile),id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
