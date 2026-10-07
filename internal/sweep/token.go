package sweep

import (
	"hash/fnv"
	"os"
	"strings"
)

const (
	// EnvRunID identifies the acceptance-test run that created a resource.
	// scripts/acctest generates a value when it is unset, so a local run is
	// tokenised without the developer doing anything.
	EnvRunID = "GCORE_ACCTEST_RUN_ID"

	// EnvSweepRunID names the run whose resources a sweep may delete. Without
	// it no tokenised resource is swept: deletion is only ever safe once the
	// owning run is known to have finished, and only the caller knows that.
	// EnvSweepUntokenized is the one path that does not consult it.
	EnvSweepRunID = "GCORE_SWEEP_RUN_ID"

	// EnvSweepUntokenized allows sweeping resources that carry no token, which
	// means everything created before tokens existed. Intended for a one-time
	// audit of a project, never for routine use.
	//
	// It works without EnvSweepRunID because an untokenised name records no run
	// to compare against, so requiring one would scope nothing. That also means
	// it cannot tell a leaked resource from a live one: only use it once no
	// branch or job still runs code that creates untokenised names.
	EnvSweepUntokenized = "GCORE_SWEEP_UNTOKENIZED"
)

// TokenLen is the fixed width of the run token inside a resource name. Names
// are the only channel the sweeper has, and several Gcore APIs cap their
// length, so the token is kept as short as collision resistance allows.
const TokenLen = 4

// tokenCharset is restricted to characters every Gcore resource name accepts.
const tokenCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

// Token reduces an arbitrary run identifier to a fixed-width token. It hashes
// rather than truncates so that a GitHub run ID, a UUID and a developer's ad
// hoc string all yield something name-safe, and so the test side and the
// sweeper derive the same token from the same run ID without coordinating.
func Token(runID string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(runID))
	v := h.Sum32()

	out := make([]byte, TokenLen)
	for i := range out {
		out[i] = tokenCharset[v%uint32(len(tokenCharset))]
		v /= uint32(len(tokenCharset))
	}
	return string(out)
}

// TokenFromName extracts the run token from a generated resource name.
//
// It reports false for any name that carries none, which includes every
// resource created before tokens existed. Those are never swept by token, only
// by the explicit untokenised opt-in.
func TokenFromName(name string) (string, bool) {
	for _, prefix := range []string{testPrefixHyphen, testPrefixUnderscore} {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		rest := name[len(prefix):]
		if len(rest) < TokenLen+1 {
			return "", false
		}
		if sep := rest[TokenLen]; sep != '-' && sep != '_' {
			return "", false
		}
		token := rest[:TokenLen]
		if strings.ContainsFunc(token, func(r rune) bool {
			return !strings.ContainsRune(tokenCharset, r)
		}) {
			return "", false
		}
		return token, true
	}
	return "", false
}

// RunToken returns the token for the current run, or "" when EnvRunID is unset.
func RunToken() string {
	if v, ok := os.LookupEnv(EnvRunID); ok && v != "" {
		return Token(v)
	}
	return ""
}
