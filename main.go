package main

import (
	"github.com/gin-gonic/gin"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	// "time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Entry is mongo structure
// type Entry struct {
// 	ID   bson.ObjectId `bson:"_id,omitempty"`
// 	Name string
// 	Time time.Time
// }

// Twitch is the json from top games
type Twitch struct {
	Total int                      `json:"_total"`
	Links map[string]interface{}   `json:"_links"`
	Top   []map[string]interface{} `json:"top"`
}

func main() {
	// lab := os.Getenv("MONGOLAB_URI")

	// session, err := mgo.Dial(lab)
	// col := session.DB("go-test").C("names")
	// if err != nil {
	// 	panic(err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/public", "public")

	r.GET("/", func(c *gin.Context) {
		// var results []Entry
		// col.Find(nil).All(&results)
		// c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
		// 	"title": results,
		// })
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	r.GET("/browse", func(c *gin.Context) {
		url := "https://api.twitch.tv/kraken/games/top"
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
		// games := dat["top"].([]interface{})
		fmt.Println(dat.Top)
		c.HTML(http.StatusOK, "browse.tmpl.html", gin.H{
			"games": dat.Top,
		})
	})

	r.GET("/browse/:game", func(c *gin.Context) {
		c.HTML(http.StatusOK, "show.tmpl.html", nil)
	})

	r.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search.tmpl.html", nil)
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl.html", nil)
	})

	// r.POST("/", func(c *gin.Context) {
	// 	name := c.PostForm("name")
	// 	err = col.Insert(&Entry{Name: name, Time: time.Now()})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	c.Redirect(http.StatusMovedPermanently, "/")
	// })

	r.Run(":" + port)
}
