package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

type Skill struct { Name string `json:"name"`; Version string `json:"version"`; TimeoutSeconds int `json:"timeoutSeconds"`; Inputs []string `json:"inputs"` }
type Passport struct { Agent string `json:"agent"`; Protocol string `json:"protocol"`; Skills []Skill `json:"skills"` }
type Finding struct { Skill string `json:"skill"`; Code string `json:"code"`; Message string `json:"message"` }

type Report struct { Agent string `json:"agent"`; SkillCount int `json:"skillCount"`; Findings []Finding `json:"findings"` }

func evaluate(p Passport) []Finding {
	var out []Finding
	if p.Agent == "" { out = append(out, Finding{"", "missing-agent", "agent is required"}) }
	if p.Protocol != "a2a" { out = append(out, Finding{"", "unsupported-protocol", "protocol must be a2a"}) }
	seen := map[string]bool{}
	for _, s := range p.Skills {
		if s.Name == "" { out = append(out, Finding{"", "missing-skill-name", "skill name is required"}); continue }
		if seen[s.Name] { out = append(out, Finding{s.Name, "duplicate-skill", "skill is declared more than once"}) }; seen[s.Name] = true
		if s.Version == "" { out = append(out, Finding{s.Name, "missing-version", "skill version is required"}) }
		if s.TimeoutSeconds <= 0 { out = append(out, Finding{s.Name, "invalid-timeout", "timeoutSeconds must be positive"}) }
		if len(s.Inputs) == 0 { out = append(out, Finding{s.Name, "no-inputs", "skill must declare at least one input"}) }
	}
	return out
}

func main() {
	file := flag.String("manifest", "examples/passport.json", "passport manifest")
	jsonOut := flag.Bool("json", false, "print machine-readable findings")
	flag.Parse()
	b, err := os.ReadFile(*file); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }
	var p Passport; if err := json.Unmarshal(b, &p); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }
	findings := evaluate(p)
	if *jsonOut { enc := json.NewEncoder(os.Stdout); enc.SetIndent("", "  "); _ = enc.Encode(Report{Agent: p.Agent, SkillCount: len(p.Skills), Findings: findings}) } else if len(findings) == 0 { fmt.Printf("PASS %s: %d skills ready\n", p.Agent, len(p.Skills)) } else { for _, f := range findings { fmt.Printf("FAIL %-18s %s\n", f.Code, f.Message) } }
	if len(findings) > 0 { os.Exit(1) }
}

// Keep deterministic output useful to callers that compare reports.
func sortedSkills(p Passport) []string { names := make([]string, 0, len(p.Skills)); for _, s := range p.Skills { names = append(names, s.Name) }; sort.Strings(names); return names }
