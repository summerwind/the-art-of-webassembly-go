package internal

import (
	"log"
	"os"
	"path"
	"regexp"
)

var numbers = regexp.MustCompile(`(\d+)`)

// ScrubNumbers scrubs numbers from stdout because Go's example tests don't allow patterns.
// See https://github.com/golang/go/issues/18831
func ScrubNumbers(main func()) {
	tmp, err := os.MkdirTemp("", "ScrubNumbers")
	if err != nil {
		log.Fatal("TempDir failed: ", err)
	}

	stdoutPath := path.Join(tmp, "stdout.txt")
	stdout, err := os.Create(stdoutPath)
	if err != nil {
		log.Fatal(err)
	}

	// Save the old os.Stdout and revert it regardless of the outcome.
	oldStdout := os.Stdout
	os.Stdout = stdout
	revertStdout := func() { stdout.Close(); os.Stdout = oldStdout }
	defer revertStdout()

	// Run the main command.
	main()

	// Revert os.Stdout so that test output is visible on failure.
	revertStdout()

	// Replay the captured stdout, replacing any number with "NNN"
	capturedStdout, err := os.ReadFile(stdoutPath)
	if err != nil {
		log.Fatal(err)
	}
	replaced := numbers.ReplaceAllString(string(capturedStdout), `NNN`)
	os.Stdout.WriteString(replaced)
}
