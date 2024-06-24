package handlers

import (
	"broker/gRPC-client/payment"
	"broker/helpers"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

func (lac *LocalApiConfig) PaymentViaGRPC(c *gin.Context) {
	var paymentPayload PaymentPayload

	err := helpers.ReadJSON(c, &paymentPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	conn, err := grpc.NewClient("payment-service:50002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	cc := payment.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	paymentResponse, err := cc.ProcessPayment(ctx, &payment.PaymentRequest{
		NewPayment: &payment.Payment{
			CardNumber:     paymentPayload.CardNumber,
			CardHolderName: paymentPayload.CardHolderName,
			CardCvv:        paymentPayload.CardCVV,
			CardExpiry:     paymentPayload.CardExpiry,
			Amount:         paymentPayload.Amount,
			Currency:       paymentPayload.Currency,
		},
	})

	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "payment processed via grpc" + paymentResponse.String()

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}
