package handlers

import (
	"broker/gRPC-client/users"
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

func (lac *LocalApiConfig) CreateNewUser(c *gin.Context, userPayload types.UserPayload) {
	log.Printf("entered in CreateNewUser method:GRPC %+v\n", userPayload)

	conn, err := grpc.NewClient("enquiry-service:50004", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	log.Println("created grpc connection for enquiry:[EnquiryViaGRPC]")

	cc := users.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newUserResponse, err := cc.CreateNewUser(ctx, &users.CreateUserRequest{
		User: &users.CreateUser{
			Email:            userPayload.Email,
			Name:             userPayload.Name,
			PreferredMethod:  userPayload.PreferredMethod,
			AvailableTimings: userPayload.AvailabelTimings,
		},
	})

	log.Println("new user created:[CreateNewUserGRPC]")
	if err != nil {
		errRes := helpers.ParseDatabaseError(err)
		helpers.ErrorJSON(c, errors.New(errRes))
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: newUserResponse.Message,
	}

	helpers.WriteJSON(c, http.StatusAccepted, payload)
	return

}
