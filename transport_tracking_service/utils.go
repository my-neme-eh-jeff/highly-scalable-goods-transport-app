// utils.go

package main

import (
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func GetLocationsFromMongoDB(bookingID string) ([]LocationRecord, error) {
    collection := mongoClient.Database("transport").Collection("locations")
    cursor, err := collection.Find(ctx, bson.M{
        "booking_id": bookingID,
    }, options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}))
    if err != nil {
        log.Printf("Error retrieving locations from MongoDB: %v", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    var locations []LocationRecord
    for cursor.Next(ctx) {
        var result LocationRecord
        err := cursor.Decode(&result)
        if err != nil {
            log.Printf("Error decoding location record: %v", err)
            continue
        }
        locations = append(locations, result)
    }
    return locations, nil
}

func GetMongoChangeStream(bookingID string) (*mongo.ChangeStream, error) {
    collection := mongoClient.Database("transport").Collection("locations")

    pipeline := mongo.Pipeline{
        {{"$match", bson.D{
            {"fullDocument.booking_id", bookingID},
        }}},
    }

    opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
    changeStream, err := collection.Watch(ctx, pipeline, opts)
    if err != nil {
        log.Printf("Error creating change stream: %v", err)
        return nil, err
    }
    return changeStream, nil
}
