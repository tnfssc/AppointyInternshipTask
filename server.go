package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbURI string = "mongodb+srv://admin:admin@cluster0.nqqfp.mongodb.net/main?retryWrites=true&w=majority"

// Participant is ...
type Participant struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	RSVP  string `json:"rsvp" bson:"rsvp"`
}

// Meeting is ...
type Meeting struct {
	ID                primitive.ObjectID  `json:"_id" bson:"_id"`
	Title             string              `json:"title" bson:"title"`
	Participants      []Participant       `json:"participants" bson:"participants"`
	StartTime         primitive.Timestamp `json:"startTime" bson:"startTime"`
	EndTime           primitive.Timestamp `json:"endTime" bson:"endTime"`
	CreationTimestamp primitive.Timestamp `json:"createdAt" bson:"createdAt"`
}

// NewMeeting is ...
type NewMeeting struct {
	Title             string              `json:"title" bson:"title"`
	Participants      []Participant       `json:"participants" bson:"participants"`
	StartTime         string              `json:"startTime" bson:"startTime"`
	EndTime           string              `json:"endTime" bson:"endTime"`
	CreationTimestamp primitive.Timestamp `json:"createdAt" bson:"createdAt"`
}

// NewMeetingToDB is ...
type NewMeetingToDB struct {
	Title             string              `json:"title" bson:"title"`
	Participants      []Participant       `json:"participants" bson:"participants"`
	StartTime         primitive.Timestamp `json:"startTime" bson:"startTime"`
	EndTime           primitive.Timestamp `json:"endTime" bson:"endTime"`
	CreationTimestamp primitive.Timestamp `json:"createdAt" bson:"createdAt"`
}

func scheduleNewMeeting(meetingDetails NewMeeting) (*mongo.InsertOneResult, int) {
	lock.Lock()
	defer lock.Unlock()
	var result *mongo.InsertOneResult
	startTime, err := time.Parse(time.RFC3339, meetingDetails.StartTime)
	endTime, err := time.Parse(time.RFC3339, meetingDetails.EndTime)
	if err != nil {
		fmt.Println("Invalid Date format", err)
		return result, 400
	}
	var meetingDetailsToDB = NewMeetingToDB{
		Title:             meetingDetails.Title,
		Participants:      meetingDetails.Participants,
		StartTime:         primitive.Timestamp{T: uint32(startTime.Unix())},
		EndTime:           primitive.Timestamp{T: uint32(endTime.Unix())},
		CreationTimestamp: meetingDetails.CreationTimestamp,
	}
	result, err = collection.InsertOne(ctx, meetingDetailsToDB)
	if err != nil {
		fmt.Println("Insert failed", err)
		return result, 500
	}
	return result, 200
}

func getMeetingCollision(participants []Participant, startTime string, endTime string) (bool, int) {
	startT, err := time.Parse(time.RFC3339, startTime)
	endT, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		fmt.Println("Invalid Date format", err)
		return false, 400
	}
	start := primitive.Timestamp{T: uint32(startT.Unix())}
	end := primitive.Timestamp{T: uint32(endT.Unix())}
	var emails []string
	for _, p := range participants {
		if p.RSVP == "Yes" {
			emails = append(emails, p.Email)
		}
	}
	response := collection.FindOne(
		context.TODO(),
		bson.D{
			{"participants", bson.D{
				{"$elemMatch", bson.D{
					{"rsvp", "Yes"},
					{"email", bson.D{
						{"$in", emails},
					}},
				}},
			}},
			{"$or", bson.A{
				bson.D{
					{"startTime", bson.D{
						{"$gte", start},
						{"$lt", end},
					}},
				},
				bson.D{
					{"endTime", bson.D{
						{"$gte", start},
						{"$lt", end},
					}},
				}},
			},
		},
	)
	if response == nil {
		return false, 500
	}
	var result Meeting
	err = response.Decode(&result)
	if err != nil {
		return false, 200
	}
	return true, 200
}

func getMeetingByTimeRange(startTime string, endTime string) ([]Meeting, int) {
	var result []Meeting
	startT, err := time.Parse(time.RFC3339, startTime)
	endT, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		fmt.Println("Invalid Date format", err)
		return result, 400
	}
	start := primitive.Timestamp{T: uint32(startT.Unix())}
	end := primitive.Timestamp{T: uint32(endT.Unix())}
	cur, err := collection.Find(
		context.TODO(),
		bson.D{{"startTime", bson.D{{"$gte", start}, {"$lt", end}}}, {"endTime", bson.D{{"$gte", start}, {"$lt", end}}}},
	)
	if err != nil {
		fmt.Println("Not found!", err)
		return result, 500
	}
	for cur.Next(ctx) {
		var res Meeting
		err := cur.Decode(&res)
		result = append(result, res)
		if err != nil {
			fmt.Println("Bad DB!", err)
			return result, 500
		}
	}
	return result, 200
}

func getMeetingByEmail(email string) ([]Meeting, int) {
	var result []Meeting
	cur, err := collection.Find(
		context.TODO(),
		bson.D{{"participants", bson.D{{"$elemMatch", bson.D{{"email", email}}}}}},
	)
	if err != nil {
		fmt.Println("Not found!", err)
		return result, 500
	}
	for cur.Next(ctx) {
		var res Meeting
		err := cur.Decode(&res)
		result = append(result, res)
		if err != nil {
			fmt.Println("Bad DB!", err)
			return result, 500
		}
	}
	return result, 200
}

func getMeetingByID(meetingID string) (Meeting, int) {
	var result Meeting
	docID, err := primitive.ObjectIDFromHex(meetingID)
	if err != nil {
		fmt.Println("Invalid ID", err)
		return result, 400
	}
	response := collection.FindOne(context.TODO(), bson.M{"_id": docID})
	if response == nil {
		fmt.Println("Not found", err)
		return result, 404
	}
	err = response.Decode(&result)
	if err != nil {
		fmt.Println("Bad data", err)
		return result, 500
	}
	return result, 200
}

func handleMeetings(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // Save meeting
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(400)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		var meetingDetails = NewMeeting{
			CreationTimestamp: primitive.Timestamp{T: uint32(time.Now().Unix())},
		}
		err := json.NewDecoder(r.Body).Decode(&meetingDetails)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println(err)
			return
		}
		colliding, possibleStatusCode := getMeetingCollision(meetingDetails.Participants, meetingDetails.StartTime, meetingDetails.EndTime)
		if possibleStatusCode == 200 {
			if colliding {
				w.WriteHeader(403)
				return
			}
			response, statusCode := scheduleNewMeeting(meetingDetails)
			w.WriteHeader(statusCode)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(possibleStatusCode)
	} else if r.Method == "GET" { // Get meeting by time range
		r.ParseForm()
		if len(r.Form["start"]) != 0 && len(r.Form["end"]) != 0 {
			response, statusCode := getMeetingByTimeRange(r.Form["start"][0], r.Form["end"][0])
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(response)
		} else if len(r.Form["participant"]) != 0 {
			response, statusCode := getMeetingByEmail(r.Form["participant"][0])
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(400)
		}
	}
}

func handleMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // Get meeting by ID
		id := r.URL.Path[len("/meeting/"):]
		response, statusCode := getMeetingByID(id)
		if statusCode != 200 {
			w.WriteHeader(statusCode)
		} else {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}
}

func initiateServer() {
	http.HandleFunc("/meetings/", handleMeetings)
	http.HandleFunc("/meeting/", handleMeeting)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var client *mongo.Client
var collection *mongo.Collection
var ctx = context.TODO()

var lock sync.Mutex

func connectToDB() {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("main").Collection("meetings")
}

func disconnectFromDB() {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func main() {
	connectToDB()
	defer disconnectFromDB()
	initiateServer()
}
