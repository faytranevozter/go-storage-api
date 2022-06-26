package mongorepo

import (
	"context"
	"storage-api/domain"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterMedia(options map[string]interface{}) (bson.M, *moptions.FindOptions) {
	query := bson.M{
		"deleted_at": bson.M{
			"$eq": nil,
		},
	}

	if id, ok := options["id"].(string); ok {
		objID, _ := primitive.ObjectIDFromHex(id)
		query["_id"] = objID
	}

	if id, ok := options["id"].(primitive.ObjectID); ok {
		query["_id"] = id
	}

	if provider, ok := options["provider"].(string); ok {
		query["provider"] = provider
	}

	if xtype, ok := options["type"].(string); ok {
		query["type"] = xtype
	}

	if q, ok := options["q"].(string); ok && q != "" {
		query["title"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: q,
				Options: "i",
			},
		}
	}

	// limit, offset & sort
	mongoOptions := moptions.Find()
	if offset, ok := options["offset"].(int64); ok {
		mongoOptions.SetSkip(offset)
	}
	if offset, ok := options["offset"].(int); ok {
		mongoOptions.SetSkip(int64(offset))
	}

	if limit, ok := options["limit"].(int64); ok {
		mongoOptions.SetLimit(limit)
	}
	if limit, ok := options["limit"].(int); ok {
		mongoOptions.SetLimit(int64(limit))
	}

	if sortBy, ok := options["sort"].(string); ok {
		sortDir, ok := options["dir"].(string)
		if !ok {
			sortDir = "asc"
		}

		sortQ := bson.M{}
		sortDirMongo := int(1)
		if strings.ToLower(sortDir) == "desc" {
			sortDirMongo = -1
		}
		sortQ[sortBy] = sortDirMongo
		mongoOptions.SetSort(sortQ)
	}

	return query, mongoOptions
}

func (r *mongoRepo) FetchMedia(ctx context.Context, options map[string]interface{}) (list []domain.MediaMongo, total int64, err error) {
	list = make([]domain.MediaMongo, 0)

	query, findOptions := generateQueryFilterMedia(options)
	cur, err := r.db.Collection(r.mediaCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("Error (FetchMedia) Query : ", err)
		return list, 0, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		row := domain.MediaMongo{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Error (FetchMedia) Decode : ", err)
		}
		list = append(list, row)
	}

	// stop if limit 1
	single := false
	if singleOpt, ok := options["single"].(bool); ok {
		single = singleOpt
		findOptions.SetLimit(1)
	}
	if findOptions.Limit != nil && *findOptions.Limit == 1 && single {
		return list, int64(len(list)), nil
	}

	total, err = r.db.Collection(r.mediaCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("Error (FetchMedia) Count : ", err)
		return list, 0, err
	}

	return list, total, nil
}

func (r *mongoRepo) InsertMedia(ctx context.Context, payload *domain.MediaMongo) error {
	payload.ID = primitive.NewObjectID()

	_, err := r.db.Collection(r.mediaCollection).InsertOne(ctx, *payload)
	if err != nil {
		logrus.Error("Error (InsertMedia) Query : ", err)
		return err
	}

	return nil
}

func (r *mongoRepo) UpdateMedia(ctx context.Context, payload *domain.MediaMongo) error {
	query := bson.M{
		"_id": payload.ID,
	}
	_, err := r.db.Collection(r.mediaCollection).UpdateOne(ctx, query, bson.D{
		{Key: "$set", Value: *payload},
	})
	if err != nil {
		logrus.Error("Error (UpdateMedia) Query : ", err)
		return err
	}

	return nil
}
