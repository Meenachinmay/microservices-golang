package handlers

import (
	"broker/gRPC-client/enquiries"
	"broker/helpers"
	"context"
	"errors"
	"github.com/Meenachinmay/microservice-shared/types"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

func (lac *LocalApiConfig) EnquiryViaGRPC(c *gin.Context, enquiryPayload types.EnquiryPayload) {
	//var enquiryPayload types.EnquiryPayload

	//err := helpers.ReadJSON(c, &enquiryPayload)
	//if err != nil {
	//	helpers.ErrorJSON(c, err)
	//	return
	//}
	log.Printf("loaded json payload for enquiry:[EnquiryViaGRPC]%+v\n", enquiryPayload)

	conn, err := grpc.NewClient("enquiry-service:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	log.Println("created grpc connection for enquiry:[EnquiryViaGRPC]")

	cc := enquiries.NewEnquiryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	enquiryResponse, err := cc.HandleCustomerEnquiry(ctx, &enquiries.CustomerEnquiryRequest{
		Enquiry: &enquiries.CustomerEnquiry{
			UserId:           enquiryPayload.UserID,
			PropertyId:       enquiryPayload.PropertyID,
			Name:             enquiryPayload.PropertyName,
			Location:         enquiryPayload.PropertyLocation,
			PreferredMethod:  enquiryPayload.PreferredMethod,
			AvailableTimings: enquiryPayload.AvailabelTimings,
			FudousanId:       enquiryPayload.FudousanID,
		},
	})

	log.Println("made enquiry:[EnquiryViaGRPC]")

	if err != nil {
		helpers.ErrorJSON(c, errors.New("Error after making gRPC request from api-gateway to enquiry service"+err.Error()))
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = enquiryResponse.Message

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}
