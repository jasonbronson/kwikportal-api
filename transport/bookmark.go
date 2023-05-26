package transport

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwikportal-api/models"
	"github.com/jasonbronson/kwikportal-api/repositories"
	"golang.org/x/net/html"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileNameLength = 10
)

func getBookmarks(g *gin.Context) {

	userID := GetUserIDFromRequest(g)

	bookmarks, err := repositories.GetUsersBookmarks(userID)
	if err != nil {
		responseError(g, err)
		return
	}
	log.Println(bookmarks)
	responseData(g, bookmarks)
}

func uploadBookmarks(g *gin.Context) {
	// Get the uploaded file from the form data
	file, err := g.FormFile("bookmarkFile")
	if err != nil {
		responseError(g, fmt.Errorf("Failed to retrieve the uploaded file"))
		return
	}

	// Generate a random filename
	fileName := generateRandomFileName()

	// Save the uploaded file to /tmp
	if err := g.SaveUploadedFile(file, "/tmp/"+fileName); err != nil {
		responseError(g, fmt.Errorf("Failed to save uploaded file"))
		return
	}

	// Parse the uploaded file
	bookmarks, err := ParseBookmarks(g, "/tmp/"+fileName)
	if err != nil {
		responseError(g, fmt.Errorf("Failed to parse bookmarks"))
		return
	}

	log.Println("Attempting to save filename ", fileName)
	var bookmarksUnique []models.Bookmark
	for _, b := range bookmarks {
		if !bookmarkExist(b.URL, bookmarksUnique) {
			bookmarksUnique = append(bookmarksUnique, b)
		}
	}
	err = repositories.SaveAllBookmarks(bookmarksUnique)
	if err != nil {
		responseError(g, fmt.Errorf("Failed to save all bookmarks"))
		return
	}

	responseSuccess(g, "success", "ok")
}

func ParseBookmarks(g *gin.Context, filename string) ([]models.Bookmark, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	var parse func(*html.Node) []models.Bookmark
	parse = func(n *html.Node) []models.Bookmark {
		var bookmarks []models.Bookmark

		if n.Type == html.ElementNode && n.Data == "a" {
			bookmark := models.Bookmark{}

			bookmark.UserID = GetUserIDFromRequest(g)
			for _, attr := range n.Attr {
				switch attr.Key {
				case "href":
					var url string
					if attr.Val == "" {
						url = fmt.Sprintf("#%v", generateRandomFileName())
					} else {
						url = attr.Val
					}
					bookmark.URL = url
				case "add_date":
					addDate, err := strconv.ParseInt(attr.Val, 10, 64)
					if err != nil {
						return nil
					}
					bookmark.AddDate = addDate
				case "icon":
					bookmark.Icon = attr.Val
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					bookmark.Name = c.Data
					break
				}
			}

			bookmarks = append(bookmarks, bookmark)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			bookmarks = append(bookmarks, parse(c)...)
		}

		return bookmarks
	}

	return parse(doc), nil
}

func generateRandomFileName() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func bookmarkExist(url string, bookmarks []models.Bookmark) bool {
	for _, b := range bookmarks {
		if b.URL == url {
			return true
		}
	}
	return false
}
