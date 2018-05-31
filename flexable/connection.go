package flexable

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/nemesisesq/flexable/account"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
)

type Exception struct {
	Message string
}

func (e Exception) Error() string {
	return e.Message
}

func VerifyJWT(t string) (account.User, error) {
	var user account.User
	var exp Exception
	var cognitoData account.CognitoData

	//TODO Actualy verify the token and sign the secret properly refernced in https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	//TODO Check the token claims to make sure that they are valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		mapstructure.Decode(claims, &cognitoData)
		//Rehydrate User

		user.Find(bson.M{"cognito_data.sub": cognitoData.Sub})
		//user.CognitoData = cognitoData
		//user.Upsert(bson.M{"user_info.sub": cognitoData.Sub})
		return user, nil
		//json.NewEncoder(w).Encode(user)
	} else {
		//json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
		exp = Exception{Message: "Invalid authorization token"}
		return user, exp
	}
}

func SocketServerConnections(server socketio.Server, namespace string) {
	//TODO Set a listener to authenticat this name space

	server.OnEvent(fmt.Sprintf("/%v", namespace), "AUTHORIZATION", func(s socketio.Conn, data interface{}) {
		if token, ok := data.(string); ok {

			user, err := VerifyJWT(token)

			if err != nil {
				s.Close()
				log.Error(err)
				return
			}

			ctx := s.Context().(context.Context)
			ctx = context.WithValue(ctx, "user", user)
			InitWatchers(s)



			s.SetContext(ctx)
			OpenShiftHandler(s, nil)

		}
	})

	server.OnConnect(fmt.Sprintf("/%v", namespace), func(s socketio.Conn) error {

		s.Emit("authenticate", "please authenticate my bebe ")

		log.Info("Connecting to ", namespace)
		//set context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		ctx = context.WithValue(ctx, "cancel", cancel)
		//ctx = SetMongoSession(ctx)

		s.SetContext(ctx)

		log.WithFields(log.Fields{
			"namespace": s.Namespace(),
			"ID":        s.ID(),
		}).Info("connected:")
		return nil
	})

	server.OnError(fmt.Sprintf("/%v", namespace), func(s socketio.Conn, e error) {
		//ctx := s.Context().(context.Context)
		//cancel := ctx.Value("cancel").(context.CancelFunc)
		//cancel()
		fmt.Println("meet error:", e)
		fmt.Println("everything cancelled", e)
	})
	server.OnDisconnect(fmt.Sprintf("/%v", namespace), func(s socketio.Conn, msg string) {
		//ctx := s.Context().(context.Context)

		//cancel := ctx.Value("cancel").(context.CancelFunc)
		//cancel()
		fmt.Println("closed and cancelled", msg)
	})
}
