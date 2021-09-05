package apis

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/pkbhowmick/url-lite/grpc/url_gen/urlgen"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/pkbhowmick/url-lite/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDbUrl     = "mongodb://localhost:27017"
	Database       = "url-lite"
	UserCollection = "user"
	gRPCServerAddr = "localhost:50001"
)

func getNewDevKey() string {
	return uuid.New().String()
}

//  function will return api-dev-key for successful registration
func getAPIDevKey(w http.ResponseWriter, r *http.Request) {
	var user model.User

	// decode request body into user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "can't parse json body", http.StatusBadRequest)
		return
	}

	// connect with mongodb database
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDbUrl))
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// retrieve the user collection from url-lite db
	collection := client.Database(Database).Collection(UserCollection)

	ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.D{{"email", user.Email}}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		user.APIDevKey = getNewDevKey()
		user.AvailableRequest = 10

		ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()
		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GenKey() (string, error) {
	conn, err := grpc.Dial(gRPCServerAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewKeyGenClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rep, err := c.GenerateKey(ctx, &empty.Empty{})
	if err != nil {
		return "", err
	}
	return rep.Key, nil
}

func Serve() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	r.Post("/api/get-api-key", getAPIDevKey)

	s := http.Server{
		Addr:              ":3000",
		Handler:           r,
	}
	log.Printf("API server is running at %v",s.Addr)
	log.Fatalln(s.ListenAndServe())
}
