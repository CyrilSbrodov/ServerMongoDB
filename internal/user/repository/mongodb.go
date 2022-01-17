package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"server/internal/user"
)

type db struct {
	collection *mongo.Collection
}

func NewStorage(database *mongo.Database, collection string) user.Storage {
	return &db{
		collection: database.Collection(collection),
	}
}

func (d *db) Create(user *user.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		response := fmt.Sprintf("Пользователь %s создан.\n", user.Name)
		response += fmt.Sprintf("ID пользователя %s: %v", user.Name, oid.Hex())
		return response, nil
	}
	return "", fmt.Errorf("failed to convert objectid to hex. oid: %s", oid)
}

func (d *db) GetAll() (string, error) {
	var u []user.User
	result, err := d.collection.Find(context.Background(), bson.M{})
	if result.Err() != nil {
		return "", fmt.Errorf("failed to find all users. due to error: %v", err)
	}
	if err = result.All(context.Background(), &u); err != nil {
		return "", fmt.Errorf("failed to read all documents from db. due to error: %v", err)
	}
	response := ""
	for _, user := range u{
		response += user.ToString()
	}
	return response, nil

}

func (d *db) Get(id string) (string, error) {
	var u []user.User
	result, err := d.collection.Find(context.Background(), bson.M{})
	if result.Err() != nil {
		return "", fmt.Errorf("failed to find users due to error: %v", err)
	}
	if err = result.All(context.Background(), &u); err != nil {
		return "", fmt.Errorf("failed to read all documents from db, due to error: %v", err)
	}
	response := ""
	for _, user := range u {
		if user.Id == id {
			response += fmt.Sprintf("Друзья пользователя %s, ID пользователя %s \n", user.Name, user.Id)
			for i := 0; i < len(user.Friends); i++ {
				for _, j := range u {
					if user.Friends[i] == j.Id {
						response += fmt.Sprintf("Имя %s, возраст %d, ID пользователя %s, ID друзей %s \n", j.Name, j.Age, j.Id, j.Friends)
					}
				}
			}
		}
	}
	return response, nil
}

func (d *db) Update(id string, age int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("failed to convert user ID to ObjectID, ID: %s", id)
	}
	filter := bson.M{"_id": objectID}
	update := bson.D{
		{"$set", bson.D{{"age", age}}},
	}

	var u user.User
	result := d.collection.FindOneAndUpdate(ctx, filter, update)
	if err := result.Decode(&u); err != nil {
		return "", fmt.Errorf("failed to update user (id: %s). error: %v", id, err)
	}
	return fmt.Sprintf("Возраст пользователя %s изменен на %d.", u.Name, age), nil
}

func (d *db) Delete(u *user.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(u.Id)
	if err != nil {
		return "", fmt.Errorf("failed to convert user ID to ObjectID, ID: %s", u.Id)
	}
	filter := bson.M{"_id": objectID}

	result := d.collection.FindOneAndDelete(ctx, filter)
	if err = result.Decode(&u); err != nil {
		return "", fmt.Errorf("ошибка декодирования пользователя (id: %s) из базы данных. Ошибка: %v", u.Id, err)
	}

	filterSource := bson.M{"friends": objectID}
	updateUserObjSource := bson.D{
		{"$pull", bson.D{{"friends", objectID}}},
	}
	deleteFromFriends, err := d.collection.UpdateMany(context.Background(),filterSource,updateUserObjSource)
	if err != nil {
		return "", fmt.Errorf("failed to update users from db, ID: %s", u.Id)
	}
	response := ""
	if deleteFromFriends.MatchedCount == 1 {
		response += fmt.Sprintf("Пользователь удален из списка всех друзей\n")
	} else {
		response += fmt.Sprintf("У пользователя не было друзей\n")
	}
	response += fmt.Sprintf("Пользователь %s c ID %s удален", u.Name, u.Id)
	return response, nil
}

func (d *db) MakeFriends(id *user.ID) (string, error) {

	objectIDSource, err := primitive.ObjectIDFromHex(id.Source_id)
	if err != nil {
		return "", fmt.Errorf("failed to convert user ID to objectID. Due to error: %v", err)
	}
	objectIDTarget, err := primitive.ObjectIDFromHex(id.Target_id)
	if err != nil {
		return "", fmt.Errorf("failed to convert user ID to objectID. Due to error: %v", err)
	}

	filterSource := bson.M{"_id": objectIDSource}
	filterTarget := bson.M{"_id": objectIDTarget}

	var userSource user.User
	var userTarget user.User

	updateUserObjSource := bson.D{
		{"$addToSet", bson.D{{"friends", objectIDTarget}}},
	}
	updateUserObjTarget := bson.D{
		{"$addToSet", bson.D{{"friends", objectIDSource}}},
	}

	resultSource := d.collection.FindOneAndUpdate(context.Background(),filterSource, updateUserObjSource)
	if err := resultSource.Decode(&userSource); err != nil {
		return "", fmt.Errorf("failed to decode user from db. due to error: %v", err)
	}

	resultTarget := d.collection.FindOneAndUpdate(context.Background(),filterTarget, updateUserObjTarget)
	if err := resultTarget.Decode(&userTarget); err != nil {
		return "", fmt.Errorf("failed to decode user from db. due to error: %v", err)
	}
	response := fmt.Sprintf("Пользователь %s и пользователь %s теперь друзья.", userSource.Name, userTarget.Name)
	return response, nil
}

