package handlers

import (
	"broker/gRPC-client/enquiries"
	"broker/helpers"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

func (lac *LocalApiConfig) EnquiryViaGRPC(c *gin.Context) {
	var enquiryPayload EnquiryPayload

	err := helpers.ReadJSON(c, &enquiryPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	conn, err := grpc.NewClient("enquiry-service:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	cc := enquiries.NewEnquiryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	enquiryResponse, err := cc.HandleCustomerEnquiry(ctx, &enquiries.CustomerEnquiryRequest{
		Enquiry: &enquiries.CustomerEnquiry{
			UserId:     enquiryPayload.UserID,
			PropertyId: enquiryPayload.PropertyID,
			Name:       enquiryPayload.Name,
			Location:   enquiryPayload.Location,
			FudousanId: enquiryPayload.FudousanID,
		},
	})

	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = enquiryResponse.String()

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}
