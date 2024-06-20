package handlers

import (
	"enquiry-service/helpers"
	"enquiry-service/internal/database"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type EnquiryPayload struct {
	UserID     int32  `json:"user_id"`
	PropertyID int32  `json:"property_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

type EnquiryMailPayload struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// HandleANewEnquiry To handle a new enquiry we need following data
/*
	start a counter to keep the track of time taken
	enquiryPayload, fetch user from database, fetch property from database, calculate priority

	NOTE: When a user makes an enquiry from the client side, we can send all the required details from the client itself
    and can save few database operations in this handler. But now for the practice purpose only we are fetching data here inside the handler.
*/
func (localApiConfig *LocalApiConfig) HandleANewEnquiry(c *gin.Context) {
	// starting counter
	var startTimer = time.Now()

	// read json
	var payload EnquiryPayload

	if err := helpers.ReadJSON(c, &payload); err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	// validate

	// insert into database
	newEnquiry, err := localApiConfig.DB.CreateEnquiry(c, database.CreateEnquiryParams{
		UserID:     payload.UserID,
		PropertyID: payload.PropertyID,
	})
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't save enquiry in DB"), http.StatusInternalServerError)
		return
	} else {
		log.Println("new enquiry created", newEnquiry)
	}

	// update the count of total enquiries made by current user
	updatedUser, err := localApiConfig.DB.AddNewEnquiryToUserById(c, payload.UserID)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't update enquiry count in DB for user"), http.StatusInternalServerError)
		return
	} else {
		log.Println("enquiry count updated for user", updatedUser)
	}

	// decide user communication
	totalEnquiries, err := localApiConfig.getTotalEnquiriesLastWeek(c, updatedUser)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't get total enquiries"), http.StatusInternalServerError)
		return
	}

	// fetch the property detail here - (only this database operation can be saved if we get data from client)
	foundProperty, err := localApiConfig.DB.GetAPropertyDetailsById(c, payload.PropertyID)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't get property details"), http.StatusInternalServerError)
		return
	}

	// prepare notifyPayload (for now it's same as EnquiryEmailPayload)
	_ = EnquiryMailPayload{
		From:      "chinmayanand896@gmail.com",
		To:        updatedUser.Email,
		Subject:   "Thank you for your enquiry.",
		Message:   fmt.Sprintf("Thank you for your enquiry about propery name %s at location %s", foundProperty.Name, foundProperty.Location),
		Timestamp: startTimer,
	}

	// execute communication
	localApiConfig.notifyUserAboutEnquiry(updatedUser, totalEnquiries)

	// send response to the user
}

func (localApiConfig *LocalApiConfig) getTotalEnquiriesLastWeek(c *gin.Context, updatedUser database.User) (int, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	var count int64
	count, err := localApiConfig.DB.CountEnquiriesForUserInLastWeek(c, database.CountEnquiriesForUserInLastWeekParams{
		UserID:      updatedUser.ID,
		EnquiryDate: oneWeekAgo,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (localApiConfig *LocalApiConfig) notifyUserAboutEnquiry(user database.User, totalEnquiries int) {
	if totalEnquiries >= 10 {
		log.Printf("Calling to the user %d...\n", user.ID)
	} else if totalEnquiries >= 1 && totalEnquiries <= 3 {
		log.Printf("Sending sms to the user %d...\n", user.ID)
	} else {
		log.Printf("Sending Email to the user %d...\n", user.Email)
		// communicate with mail-service to send an email using rabbitmq.
	}
}

func (localApiConfig *LocalApiConfig) notifyUserWithEmail(payload EnquiryPayload, userEmailAddress string, propertyDetails database.Property) {

}
