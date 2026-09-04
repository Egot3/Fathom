package quizparser_test

// NOOP for now
// func Test_parsers(t *testing.T) {
// 	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if d.IsDir() {
// 			t.Logf("Directory: %s\n", path)
// 		} else {
// 			t.Run(path, func(t *testing.T) {
// 				f, err := os.Open(path)
// 				require.NoError(t, err)
// 				defer f.Close()

// 				q, err := quizparser.ParseQuiz(f)
// 				if strings.Contains(path, "valid") {
// 					require.NoError(t, err)
// 					t.Logf("%+v", q)
// 				} else {
// 					t.Logf("path:%v\n%v", path, err)
// 					require.Error(t, err)
// 				}
// 			})

// 		}

// 		return nil
// 	})
// 	require.NoError(t, err)
// }
