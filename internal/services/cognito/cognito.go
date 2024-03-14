package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"go-aws/internal/types"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cog "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cogTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoProviderIface interface {
	AssociateSoftwareToken(
		context.Context,
		*cog.AssociateSoftwareTokenInput,
		...func(*cog.Options)) (*cog.AssociateSoftwareTokenOutput, error)
	ConfirmSignUp(context.Context, *cog.ConfirmSignUpInput, ...func(*cog.Options)) (*cog.ConfirmSignUpOutput, error)
	InitiateAuth(context.Context, *cog.InitiateAuthInput, ...func(*cog.Options)) (*cog.InitiateAuthOutput, error)
	RespondToAuthChallenge(
		context.Context,
		*cog.RespondToAuthChallengeInput,
		...func(*cog.Options)) (*cog.RespondToAuthChallengeOutput, error)
	SignUp(context.Context, *cog.SignUpInput, ...func(*cog.Options)) (*cog.SignUpOutput, error)
}

type CsrpIface interface {
	GetAuthParams() map[string]string
	PasswordVerifierChallenge(map[string]string, time.Time) (map[string]string, error)
}

type CognitoAuthTokens struct {
	AccessToken  string
	IdToken      string
	RefreshToken string
}

type CognitoSecret struct {
	SecretCode string
	Session    string
}

type CognitoService struct {
	client       CognitoProviderIface
	clientId     string
	clientSecret string
	getCsrp      func(string, string) (CsrpIface, error)
	csrp         CsrpIface
	ctx          context.Context
}

func (c *CognitoService) setCsrp(userName, password string) {
	if c.csrp == nil {
		csrp, err := c.getCsrp(userName, password)
		if err != nil {
			panic(err)
		}
		c.csrp = csrp
	}
}

func (c *CognitoService) RegisterUser(user types.UserRecord, password string) (bool, error) {
	out, err := c.client.SignUp(c.ctx, &cog.SignUpInput{
		ClientId:       aws.String(c.clientId),
		Password:       aws.String(password),
		Username:       aws.String(user.UserName),
		UserAttributes: []cogTypes.AttributeType{{Name: aws.String("email"), Value: aws.String(user.Email)}},
		SecretHash:     aws.String(computeSecretHash(c.clientSecret, user.UserName, c.clientId)),
	})
	if err != nil {
		return false, err
	}
	return out.UserConfirmed, nil
}

func (c *CognitoService) ConfirmUserSignup(confirmationCode, userName string) (bool, error) {
	_, err := c.client.ConfirmSignUp(c.ctx, &cog.ConfirmSignUpInput{
		ClientId:         &c.clientId,
		ConfirmationCode: &confirmationCode,
		Username:         &userName,
		SecretHash:       aws.String(computeSecretHash(c.clientSecret, userName, c.clientId)),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CognitoService) AuthenticateUser(userName, password string) (map[string]string, error) {
	c.setCsrp(userName, password)
	out, err := c.client.InitiateAuth(c.ctx, &cog.InitiateAuthInput{
		AuthFlow:       cogTypes.AuthFlowTypeUserSrpAuth,
		AuthParameters: c.csrp.GetAuthParams(),
		ClientId:       aws.String(c.clientId),
	})
	if err != nil {
		return nil, err
	}
	if out.ChallengeName == "PASSWORD_VERIFIER" {
		challengeResponses, err := c.csrp.PasswordVerifierChallenge(out.ChallengeParameters, time.Now())
		if err != nil {
			return nil, err
		}
		return challengeResponses, nil
	}
	return nil, errors.New("only password verifiction supported")
}

func (c *CognitoService) RespondToChallenge(challengeResponses map[string]string) (string, error) {
	out, err := c.client.RespondToAuthChallenge(c.ctx, &cog.RespondToAuthChallengeInput{
		ChallengeName:      cogTypes.ChallengeNameTypePasswordVerifier,
		ChallengeResponses: challengeResponses,
		ClientId:           &c.clientId,
	})
	if err != nil {
		return "", err
	}
	return *out.Session, nil
}

func (c *CognitoService) AssociateToken(sessionToken string) *CognitoSecret {
	out, err := c.client.AssociateSoftwareToken(c.ctx, &cog.AssociateSoftwareTokenInput{Session: &sessionToken})
	if err != nil {
		panic(err)
	}
	return &CognitoSecret{SecretCode: *out.SecretCode, Session: *out.Session}
}

func New(
	ctx context.Context,
	client CognitoProviderIface,
	getCsrp func(string, string) (CsrpIface, error)) *CognitoService {
	clientId := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")
	return &CognitoService{
		client:       client,
		clientId:     clientId,
		clientSecret: clientSecret,
		csrp:         nil,
		ctx:          ctx,
		getCsrp:      getCsrp,
	}
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func SaveJson(fileName string, data any) {
	bts, _ := json.Marshal(data)
	os.WriteFile(fileName, bts, 0777)
}
