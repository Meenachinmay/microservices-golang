package handlers

import (
	"context"
	"enquiry-service/grpc-proto-files/users"
	"enquiry-service/internal/database"
	"log"
	"time"
)

type UserServer struct {
	users.UnimplementedUserServiceServer
	LocalApiConfig *LocalApiConfig
}

func (u *UserServer) CreateNewUser(ctx context.Context, request *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	input := request.GetUser()

	newUser, err := u.LocalApiConfig.DB.CreateUser(ctx, database.CreateUserParams{
		Email:                  input.Email,
		Name:                   input.Name,
		AvailableTimings:       input.AvailableTimings,
		PreferredContactMethod: input.PreferredMethod,
	})
	if err != nil {
		return nil, err
	} else {
		log.Println("New user created:[HandlerNewUser:GRPC]", newUser)
	}

	// Converting database user to a protobuf user.
	newUserResponse := &users.User{
		Id:               newUser.ID,
		Email:            newUser.Email,
		Name:             newUser.Name,
		PreferredMethod:  newUser.PreferredContactMethod,
		AvailableTimings: newUser.AvailableTimings,
		CreatedAt:        newUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        newUser.UpdatedAt.Format(time.RFC3339),
	}

	res := &users.CreateUserResponse{
		Success: true,
		Message: "New user created:[HandlerNewUser:GRPC]",
		User:    newUserResponse,
	}

	return res, nil
}
