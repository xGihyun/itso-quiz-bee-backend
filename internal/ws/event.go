package ws

type Event string

const (
	QuizUpdateStatus Event = "quiz-update-status"

	QuizStart Event = "quiz-start"
	// QuizOpen  Event = "quiz-open"
	// QuizPause Event = "quiz-pause"
	// QuizClose Event = "quiz-close"

	QuizUpdateQuestion   Event = "quiz-update-question"
	QuizDisableAnswering Event = "quiz-disable-answering"

	TimerPass       Event = "timer-pass"
	TimerUpdateMode Event = "timer-update-mode"
	TimerDone       Event = "timer-done"

	PlayerJoin         Event = "player-join"
	PlayerLeave        Event = "player-leave"
	PlayerTypeAnswer   Event = "player-type-answer"
	PlayerSubmitAnswer Event = "player-submit-answer"

	Heartbeat Event = "heartbeat"
)
