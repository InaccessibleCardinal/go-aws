package types

type User struct {
	UserID   string `json:"userId" dynamodbav:"USERID"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}
