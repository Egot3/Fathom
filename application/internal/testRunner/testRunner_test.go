package testrunner

import (
	"os"
	"testing"
	"time"

	"github.com/egot3/fathom/internal/testutils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func TestTestRunner_Start(t *testing.T) {
	t.Parallel()

	quizUUID, err := uuid.NewV7()
	require.NoError(t, err, "are we for real?")
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	i := do.New()
	do.Provide(i, NewTestRunner)

	t.Run("Valid start", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)
	})

	t.Run("Invalid start", func(t *testing.T) {
		t.Run("Orphans", func(t *testing.T) {
			testCases := []struct {
				desc       string
				quizPathes []string
				quizUUIDs  uuid.UUIDs
			}{
				{
					desc:       "No pathes",
					quizPathes: nil,
					quizUUIDs:  quizUUIDs,
				},
				{
					desc:       "No uuids",
					quizPathes: quizPathes,
					quizUUIDs:  nil,
				},
			}
			for _, tC := range testCases {
				t.Run(tC.desc, func(t *testing.T) {
					i := i.Scope(tC.desc)

					r := do.MustInvoke[TestRunner](i)

					err := r.Start(t.Context(), time.Minute, tC.quizPathes, tC.quizUUIDs, nil, testUUID)
					require.Error(t, err)
				})
			}

		})

		t.Run("Not found", func(t *testing.T) {
			i := do.New()
			do.Provide(i, NewTestRunner)

			r := do.MustInvoke[TestRunner](i)

			err = r.Start(t.Context(), 10*time.Second, []string{"/unknown.md"}, quizUUIDs, nil, testUUID)
			require.Error(t, err)
		})

		t.Run("Not absolute", func(t *testing.T) {
			i := do.New()
			do.Provide(i, NewTestRunner)

			r := do.MustInvoke[TestRunner](i)

			err = r.Start(t.Context(), 10*time.Second, []string{"./unknown.md"}, quizUUIDs, nil, testUUID)
			require.Error(t, err)
		})
	})
}

func TestTestRunner_Get(t *testing.T) {
	t.Parallel()

	quizUUID, err := uuid.NewV7()
	require.NoError(t, err, "are we for real?")
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	i := do.New()
	do.Provide(i, NewTestRunner)

	t.Run("Valid get", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		q, err := r.Get(quizUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)
	})

	t.Run("Not found", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		q, err := r.Get(uuid.Nil)
		require.Error(t, err)
		require.Nil(t, q)

		require.ErrorIs(t, err, ErrQuizNotCached)
	})
}

func TestTestRunner_Stop(t *testing.T) {
	t.Parallel()

	quizUUID, err := uuid.NewV7()
	require.NoError(t, err, "are we for real?")
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	i := do.New()
	do.Provide(i, NewTestRunner)

	t.Run("Valid stop", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		r.Stop() /* this doesn't even return an error
		because there is no "invalid stop"
		Server always have valid testRunner
		running multiple stops will do nothing */

		q, err := r.Get(quizUUID)
		require.Error(t, err)
		require.Nil(t, q)

		require.ErrorIs(t, err, ErrRunnerInactive)
	})
}

func TestTestRunner_Upsert(t *testing.T) {
	t.Parallel()

	quizUUID, err := uuid.NewV7()
	require.NoError(t, err, "are we for real?")
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	i := do.New()
	do.Provide(i, NewTestRunner)

	t.Run("Valid upsert", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		q, err := r.Get(quizUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)

		newQ := testutils.TestQuiz(t)
		defer os.Remove(newQ.Name())
		defer newQ.Close()
		newQUUID, err := uuid.NewV7()
		require.NoError(t, err)

		err = r.UpsertQuiz([]string{newQ.Name()}, uuid.UUIDs{newQUUID})
		require.NoError(t, err)

		q, err = r.Get(newQUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)
	})

	t.Run("Unequal", func(t *testing.T) {
		t.Run("No UUID", func(t *testing.T) {
			i := do.New()
			do.Provide(i, NewTestRunner)

			r := do.MustInvoke[TestRunner](i)

			err = r.Start(t.Context(), 10*time.Second, nil, nil, nil, testUUID)
			require.NoError(t, err)

			newQ := testutils.TestQuiz(t)
			defer os.Remove(newQ.Name())
			defer newQ.Close()

			err = r.UpsertQuiz([]string{newQ.Name()}, nil)
			require.Error(t, err)

			require.ErrorIs(t, err, ErrBadQuizzes)
		})

		t.Run("No path", func(t *testing.T) {
			i := do.New()
			do.Provide(i, NewTestRunner)

			r := do.MustInvoke[TestRunner](i)

			err = r.Start(t.Context(), 10*time.Second, nil, nil, nil, testUUID)
			require.NoError(t, err)

			newQUUID, err := uuid.NewV7()
			require.NoError(t, err)

			err = r.UpsertQuiz(nil, uuid.UUIDs{newQUUID})
			require.Error(t, err)

			require.ErrorIs(t, err, ErrBadQuizzes)
		})
	})

	t.Run("Empty", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, nil, nil, nil, testUUID)
		require.NoError(t, err)

		err = r.UpsertQuiz(nil, nil)
		require.NoError(t, err)
	})
}

func TestTestRunner_Remove(t *testing.T) {
	t.Parallel()
	quizUUID := uuid.Must(uuid.NewV7())
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	i := do.New()
	do.Provide(i, NewTestRunner)

	t.Run("Valid remove", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		q, err := r.Get(quizUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)

		err = r.RemoveQuiz(quizUUIDs)
		require.NoError(t, err)

		q, err = r.Get(quizUUID)
		require.Error(t, err)
		require.Nil(t, q)

		require.ErrorIs(t, err, ErrQuizNotCached)
	})

	t.Run("Empty", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, nil, nil, nil, testUUID)
		require.NoError(t, err)

		err = r.RemoveQuiz(nil)
		require.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		q, err := r.Get(quizUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)

		err = r.RemoveQuiz(uuid.UUIDs{uuid.Nil})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrQuizNotCached)

		q, err = r.Get(quizUUID)
		require.NoError(t, err)
		require.Equal(t, "there is a body!", q.Body)
		require.Equal(t, "quiz!", q.Title)
		require.Equal(t, "yeah!", q.Answer.Input.Input)
	})
}

func TestTestRunner_Deadline(t *testing.T) {
	quizUUID := uuid.Must(uuid.NewV7())
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	t.Run("Valid deadline", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		timeEnd := time.Now().Add(10 * time.Second)
		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		ti, err := r.Deadline()
		require.NoError(t, err)
		require.WithinDuration(t, timeEnd, *ti, 2*time.Second)
	})

	t.Run("Not active", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 0*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		ti, err := r.Deadline()
		require.Error(t, err)
		require.Nil(t, ti)

		require.ErrorIs(t, err, ErrRunnerInactive)
	})
}

func TestTestRunner_CurrentTestUUID(t *testing.T) {
	quizUUID := uuid.Must(uuid.NewV7())
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	t.Run("Valid deadline", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		ti := r.CurrentTestUUID()
		require.Equal(t, nil, testUUID, ti)
	})

	t.Run("Not active", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		ti := r.CurrentTestUUID()
		require.Equal(t, uuid.Nil, ti)
	})
}

func TestTestRunner_Extend(t *testing.T) {
	quizUUID := uuid.Must(uuid.NewV7())
	quizUUIDs := uuid.UUIDs{quizUUID}

	quizFile := testutils.TestQuiz(t)
	defer os.Remove(quizFile.Name())
	defer quizFile.Close()

	quizPathes := []string{quizFile.Name()}
	testUUID, err := uuid.NewV7()
	require.NoError(t, err)

	t.Run("Valid extension", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		timeEnd := time.Now().Add(10 * time.Second)
		err = r.Start(t.Context(), 10*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		err := r.ExtendTime(time.Hour)
		require.NoError(t, err)

		ti, err := r.Deadline()
		require.NoError(t, err)
		require.WithinDuration(t, timeEnd.Add(time.Hour), *ti, 10*time.Second)
	})

	t.Run("Not active", func(t *testing.T) {
		i := do.New()
		do.Provide(i, NewTestRunner)

		r := do.MustInvoke[TestRunner](i)

		err = r.Start(t.Context(), 0*time.Second, quizPathes, quizUUIDs, nil, testUUID)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		err := r.ExtendTime(time.Hour)
		require.Error(t, err)

		require.ErrorIs(t, err, ErrRunnerInactive)
	})
}

func TestTestRunner_PauseResume(t *testing.T) {
	t.Parallel()

	t.Run("Valid pause resumance", func(t *testing.T) {
		t.Parallel()
		t.Run("By waiting", func(t *testing.T) {
			t.Parallel()

			i := do.New()
			do.Provide(i, NewTestRunner)

			quizUUID := uuid.Must(uuid.NewV7())
			quizUUIDs := uuid.UUIDs{quizUUID}

			quizFile := testutils.TestQuiz(t)
			defer os.Remove(quizFile.Name())
			defer quizFile.Close()

			quizPathes := []string{quizFile.Name()}
			testUUID, err := uuid.NewV7()
			require.NoError(t, err)

			r := do.MustInvoke[TestRunner](i)

			err = r.Start(t.Context(), 10*time.Minute, quizPathes, quizUUIDs, nil, testUUID)
			require.NoError(t, err)

			err = r.Pause()
			require.NoError(t, err)

			o, err := r.Deadline()

			oldDeadline := *o

			require.NoError(t, err)
			require.WithinDuration(t, time.Now().Add(10*time.Minute), oldDeadline, time.Second)

			t.Log("Please, stand by. Get yourself some tea")
			time.Sleep(10 * time.Second)
			t.Log("Thanks for your patience!")

			n, err := r.Deadline()
			newNotch := *n
			require.NoError(t, err)
			require.Equal(t, oldDeadline, newNotch)

			err = r.Resume()
			require.NoError(t, err)

			n, err = r.Deadline()
			newNotch = *n

			require.NoError(t, err)
			require.WithinDuration(t, oldDeadline.Add(10*time.Second), newNotch, 2*time.Second)
		})

		t.Run("By extending in the meantime", func(t *testing.T) {
			t.Parallel()

			i := do.New()
			do.Provide(i, NewTestRunner)

			r := do.MustInvoke[TestRunner](i)

			quizUUID := uuid.Must(uuid.NewV7())
			quizUUIDs := uuid.UUIDs{quizUUID}

			quizFile := testutils.TestQuiz(t)
			defer os.Remove(quizFile.Name())
			defer quizFile.Close()

			quizPathes := []string{quizFile.Name()}
			testUUID, err := uuid.NewV7()
			require.NoError(t, err)

			err = r.Start(t.Context(), 10*time.Minute, quizPathes, quizUUIDs, nil, testUUID)
			require.NoError(t, err)

			err = r.Pause()
			require.NoError(t, err)

			o, err := r.Deadline()

			oldDeadline := *o

			require.NoError(t, err)
			require.WithinDuration(t, time.Now().Add(10*time.Minute), oldDeadline, time.Second)

			err = r.ExtendTime(10 * time.Second)

			n, err := r.Deadline()
			newNotch := *n
			require.NoError(t, err)
			require.NotEqual(t, oldDeadline, newNotch)
			require.Equal(t, oldDeadline.Add(10*time.Second), newNotch)

			err = r.Resume()
			require.NoError(t, err)

			n, err = r.Deadline()
			newNotch = *n

			require.NoError(t, err)
			require.WithinDuration(t, oldDeadline.Add(10*time.Second), newNotch, 2*time.Second)
		})
	})

	t.Run("Invalid pause+resume sequence", func(t *testing.T) {
		t.Parallel()

		t.Run("Pausing", func(t *testing.T) {
			t.Parallel()

			t.Run("Paused", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				quizUUID := uuid.Must(uuid.NewV7())
				quizUUIDs := uuid.UUIDs{quizUUID}

				quizFile := testutils.TestQuiz(t)
				defer os.Remove(quizFile.Name())
				defer quizFile.Close()

				quizPathes := []string{quizFile.Name()}
				testUUID, err := uuid.NewV7()
				require.NoError(t, err)

				err = r.Start(t.Context(), 10*time.Minute, quizPathes, quizUUIDs, nil, testUUID)
				require.NoError(t, err)

				err = r.Pause()
				require.NoError(t, err)

				err = r.Pause()
				require.Error(t, err)
				require.ErrorIs(t, err, ErrRunnerPaused)
			})
			t.Run("Inactive", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				err := r.Pause()
				require.Error(t, err)
				require.ErrorIs(t, err, ErrRunnerInactive)
			})
		})

		t.Run("Resuming", func(t *testing.T) {
			t.Parallel()

			t.Run("Running", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				quizUUID := uuid.Must(uuid.NewV7())
				quizUUIDs := uuid.UUIDs{quizUUID}

				quizFile := testutils.TestQuiz(t)
				defer os.Remove(quizFile.Name())
				defer quizFile.Close()

				quizPathes := []string{quizFile.Name()}
				testUUID, err := uuid.NewV7()
				require.NoError(t, err)

				err = r.Start(t.Context(), 10*time.Minute, quizPathes, quizUUIDs, nil, testUUID)
				require.NoError(t, err)

				err = r.Pause()
				require.NoError(t, err)

				err = r.Resume()
				require.NoError(t, err)

				err = r.Resume()
				require.Error(t, err)
				require.ErrorIs(t, err, ErrRunnerNotPaused)
			})

			t.Run("Inactive", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				err := r.Resume()
				require.Error(t, err)
				require.ErrorIs(t, err, ErrRunnerInactive)
			})
		})

		t.Run("Expired", func(t *testing.T) {
			t.Parallel()

			t.Run("Expiration during pause", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				quizUUID := uuid.Must(uuid.NewV7())
				quizUUIDs := uuid.UUIDs{quizUUID}

				quizFile := testutils.TestQuiz(t)
				defer os.Remove(quizFile.Name())
				defer quizFile.Close()

				quizPathes := []string{quizFile.Name()}
				testUUID, err := uuid.NewV7()
				require.NoError(t, err)

				err = r.Start(t.Context(), 2*time.Second, quizPathes, quizUUIDs, nil, testUUID)
				require.NoError(t, err)

				err = r.Pause()
				require.NoError(t, err)

				d, err := r.Deadline()
				require.NoError(t, err)
				deadline := *d

				t.Log("Please, stand by. Get yourself some tea")
				time.Sleep(5 * time.Second)
				t.Log("Thanks for your patience!")

				err = r.Resume()
				require.NoError(t, err)

				n, err := r.Deadline()
				require.NoError(t, err)
				newDeadline := *n

				require.WithinDuration(t, deadline, newDeadline, 6*time.Second)
			})
			t.Run("Manual expiration during pause", func(t *testing.T) {
				i := do.New()
				do.Provide(i, NewTestRunner)

				r := do.MustInvoke[TestRunner](i)

				quizUUID := uuid.Must(uuid.NewV7())
				quizUUIDs := uuid.UUIDs{quizUUID}

				quizFile := testutils.TestQuiz(t)
				defer os.Remove(quizFile.Name())
				defer quizFile.Close()

				quizPathes := []string{quizFile.Name()}
				testUUID, err := uuid.NewV7()
				require.NoError(t, err)

				err = r.Start(t.Context(), 2*time.Second, quizPathes, quizUUIDs, nil, testUUID)
				require.NoError(t, err)

				err = r.Pause()
				require.NoError(t, err)

				err = r.ExtendTime(-20 * time.Minute) //yes, VSCode, -20*time.Saturday
				require.NoError(t, err)

				err = r.Resume()
				require.Error(t, err)
				require.ErrorIs(t, err, ErrRunnerExpired)
			})
		})
	})
}
