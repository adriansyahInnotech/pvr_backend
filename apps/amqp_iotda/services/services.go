package services

import (
	"apps/amqp_iotda/services/absence"
	"apps/amqp_iotda/services/enrollment"
	syncstudent "apps/amqp_iotda/services/sync_student"
	"pvr_backend/helper"
	"pvr_backend/repository"
)

type Services struct {
	Enrollment   enrollment.Enrollment
	Absence      absence.Absence
	SyncStudendt syncstudent.SyncStudent
}

func NewServices(helper *helper.Helper, platformRepository repository.PlatformRepository) *Services {
	return &Services{
		Enrollment:   enrollment.NewEnrollment(helper, platformRepository),
		Absence:      absence.NewAbsence(helper, platformRepository),
		SyncStudendt: syncstudent.NewSyncStudent(helper, platformRepository),
	}
}
