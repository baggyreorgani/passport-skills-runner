package main

import "testing"

func validPassport() Passport { return Passport{Agent: "weather-seller", Protocol: "a2a", Skills: []Skill{{Name: "forecast", Version: "1.0.0", TimeoutSeconds: 10, Inputs: []string{"location"}}}} }

func TestValidPassport(t *testing.T) { if got := evaluate(validPassport()); len(got) != 0 { t.Fatalf("unexpected findings: %#v", got) } }

func TestDuplicateAndInvalidSkill(t *testing.T) {
	p := validPassport(); p.Skills = append(p.Skills, Skill{Name: "forecast", TimeoutSeconds: 0})
	got := evaluate(p); if len(got) != 3 { t.Fatalf("got %d findings, want 3: %#v", len(got), got) }
}

func TestProtocolIsRequired(t *testing.T) { p := validPassport(); p.Protocol = "http"; if got := evaluate(p); len(got) != 1 || got[0].Code != "unsupported-protocol" { t.Fatalf("unexpected result: %#v", got) } }
