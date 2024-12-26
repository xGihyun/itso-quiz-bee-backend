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

	QuizStartTimer Event = "quiz-start-timer"
	QuizTimerPass  Event = "quiz-timer-pass"

	PlayerJoin         Event = "player-join"
	PlayerLeave        Event = "player-leave"
	PlayerTypeAnswer   Event = "player-type-answer"
	PlayerSubmitAnswer Event = "player-submit-answer"

	Heartbeat Event = "heartbeat"
)
