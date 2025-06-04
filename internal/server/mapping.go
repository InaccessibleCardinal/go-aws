package server

import (
	"go-aws/internal/types"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func UserToAttributeValueMap(user types.User) (map[string]dbTypes.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Println("error marshaling user", err)
		return av, err
	}
	av["ORGID"] = &dbTypes.AttributeValueMemberS{Value: os.Getenv("ORGID")}
	return av, nil
}

func UserFromAttributeValueMap(attrMap map[string]dbTypes.AttributeValue) (*types.User, error) {
	var user types.User
	err := attributevalue.UnmarshalMap(attrMap, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IDToAttributeValue(userId string) map[string]dbTypes.AttributeValue {
	return map[string]dbTypes.AttributeValue{
		"ORGID":  &dbTypes.AttributeValueMemberS{Value: "ORGID#01HRR1VZRMQDJZR4FAP4M2ZF9V"},
		"USERID": &dbTypes.AttributeValueMemberS{Value: userId}}
}
