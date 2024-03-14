package types

type UserRecord struct {
	UserId   string `json:"userId" dynamodbav:"USERID"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}
