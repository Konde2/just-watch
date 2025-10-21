package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Video struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	OwnerID   string             `bson:"owner_id" json:"owner_id"`
	MinIOKey  string             `bson:"minio_key" json:"minio_key"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type VideoRepository struct {
	collection *mongo.Collection
}

func NewVideoRepository(client *mongo.Client) *VideoRepository {
	return &VideoRepository{
		collection: client.Database("justwatch").Collection("videos"),
	}
}

func (r *VideoRepository) Create(video *Video) error {
	video.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.TODO(), video)
	return err
}

func (r *VideoRepository) UpdateStatus(id string, status string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *VideoRepository) FindByID(id string) (*Video, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	var video Video
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&video)
	return &video, err
}

func (r *VideoRepository) List() ([]*Video, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var videos []*Video
	err = cursor.All(context.TODO(), &videos)
	return videos, err
}
