package routes

import (
	"apps/amqp_iotda/controllers"
	"apps/amqp_iotda/services"
	"fmt"
	"log"

	"pvr_backend/helper"
	"pvr_backend/helper/utils/amqp/dtos"
)

type Router struct {
	controllers controllers.Controllers
}

func NewRouter(helper *helper.Helper, allServices *services.Services) *Router {

	return &Router{
		controllers: *controllers.NewControllers(helper, allServices),
	}
}

func (s *Router) Routes(payload dtos.DeviceMessageReport) error {

	switch payload.NotifyData.Body.Content.Command {
	case "enrollment":

		log.Println("MASUK ROUTES ENROLLMENT")
		return s.controllers.Enrollment.Handle(payload)

	case "absence":

		log.Println("MASUK ROUTES Absence")
		return s.controllers.Absence.Handle(payload)

	case "sync_student":

		log.Println("MASUK ROUTES Sync Students")
		return s.controllers.SyncStudent.Handle(payload)
		//
	default:

		return fmt.Errorf("unknown command: %s", payload.NotifyData.Body.Content.Command)
	}

}
