package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
	FocusedCollection *mongo.Collection
}

func (d *DB) SavePosts(posts []Post) (*mongo.InsertManyResult, error) {
	var postsToSave []any
	if len(posts) == 0 {
		return nil, nil
	}
	for _, post := range posts {
		createdAt := time.Unix(post.CreatedTimestamp/1000, (post.CreatedTimestamp%1000)*int64(time.Millisecond))

		doc := bson.M{
			"title":             post.Title,
			"author_id":         post.AuthorID,
			"subreddit_name":    post.SubredditName,
			"domain":            post.Domain,
			"id":                post.ID,
			"subreddit_id":      post.SubredditID,
			"url":               post.URL,
			"type":              post.Type,
			"number_comments":   post.NumberComments,
			"score":             post.Score,
			"created_timestamp": createdAt,
			"spoiler":           post.Spoiler,
			"promoted":          post.Promoted,
			"archived":          post.Archived,
			"nsfw":              post.NSFW,
		}

		postsToSave = append(postsToSave, doc)
	}
	insertManyResult, err := d.FocusedCollection.InsertMany(context.TODO(), postsToSave)
	if err != nil {
		return nil, err
	}
	return insertManyResult, err
}

func (d *DB) GetNewestPost() (*Post, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "created_timestamp", Value: -1}})

	result := d.FocusedCollection.FindOne(context.Background(), bson.D{}, opts)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return &Post{}, nil
		}
		return nil, fmt.Errorf("error finding newest post: %v", result.Err())
	}

	var post Post
	if err := result.Decode(&post); err != nil {
		return nil, fmt.Errorf("error decoding post: %v", err)
	}
	return &post, nil
}

func GetNewestPostTimestampOrDefault(post *Post) int64 {
	if post != nil {
		return post.CreatedTimestamp
	} else {
		return 0
	}

}
