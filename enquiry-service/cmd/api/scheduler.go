package main

import (
	"context"
	"encoding/json"
	"enquiry-service/handlers"
	"github.com/Meenachinmay/microservice-shared/utils"
	"log"
	"time"
)

type TaskDetails struct {
	UserID  int32           `json:"user_id"`
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}

func ProcessScheduledTasks(s *handlers.EnquiryServer) {
	if s == nil || s.LocalApiConfig == nil || s.LocalApiConfig.DB == nil {
		log.Println("EnquiryServer or its dependencies are not initialized properly")
		return
	}

	log.Println("EnquiryServer or its dependencies are initialized properly and process started.")

	for {
		// fetch
		tasks, err := s.LocalApiConfig.DB.GetDueSchedules(context.Background())
		if err != nil {
			log.Println("error fetching due schedules", err.Error())
			time.Sleep(1 * time.Minute)
			continue
		}

		// process tasks
		for _, task := range tasks {
			if utils.CheckIfSlotIsInCurrentTimeWindow(task.ScheduledTime) {
				switch task.TaskType {
				case "phone":
					log.Printf("Calling user %d...\n", task.UserID)
				case "sms":
					log.Printf("Sending sms to user %d...\n", task.UserID)
				case "email":
					var taskDetails TaskDetails
					err := json.Unmarshal(task.TaskDetails, &taskDetails)
					if err != nil {
						log.Println("error unmarshalling task details", err.Error())
						continue
					}
					log.Println("task details ", task.TaskDetails)

					var mailPayload handlers.EnquiryMailPayloadUsingSendgrid
					err = json.Unmarshal(taskDetails.Payload, &mailPayload)
					if err != nil {
						log.Println("error unmarshalling task details", err.Error())
						continue
					}
					log.Println("mail payload", mailPayload)

					err = s.SendEmail(mailPayload)
					if err != nil {
						log.Println("error sending scheduled email", err.Error())
						continue
					}
				}
				err := s.LocalApiConfig.DB.DeleteSchedule(context.Background(), task.ID)
				if err != nil {
					log.Println("error deleting scheduled task", err.Error())
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func isCurrentTimeInSlot(scheduleTime time.Time) bool {
	now := utils.ConvertToTokyoTime()
	return now.After(scheduleTime) || now.Equal(scheduleTime)
}
