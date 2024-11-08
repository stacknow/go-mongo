package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// User struct represents a user in MongoDB
type User struct {
    ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name  string             `json:"name,omitempty" bson:"name,omitempty"`
    Email string             `json:"email,omitempty" bson:"email,omitempty"`
}

var client *mongo.Client

// Connect to MongoDB
func connectToMongoDB() *mongo.Client {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    return client
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
    collection := client.Database("go_mongo_db").Collection("users")
    var users []User
    cur, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cur.Close(context.Background())

    for cur.Next(context.Background()) {
        var user User
        if err := cur.Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
    collection := client.Database("go_mongo_db").Collection("users")
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    user.ID = primitive.NewObjectID()

    _, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func main() {
    client = connectToMongoDB()
    defer client.Disconnect(context.Background())

    router := mux.NewRouter()
    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")

    log.Println("Server is running on port 8000")
    log.Fatal(http.ListenAndServe(":8000", router))
}
