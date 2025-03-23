package ws

// TODO: Delete this file

// import "github.com/xGihyun/itso-quiz-bee/internal/quiz"

// WARN: This will only work if there's only one quiz in progress.
// Otherwise, the existing timer state would be overwritten.
// var QuizTimer quiz.Timer

// NOTE:
// It's much better to have a separate timer per quiz
// var quizTimer map[string]Timer

func (c *client) handleQuestionTimer() {
	for {
		// select {
		// case <-QuizTimer.Ticker.C:
		// 	QuizTimer.Duration -= 1
		//
		// 	response := Response{
		// 		Event: TimerPass,
		// 		Data:  QuizTimer.Duration,
		// 	}
		//
		// 	c.Pool.Broadcast <- response
		//
		// 	if QuizTimer.Duration <= 0 {
		// 		response = Response{
		// 			Event: TimerDone,
		// 		}
		// 		c.Pool.Broadcast <- response
		//
		// 		return
		// 	}
		// }
	}
}
