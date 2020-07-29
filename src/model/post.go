package model

import (
	"context"
	"time"
	"tree-hole/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const colNamePost = "post"

var colPost *mongo.Collection

func initModelPost() {
	colPost = MongoDatabase.Collection(colNamePost)
}

type Comment struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Content   string             `bson:"content" json:"content"`
}

type Post struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Salt      string             `bson:"salt" json:"salt"`
	Title     string             `bson:"title" json:"title"`
	CreatedAt time.Time          `bson:"created_at" `
	UpdatedAt time.Time          `bson:"updated_at"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Content   string             `bson:"content" json:"content"`
	Reply     []Comment          `bson:"reply" json:"reply"`
}

func GetAllPost() ([]Post, error) {
	var result []Post
	cursor, err := colPost.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetPostWithID(idHex string) (Post, bool, error) {
	var post Post
	id, err := primitive.ObjectIDFromHex(idHex)
	err = colPost.FindOne(context.Background(), bson.M{"_id": id}).Decode(&post)
	if err == mongo.ErrNoDocuments {
		return post, false, nil
	}
	if err != nil {
		return post, false, err
	}
	return post, true, nil
}

func AddPost(post Post) (string, error) {
	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt
	post.Salt = util.RandomString(32)
	post.Reply = make([]Comment, 0)
	_, err := colPost.InsertOne(context.Background(), post)
	if err != nil {
		return "", err
	}
	return post.ID.Hex(), nil
}

func AddReplyWithPostID(idHex string, reply Comment) (string, error) {
	reply.ID = primitive.NewObjectID()
	reply.CreatedAt = time.Now()
	reply.UpdatedAt = reply.CreatedAt
	update := bson.M{
		"$set": bson.M{
			"updated_at": reply.CreatedAt,
		},
		"$push": bson.M{
			"reply": reply,
		},
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return "", err
	}
	_, err = colPost.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return "", err
	}
	return reply.ID.Hex(), nil
}
