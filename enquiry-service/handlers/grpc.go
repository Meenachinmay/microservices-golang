package handlers

import (
	"context"
	"encoding/json"
	"enquiry-service/enquiries-grpc"
	"enquiry-service/internal/database"
	"enquiry-service/mqactions"
	"log"
	"time"
)

type EnquiryServer struct {
	enquiries.UnimplementedEnquiryServiceServer
	LocalApiConfig *LocalApiConfig
}

func (e *EnquiryServer) HandleCustomerEnquiry(ctx context.Context, request *enquiries.CustomerEnquiryRequest) (*enquiries.CustomerEnquiryResponse, error) {
	// starting counter
	var startTimer = time.Now()

	input := request.GetEnquiry()
	log.Printf("Processing customer enquiry:[DEBUG:HandleCustomerEnquiry]")

	// database insertion operation here
	log.Printf("inserting enquiry into database.%+v\n", input)
	// insert into database
	newEnquiry, err := e.LocalApiConfig.DB.CreateEnquiry(ctx, database.CreateEnquiryParams{
		UserID:     input.UserId,
		PropertyID: input.PropertyId,
	})
	if err != nil {
		//
		res := &enquiries.CustomerEnquiryResponse{
			Success: false,
			Message: err.Error() + "Could not insert customer enquiry:[HandlerCustomerEnquiry:GRPC]",
		}
		return res, err
	} else {
		log.Println("new enquiry created:[HandleCustomerEnquiry:GRPC]", newEnquiry)
	}

	// update the count of total enquiries made by current user
	updatedUser, err := e.LocalApiConfig.DB.AddNewEnquiryToUserById(ctx, input.UserId)
	if err != nil {
		res := &enquiries.CustomerEnquiryResponse{
			Success: false,
			Message: err.Error() + "Could not update user enquiry count for new customer enquiry:[HandlerCustomerEnquiry:GRPC]",
		}
		return res, err
	} else {
		log.Println("enquiry count updated for user", updatedUser)
	}

	// decide user communication
	totalEnquiries, err := e.getTotalEnquiriesLastWeek(ctx, updatedUser)
	if err != nil {
		res := &enquiries.CustomerEnquiryResponse{
			Success: false,
			Message: err.Error() + "Could not fetch user total enquiry count for new customer enquiry:[HandlerCustomerEnquiry:GRPC]",
		}
		return res, err
	} else {
		log.Printf("total enquiries count is %d for UserId %d", totalEnquiries, updatedUser.ID)
	}

	//
	mailPayloadForSendgrid := EnquiryMailPayloadUsingSendgrid{
		To:               "chinmayanand896@gmail.com",
		ToName:           "Chinmay anand",
		Subject:          "お問い合わせありがとうございます。",
		PropertyLocation: input.Location,
		PropertyName:     input.Name,
		Timestamp:        startTimer,
	}
	// execute communication
	err = e.notifyUserAboutEnquiry(updatedUser, totalEnquiries, mailPayloadForSendgrid)
	if err != nil {
		res := &enquiries.CustomerEnquiryResponse{
			Success: false,
			Message: err.Error() + "error in sending producing email to rabbitmq new customer enquiry:[HandlerCustomerEnquiry:GRPC]",
		}
		return res, err
	}
	log.Println("NotifyUserAboutEnquiry:[DEBUG_LOG]")

	//
	res := &enquiries.CustomerEnquiryResponse{
		Success: true,
		Message: "Successfully processed customer enquiry via gRPC.",
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

func (e *EnquiryServer) notifyUserAboutEnquiry(user database.User, totalEnquiries int, mailPayload EnquiryMailPayloadUsingSendgrid) error {
	if totalEnquiries >= 10 {
		log.Printf("Calling to the user %d...\n", user.ID)
		return nil
	} else if totalEnquiries >= 1 && totalEnquiries <= 3 {
		log.Printf("Sending sms to the user %d...\n", user.ID)
		return nil
	} else {
		log.Printf("Sending Email to the user %s...\n", user.Email)
		// communicate with mail-service to send an email using rabbitmq.
		err := e.sendEmail(mailPayload)
		return err
	}
	return nil
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
