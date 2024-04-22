package utils_test

import (
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncludesOnly_WithOnlyIncludedRunes(t *testing.T) {
	runes := []rune{'a', 'b', 'c'}
	s := "abc"

	assert.True(t, utils.IncludesOnly(s, runes))
}

func TestIncludesOnly_WithOnlyIncludedRunes(t *testing.T) {
	runes := []rune{'a', 'b', 'c'}
	s := "abc"

	assert.True(t, utils.IncludesOnly(s, runes))
}

// func TestIncludesOnly_WithNotIncludedRunes(t *testing.T) {
// 	runes := []rune{'a', 'b', 'c'}
// 	s := "abcd"
// 	if utils.IncludesOnly(s, runes) {
// 		t.Errorf("Expected IncludesOnly(%q, %q) to be false, got true", s, runes)
// 	}
// }
//
// func TestIncludesOnly_WithEmptyString(t *testing.T) {
// 	runes := []rune{'a', 'b', 'c'}
// 	s := ""
// 	if !IncludesOnly(s, runes) {
// 		t.Errorf("Expected IncludesOnly(%q, %q) to be true, got false", s, runes)
// 	}
// }
//
// func TestIncludesOnly_WithEmptyRunes(t *testing.T) {
// 	runes := []rune{}
// 	s := "abc"
// 	if IncludesOnly(s, runes) {
// 		t.Errorf("Expected IncludesOnly(%q, %q) to be false, got true", s, runes)
// 	}
// }
//
// func TestIncludesOnly_WithEmptyStringAndRunes(t *testing.T) {
// 	runes := []rune{}
// 	s := ""
// 	if !IncludesOnly(s, runes) {
// 		t.Errorf("Expected IncludesOnly(%q, %q) to be true, got false", s, runes)
// 	}
// }
