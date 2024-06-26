package handlers

import (
	"context"
	"encoding/json"
	"enquiry-service/grpc-proto-files"
	"enquiry-service/internal/database"
	"enquiry-service/mqactions"
	"fmt"
	"github.com/Meenachinmay/microservice-shared/utils"
	"log"
	"time"
)

type EnquiryServer struct {
	enquiries.UnimplementedEnquiryServiceServer
	LocalApiConfig *LocalApiConfig
}

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

func (e *EnquiryServer) HandleCustomerEnquiry(ctx context.Context, request *enquiries.CustomerEnquiryRequest) (*enquiries.CustomerEnquiryResponse, error) {
	// starting counter
	var startTimer = time.Now()

	input := request.GetEnquiry()

	// insert into database
	_, err := e.LocalApiConfig.DB.CreateEnquiry(ctx, database.CreateEnquiryParams{
		UserID:     input.UserId,
		PropertyID: input.PropertyId,
	})
	if err != nil {
		return nil, err
	}

	// update the count of total enquiries made by current user
	updatedUser, err := e.LocalApiConfig.DB.AddNewEnquiryToUserById(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	// decide user communication
	totalEnquiries, err := e.getTotalEnquiriesLastWeek(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	//
	mailPayloadForSendgrid := EnquiryMailPayloadUsingSendgrid{
		To:               input.Email,
		ToName:           "Chinmay anand",
		Subject:          "お問い合わせありがとうございます。",
		PropertyLocation: input.Location,
		PropertyName:     input.Name,
		Timestamp:        startTimer,
	}

	// execute communication
	err = e.notifyUserAboutEnquiry(input, totalEnquiries, mailPayloadForSendgrid)
	if err != nil {
		return nil, err
	}

	//
	res := &enquiries.CustomerEnquiryResponse{
		Success: true,
		Message: "We received your enquiry, please wait while we are contacting you back.",
	}
	return res, nil
}

func (e *EnquiryServer) getTotalEnquiriesLastWeek(c context.Context, updatedUser database.User) (int, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	var count int64
	count, err := e.LocalApiConfig.DB.CountEnquiriesForUserInLastWeek(c, database.CountEnquiriesForUserInLastWeekParams{
		UserID:      updatedUser.ID,
		EnquiryDate: oneWeekAgo,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (e *EnquiryServer) SendEmail(payload EnquiryMailPayloadUsingSendgrid) error {
	emitter, err := mqactions.NewEmitter(e.LocalApiConfig.Rabbit, "mail_exchange", "enquiry_mail")
	if err != nil {
		return err
	}

	j, _ := json.Marshal(&payload)
	err = emitter.Emit(string(j))
	if err != nil {
		return err
	}
	return nil
}

func (e *EnquiryServer) notifyUserAboutEnquiry(input *enquiries.CustomerEnquiry, totalEnquiries int, mailPayload EnquiryMailPayloadUsingSendgrid) error {

	// ------------------------------------------------------------------
	timeSlotStr := input.AvailableTimings
	if utils.CheckIfSlotIsInCurrentTimeWindow(timeSlotStr) {
		return e.executeTask(input, mailPayload)
	} else {
		return e.scheduleTask(input, mailPayload)
	}

	return nil
	// ------------------------------------------------------------------
}

func (e *EnquiryServer) executeTask(input *enquiries.CustomerEnquiry, mailPayload EnquiryMailPayloadUsingSendgrid) error {
	switch input.PreferredMethod {
	case "email":
		log.Printf("Sending Email to the user %d...\n", input.UserId)
		return e.SendEmail(mailPayload)
	case "phone":
		log.Printf("Calling to the user %d...\n", input.UserId)
		return nil
	case "sms":
		log.Printf("sending sms to the user %d...\n", input.UserId)
		return nil
	default:
		return fmt.Errorf("invalid preferred method: %s", input.PreferredMethod)
	}

	return nil
}

func (e *EnquiryServer) scheduleTask(input *enquiries.CustomerEnquiry, mailPayload EnquiryMailPayloadUsingSendgrid) error {
	taskDetails := map[string]interface{}{
		"user_id": input.UserId,
		"method":  input.PreferredMethod,
		"payload": mailPayload,
	}

	taskDetailsJSON, err := json.Marshal(taskDetails)
	if err != nil {
		return err
	}

	_, err = e.LocalApiConfig.DB.CreateSchedule(context.Background(), database.CreateScheduleParams{
		UserID:        input.UserId,
		TaskType:      input.PreferredMethod,
		TaskDetails:   taskDetailsJSON,
		ScheduledTime: input.AvailableTimings,
	})
	if err != nil {
		return err
	}

	log.Printf("Scheduled the task %s for user %d at %s.\n", input.PreferredMethod, input.UserId, input.AvailableTimings)
	return nil
}
