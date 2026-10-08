package sweep

import (
	"hash/fnv"
	"log"
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
	//
	// It applies to test-prefixed names only. Resources with no name, such as
	// floating IPs, are never swept without a token; see ShouldSweepToken.
	EnvSweepUntokenized = "GCORE_SWEEP_UNTOKENIZED"

	// EnvSweepAllReservedFixedIPs allows sweeping external reserved fixed IPs
	// by state alone, with no ownership check at all.
	//
	// That family can carry no ownership marker: the name is computed by the
	// API, there is no tags field, and an external address sits on the region's
	// shared public network rather than a test one. Without this the sweeper
	// reclaims nothing, so a run killed mid-test leaks a billable address
	// permanently.
	//
	// Unlike EnvSweepUntokenized, which still narrows to test-prefixed names,
	// this narrows to nothing beyond project, region, IsExternal and DOWN. What
	// makes that tolerable is scope rather than precision: the list call is
	// project-scoped, so it can only reach a project dedicated to acceptance
	// tests. Do not copy the pattern into an account-global sweeper, where the
	// same reasoning does not hold.
	//
	// The residual cost is that it cannot tell a leaked address from one a
	// concurrent run still holds. CI runs are serialised by a shared
	// concurrency group, so the exposure is a local run overlapping a scheduled
	// sweep, and the result is a confusing test failure rather than lost
	// production capacity. Deletions name this variable in the log so that
	// failure is diagnosable rather than mysterious.
	EnvSweepAllReservedFixedIPs = "GCORE_SWEEP_ALL_RESERVED_FIXED_IPS"
)

// TagKey carries the run token for resources whose API exposes no name to put
// it in. Floating IPs are the case that forced this: they have no name field at
// all, so before tags the sweeper had nothing to identify them by.
const TagKey = "gcore_tf_acctest_run"

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

// ShouldSweepToken reports whether a resource may be deleted, given the run
// token it carries. ShouldSweep is the name-based entry point to this; call
// this one directly for resources with no name, where the token arrives some
// other way such as a tag.
//
// A resource with no token is never swept, even with EnvSweepUntokenized set.
// For a named resource the test prefix still bounds that opt-in, but a
// resource with no name has nothing else to go on: an untagged floating IP is
// indistinguishable from a reserved production one, so honouring the opt-in
// here would bring back the sweep that deleted every unattached address.
func ShouldSweepToken(resourceType, identifier, token string) bool {
	if token == "" {
		log.Printf("[INFO] Skipping %s: %s (carries no run token; resources with no name are never swept without one, audit them manually)", resourceType, identifier)
		recordSkip(resourceType)
		return false
	}

	want, ok := os.LookupEnv(EnvSweepRunID)
	if !ok || want == "" {
		log.Printf("[WARN] Skipping %s: %s (%s is unset; a sweep must name the finished run it reclaims)", resourceType, identifier, EnvSweepRunID)
		recordSkip(resourceType)
		return false
	}

	if wantToken := Token(want); token != wantToken {
		log.Printf("[DEBUG] Skipping %s: %s (token %q, sweeping %q)", resourceType, identifier, token, wantToken)
		recordSkip(resourceType)
		return false
	}

	log.Printf("[INFO] Sweeping %s: %s", resourceType, identifier)
	return true
}
