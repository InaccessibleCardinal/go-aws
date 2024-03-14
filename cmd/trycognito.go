package main

import (
	"go-aws/internal/services/cognito"
	"go-aws/internal/types"
)

func TryCognito() {
	user := types.UserRecord{Email: "kennethlandsbaum@hotmail.com", UserName: "ken1"}
	password := "Password1!"
	cg := cognito.InitCognito()

	challengeResponses, err := cg.AuthenticateUser(user.UserName, password)
	if err != nil {
		panic(err)
	}
	cognito.SaveJson("cognito-challenge-responses.json", challengeResponses)
	sessionToken, err := cg.RespondToChallenge(challengeResponses)
	if err != nil {
		panic(err)
	}

	secrets := cg.AssociateToken(sessionToken)
	cognito.SaveJson("cognito-secrets.json", secrets)
}
