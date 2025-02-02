package quiz

type QuestionTimer struct {
	Question      Question `json:"question"`
	RemainingTime int      `json:"remaining_time"`
	IsAuto        bool     `json:"is_auto"`
}

// type UpdateTimerModeRequest struct {
// 	Mode TimerMode `json:"mode"`
// }
