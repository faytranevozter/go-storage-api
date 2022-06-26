package mongorepo

import (
	"storage-api/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	db              *mongo.Database
	mediaCollection string
}

// NewMongoRepo ...
func NewMongoRepo(db *mongo.Database) domain.MongoRepository {
	return &mongoRepo{
		db:              db,
		mediaCollection: "xmedia",
	}
}
