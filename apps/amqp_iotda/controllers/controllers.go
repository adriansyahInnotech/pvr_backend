package controllers

import (
	"apps/amqp_iotda/controllers/absence"
	"apps/amqp_iotda/controllers/enrollment"
	syncstudent "apps/amqp_iotda/controllers/sync_student"
	"apps/amqp_iotda/services"
	"pvr_backend/helper"
)

type Controllers struct {
	Enrollment  *enrollment.Enrollment
	Absence     *absence.Absence
	SyncStudent *syncstudent.SyncStudent
}

func NewControllers(helper *helper.Helper, allService *services.Services) *Controllers {

	return &Controllers{
		Enrollment:  enrollment.NewEnrollment(allService.Enrollment, helper),
		Absence:     absence.NewAbsence(allService.Absence, helper),
		SyncStudent: syncstudent.NewSyncStudent(allService.SyncStudendt, helper),
	}
}
