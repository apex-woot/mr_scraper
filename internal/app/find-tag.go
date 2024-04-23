package app

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/html"
)

func FindTagByName(doc *html.Node, name string) *html.Node {
	var node *html.Node
	var findNode func(n *html.Node)
	findNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == name {
			node = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNode(c)
		}
	}
	findNode(doc)
	return node
}

type Post struct {
	Title            string `json:"title"`
	AuthorID         string `json:"author_id"`
	SubredditName    string `json:"subreddit_name"`
	Domain           string `json:"domain"`
	ID               string `json:"id"`
	SubredditID      string `json:"subreddit_id"`
	URL              string `json:"url"`
	Type             string `json:"type"`
	NumberComments   int    `json:"number_comments"`
	Score            int    `json:"score"`
	CreatedTimestamp int64  `json:"created_timestamp"`
	Spoiler          bool   `json:"spoiler"`
	Promoted         bool   `json:"promoted"`
	Archived         bool   `json:"archived"`
	NSFW             bool   `json:"nsfw"`
}

func (p *Post) UnmarshalBSON(data []byte) error {
	type postBSON struct {
		Title            string
		AuthorID         string
		SubredditName    string
		Domain           string
		ID               string
		SubredditID      string
		URL              string
		Type             string
		NumberComments   int
		Score            int
		CreatedTimestamp primitive.DateTime `bson:"created_timestamp"`
		Spoiler          bool
		Promoted         bool
		Archived         bool
		NSFW             bool
	}
	var raw postBSON
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Title = raw.Title
	p.AuthorID = raw.AuthorID
	p.SubredditName = raw.SubredditName
	p.Domain = raw.Domain
	p.ID = raw.ID
	p.SubredditID = raw.SubredditID
	p.URL = raw.URL
	p.Type = raw.Type
	p.NumberComments = raw.NumberComments
	p.Score = raw.Score
	p.Spoiler = raw.Spoiler
	p.Promoted = raw.Promoted
	p.Archived = raw.Archived
	p.NSFW = raw.NSFW
	p.CreatedTimestamp = raw.CreatedTimestamp.Time().UnixMilli()

	return nil
}

type PostJSON struct {
	Post Post `json:"post"`
}

func GetChildren(doc *html.Node, newestCheckedPostTimestamp int64) ([]Post, error) {

	var posts []Post
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "faceplate-tracker" && c.Parent.Data == "reddit-feed" {
			for _, attr := range c.Attr {
				if attr.Key == "data-faceplate-tracking-context" {
					var postJSON PostJSON
					if err := json.Unmarshal([]byte(attr.Val), &postJSON); err != nil {
						return nil, err
					}
					if postJSON.Post.CreatedTimestamp > newestCheckedPostTimestamp {
						posts = append(posts, postJSON.Post)
					}

				}
			}

		}
		GetChildren(c, newestCheckedPostTimestamp)
	}
	return posts, nil
}
