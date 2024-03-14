package cognito

import (
	"os"

	cognitosrp "github.com/alexrudd/cognito-srp/v2"
)

func GetCsrp(userName, password string) (CsrpIface, error) {
	clientId := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")
	cognitoUserPool := os.Getenv("COGNITO_USER_POOL")
	csrp, _ := cognitosrp.NewCognitoSRP(userName, password, cognitoUserPool, clientId, &clientSecret)
	return csrp, nil
}
