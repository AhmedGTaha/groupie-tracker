package main

import "testing"

func TestContainsIgnoreCase(t *testing.T) {
    // Table of test cases
    cases := []struct {
        s     string // the string to search in
        sub   string // the substring to search for
        want  bool   // expected result
    }{
        {"hello", "ell", true},
        {"Hello", "ell", true},       // mixed case
        {"HELLO", "hello", true},     // full upper/lower
        {"hello", "x", false},
        {"hello", "", true},          // empty substring is contained
        {"", "a", false},             // empty string contains nothing except ""
        {"", "", true},               // empty contains empty
        {"queen", "queen", true},     // exact match
        {"queen", "Queen", true},     // case insensitive
        {"Scorpions", "pion", true},
        {"Scorpions", "scorp", true}, // start
    }

    for _, c := range cases {
        got := containsIgnoreCase(c.s, c.sub)
        if got != c.want {
            t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", c.s, c.sub, got, c.want)
        }
    }
}