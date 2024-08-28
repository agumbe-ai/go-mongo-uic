// Package versioned provides utilities for optimistic concurrency control
// in MongoDB documents, similar to the `update-if-current` plugin in Mongoose.
package versioned

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrVersionConflict is returned when a version conflict is detected.
var ErrVersionConflict = errors.New("version conflict occurred")

// SetInitialVersion ensures that the Version field is set to 1 if it's the first instance.
func SetInitialVersion(version *int) {
	if *version == 0 {
		*version = 1
	}
}

// UpdateIfCurrent updates a document if the current version matches the provided version.
// It atomically increments the version field upon a successful update.
func UpdateIfCurrent(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M, version int) (*mongo.SingleResult, error) {
	filter["version"] = version
	update["$inc"] = bson.M{"version": 1} // Increment the version atomically

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, ErrVersionConflict
		}
		return nil, result.Err()
	}

	return result, nil
}
