package quiz

import "context"

type QuizRepository interface {
	RegisterQuiz(ctx context.Context, path string, checksum []byte, score int) error
	DeallocateQuiz(ctx context.Context, path string) error
	ListQuizzes(ctx context.Context, page, size int) ([]string, int, error)
	CheckRegistered(ctx context.Context, path string) (bool, error)
	CheckIntegrity(ctx context.Context, path string, checksum []byte) (bool, error)
	// fuzzy search goes to app layer, as sqlite is... for storage only
}
