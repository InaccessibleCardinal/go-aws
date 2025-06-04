package server

import (
	"context"
	"go-aws/internal/types"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (d *Server) userDBPut(ctx context.Context, user types.User) error {
	user.UserID = d.idGenerator("USER#")
	marshaledUser, err := UserToAttributeValueMap(user)
	if err != nil {
		return err
	}
	_, err = d.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      marshaledUser,
		TableName: &d.userTable,
	})
	return err
}

func (d *Server) userDBGet(ctx context.Context, userID string) (*types.User, error) {
	log.Printf("getting user %s\n", userID)
	out, err := d.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key: IDToAttributeValue(userID), TableName: &d.userTable,
	})
	if err != nil {
		log.Printf("error getting user %s\n", err)
		return nil, err
	}
	if out.Item == nil {
		log.Println("no user found")
		return nil, nil
	}
	user, err := UserFromAttributeValueMap(out.Item)
	if err != nil {
		log.Printf("error parsing user %s\n", err)
		return nil, err
	}
	return user, nil
}

func (d *Server) userDBRemove(ctx context.Context, userId string) error {
	_, err := d.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: IDToAttributeValue(userId), TableName: &d.userTable})
	if err != nil {
		return err
	}
	return nil
}
