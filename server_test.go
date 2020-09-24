package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func printErrors(t *testing.T, err error, res *http.Response) {
	if err != nil {
		t.Error("Fail", err)
	}
	if res == nil {
		t.Error("Server crashed", res)
	}
	t.Log(res)
}

var participant1 = Participant{
	Name:  "Sharath",
	Email: "tnfssc@gmail.com",
	RSVP:  "Yes",
}

var participant2 = Participant{
	Name:  "Chandra",
	Email: "sharath.chandra.the.great@gmail.com",
	RSVP:  "No",
}

var participant3 = Participant{
	Name:  "Sheripally",
	Email: "s.sharath.chandra@outlook.com",
	RSVP:  "Yes",
}

func TestScheduleAMeeting(t *testing.T) {
	var meetingDetails NewMeeting
	meetingDetails.Title = "Test title"
	meetingDetails.Participants = append(meetingDetails.Participants, participant1, participant2, participant3)
	meetingDetails.StartTime = "2020-10-24T00:54:44.104Z"
	meetingDetails.EndTime = "2020-10-25T00:54:44.104Z"
	bytesData, _ := json.Marshal(meetingDetails)
	res, err := http.Post("http://localhost:9090/meetings", "application/json", bytes.NewBuffer(bytesData))
	printErrors(t, err, res)
}

func TestGetMeetingByID(t *testing.T) {
	res, err := http.Get("http://localhost:9090/meeting/5f6b92351afb1e1603fa8999")
	printErrors(t, err, res)
}

func TestGetMeetingOfParticipant(t *testing.T) {
	res, err := http.Get("http://localhost:9090/meetings?participant=tnfssc@gmail.com")
	printErrors(t, err, res)
}

func TestGetMeetingsByTimeRange(t *testing.T) {
	res, err := http.Get("http://localhost:9090/meetings?start=2020-05-19T18:38:27.628Z&end=2020-09-27T18:38:27.628Z")
	printErrors(t, err, res)
}

func TestGetMeetingOfParticipantByTimeRangeAndPagination(t *testing.T) {
	res, err := http.Get("http://localhost:9090/meetings?participant=tnfssc@gmail.com&start=2020-05-19T18:38:27.628Z&end=2020-09-27T18:38:27.628Z&limit=2&offset=1")
	printErrors(t, err, res)
}
