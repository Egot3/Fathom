package quiz

import "context"

type QuizRepository interface {
	RegisterQuiz(ctx context.Context, path string) error
	DeallocateQuiz(ctx context.Context, path string) error
	ListQuizzes(ctx context.Context, page, size int) ([]string, int, error)
	// fuzzy search goes to app layer, as sqlite is... for storage only
}
