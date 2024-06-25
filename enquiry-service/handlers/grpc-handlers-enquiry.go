package handlers

import (
	"context"
	"encoding/json"
	"enquiry-service/grpc-proto-files"
	"enquiry-service/internal/database"
	"enquiry-service/mqactions"
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

func (e *EnquiryServer) sendEmail(payload EnquiryMailPayloadUsingSendgrid) error {
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
		if input.PreferredMethod == "phone" {
			log.Printf("Calling to the user %d...\n", input.UserId)
			return nil
		} else if input.PreferredMethod == "email" {
			log.Printf("Sending Email to the user %d...\n", input.UserId)
			// communicate with mail-service to send an email using rabbitmq.
			err := e.sendEmail(mailPayload)
			return err
		} else {
			log.Printf("Sending sms to the user %d...\n", input.UserId)
			return nil
		}
	} else {
		if input.PreferredMethod == "phone" {
			log.Printf("Scheduling Calling to the user %d for time slot %s\n", input.UserId, input.AvailableTimings)
			return nil
		} else if input.PreferredMethod == "email" {
			log.Printf("Scheduling Email to the user %d for time slot %s\n", input.UserId, input.AvailableTimings)
			return nil
		} else {
			log.Printf("Scheduling sms to the user %d for the time slot %s\n", input.UserId, input.AvailableTimings)
			return nil
		}
	}

	return nil
}
