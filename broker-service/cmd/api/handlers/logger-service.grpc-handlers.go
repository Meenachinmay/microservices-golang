package handlers

import (
	"broker/gRPC-client/logs"
	"broker/helpers"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

func (lac *LocalApiConfig) GetAllLogs(c *gin.Context) {
	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	cc := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &logs.GetAllLogsRequest{}
	res, err := cc.GetAllLogs(ctx, req)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "Fetched all logs"
	payload.Data = res.Logs

	helpers.WriteJSON(c, http.StatusAccepted, payload)

}
