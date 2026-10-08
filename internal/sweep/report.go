package sweep

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

// Sweep outcomes are counted per sweeper so a run can be read rather than
// guessed at. Every one of these was previously invisible: delete failures were
// logged and swallowed, a permission error on a List looked like a clean sweep,
// and the process exited 0 either way.
type counters struct {
	listed       int
	matched      int
	deleted      int
	failed       int
	skipped      int
	firstFailure error
}

var (
	mu       sync.Mutex
	stats    = map[string]*counters{}
	apiSkips int
)

func count(resourceType string) *counters {
	c, ok := stats[resourceType]
	if !ok {
		c = &counters{}
		stats[resourceType] = c
	}
	return c
}

// RecordListed notes how many resources a sweeper saw. A sweep that reports no
// deletions is ambiguous until you know whether it looked at 12 resources or
// 4000, which is also how a single unpaginated page shows up here.
func RecordListed(resourceType string, n int) {
	mu.Lock()
	defer mu.Unlock()
	count(resourceType).listed += n
}

// RecordDeleted notes a successful delete.
func RecordDeleted(resourceType, name string) {
	mu.Lock()
	defer mu.Unlock()
	c := count(resourceType)
	c.matched++
	c.deleted++
}

// RecordFailure notes a resource that matched but could not be deleted. The
// sweep continues -- one stuck resource must not abandon the rest -- but the
// exit code will reflect it.
func RecordFailure(resourceType, name string, err error) {
	mu.Lock()
	defer mu.Unlock()
	c := count(resourceType)
	c.matched++
	c.failed++
	if c.firstFailure == nil {
		c.firstFailure = err
	}
	log.Printf("[ERROR] Failed to delete %s %s: %s", resourceType, name, err)
}

// recordSkip notes a resource left alone, whatever the reason.
func recordSkip(resourceType string) {
	mu.Lock()
	defer mu.Unlock()
	count(resourceType).skipped++
}

// recordAPISkip notes a call that could not be read for permission reasons.
// This is counted rather than ignored because an account you cannot see is
// indistinguishable from an account with nothing in it.
func recordAPISkip(status int) {
	mu.Lock()
	defer mu.Unlock()
	apiSkips++
	log.Printf("[WARN] sweep: API returned HTTP %d; treating as no visibility, not as empty", status)
}

// Finish emits one machine-readable line per sweeper, as that sweeper ends
// rather than at the end of the run. resource.TestMain can exit before any
// final aggregation, and the test binary's own timeout can panic past it.
func Finish(resourceType string) {
	mu.Lock()
	defer mu.Unlock()
	c := count(resourceType)
	fmt.Fprintf(os.Stdout,
		"SWEEP_RESULT sweeper=%s listed=%d matched=%d deleted=%d failed=%d skipped=%d\n",
		resourceType, c.listed, c.matched, c.deleted, c.failed, c.skipped)
}

// FinalExitCode turns the recorded counts into an exit status, printing a
// summary first.
//
// resource.TestMain already exits non-zero when a sweeper *returns* an error.
// The sweepers deliberately continue past individual delete failures, so "no
// sweeper returned an error" is not the same as "cleanup succeeded" -- this is
// what closes that gap.
func FinalExitCode() int {
	mu.Lock()
	defer mu.Unlock()

	names := make([]string, 0, len(stats))
	for name := range stats {
		names = append(names, name)
	}
	sort.Strings(names)

	var listed, matched, deleted, failed int
	for _, name := range names {
		c := stats[name]
		listed += c.listed
		matched += c.matched
		deleted += c.deleted
		failed += c.failed
	}

	fmt.Fprintf(os.Stdout,
		"SWEEP_TOTAL sweepers=%d listed=%d matched=%d deleted=%d failed=%d api_skips=%d\n",
		len(names), listed, matched, deleted, failed, apiSkips)

	switch {
	case failed > 0:
		fmt.Fprintf(os.Stderr, "%d resource(s) matched but could not be deleted\n", failed)
		return 1
	case apiSkips > 0:
		fmt.Fprintf(os.Stderr, "%d list call(s) were refused; the sweep could not see everything it was asked to\n", apiSkips)
		return 1
	}
	return 0
}
