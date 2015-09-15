package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	// "fmt"
)

// Stream is the stream mongo structure
type Stream struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string
	Game string
	URL  string
}

// Twitch is the json from top games
type Twitch struct {
	Total int                      `json:"_total"`
	Links map[string]interface{}   `json:"_links"`
	Top   []map[string]interface{} `json:"top"`
}

// Streams is the json from top streams
type Streams struct {
	Total   int                      `json:"_total"`
	Links   map[string]interface{}   `json:"_links"`
	Streams []map[string]interface{} `json:"streams"`
}

func main() {
	lab := os.Getenv("MONGOLAB_URI")
	db := os.Getenv("DATBASE_NAME")

	session, err := mgo.Dial(lab)
	col := session.DB(db).C("streams")
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/public", "public")

	r.GET("/", func(c *gin.Context) {
		var results []Stream
		col.Find(nil).All(&results)
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"streams": results,
		})
	})

	r.GET("/browse", func(c *gin.Context) {
		url := "https://api.twitch.tv/kraken/games/top?limit=100"
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		dat := &Twitch{}
		if err := json.Unmarshal([]byte(body), &dat); err != nil {
			panic(err)
		}
		c.HTML(http.StatusOK, "browse.tmpl.html", gin.H{
			"games": dat.Top,
		})
	})

	r.GET("/browse/:game", func(c *gin.Context) {
		game := c.Param("game")
		url := "https://api.twitch.tv/kraken/search/streams?q=" + strings.Replace(game, " ", "%20", -1) + "&limit=25"
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		dat := &Streams{}
		if err := json.Unmarshal([]byte(body), &dat); err != nil {
			panic(err)
		}
		c.HTML(http.StatusOK, "show.tmpl.html", gin.H{
			"streams": dat.Streams, "game": game,
		})
	})

	r.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search.tmpl.html", nil)
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl.html", nil)
	})

	r.POST("/addStream", func(c *gin.Context) {
		name := c.PostForm("name")
		game := c.PostForm("game")
		url := c.PostForm("url")
		err = col.Insert(&Stream{Name: name, Game: game, URL: url})
		if err != nil {
			panic(err)
		}
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.Run(":" + port)
}
