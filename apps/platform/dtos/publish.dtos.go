package dtos

// type Publish struct {
// 	DeviceID string `json:"device_id"`
// 	Message  string `json:"message"`
// }

// type PublishPayload struct {
// 	Cmd string `json:"cmd"`
// }

type PersonPublish struct {
	Cmd  string       `json:"cmd"`
	Data []InfoPerson `json:"Data"`
}

type InfoPerson struct {
	Nisn string `json:"nisn"`
	Name string `json:"name"`
	Idx  uint   `json:"idx"`
}

// type RestorePayload struct {
// 	PublishPayload
// 	Data InfoRestore `json:"Data"`
// }

// type InfoRestore struct {
// 	User   []UserRestore   `json:"user"`
// 	Record []RecordRestore `json:"record"`
// }

// type UserRestore struct {
// 	PersonID string `json:"person_id"`
// 	Name     string `json:"name"`
// 	Palm1    string `json:"palm_1"`
// 	Palm2    string `json:"palm_2"`
// }

// type RecordRestore struct {
// 	RecordID  string    `json:"record_id"`
// 	PersonID  string    `json:"person_id"`
// 	Timestamp time.Time `json:"timestamp"`
// }
