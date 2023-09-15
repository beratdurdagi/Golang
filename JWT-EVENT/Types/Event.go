package Types

import (
	"time"
)

// Type For Event
type Event struct {
	//////
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	User_id     int       `json:"user_id"`
	Date        time.Time `json:"date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
	Deleted_At  time.Time `json:"deleted_at"`
}

// constructor for Event
func CreateEvent(title, desc string, endTime, date time.Time) *Event {
	return &Event{
		Title:       title,
		Description: desc,
		Date:        date,
		EndTime:     endTime,
	}
}

// Type for creating Event req
type CreatedEventReq struct {
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	Date        time.Time `json:"date"`
	EndTime     time.Time `json:"end_time"`
}

// Type for update Event req
type UpdateEventReq struct {
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	Date        time.Time `json:"date"`
	EndTime     time.Time `json:"end_time"`
}

// Type for Delete Event req
type DeleteEventReq struct {
	Id      int `json:"id"`
	User_id int `json:"user_id"`
}

// Type for Response
type EventResponse struct {
	Message string `json:"msg"`
	Ok      bool   `json:"ok"`
}
