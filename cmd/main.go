package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", createAlbum)
	router.PUT("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server run error: ", err)
	}
}

func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM albums")
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
	row := db.QueryRow("SELECT * FROM albums WHERE id = $1", id)

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

	_, err := db.Exec("INSERT INTO albums (title, artist, year) VALUES ($1, $2, $3)",
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

	_, err := db.Exec("UPDATE albums SET title=$1, artist=$2, year=$3 WHERE id=$4",
		album.Title, album.Artist, album.Year, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteAlbum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := db.Exec("DELETE FROM albums WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
