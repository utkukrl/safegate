package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type RuleAction string

const (
	ActionBlock RuleAction = "block"
	ActionLimit RuleAction = "limit"
	ActionAllow RuleAction = "allow"
)

type Rule struct {
	Method string
	Path   string
	Action RuleAction
	Rate   string
	Burst  int
}

type RuleSet struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make(map[string]Rule),
	}
}

func (r *RuleSet) AddRule(rule Rule) {
	key := r.key(rule.Method, rule.Path)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[key] = rule
	fmt.Printf("✅ Rule added: %+v\n", rule)
}

func (r *RuleSet) DeleteRule(method, path string) {
	key := r.key(method, path)
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.rules, key)
	fmt.Printf("❌ Rule removed: %s %s\n", method, path)
}

func (r *RuleSet) GetRules() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ruleList := make([]Rule, 0, len(r.rules))
	for _, r := range r.rules {
		ruleList = append(ruleList, r)
	}
	return ruleList
}

func (r *RuleSet) GetRule(method, path string) (Rule, bool) {
	key := r.key(method, path)
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[key]
	return rule, ok
}

func (r *RuleSet) key(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

func StartREPL(rules *RuleSet) {
	fmt.Println("🚦 Firewall REPL started. Type 'help' for commands.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("repl> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		handleCommand(line, rules)
	}
}

func handleCommand(input string, rules *RuleSet) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "block":
		if len(args) < 2 {
			fmt.Println("Usage: block <path>")
			return
		}
		rules.AddRule(Rule{
			Method: "*",
			Path:   args[1],
			Action: ActionBlock,
		})
	case "unblock":
		if len(args) < 2 {
			fmt.Println("Usage: unblock <path>")
			return
		}
		rules.DeleteRule("*", args[1])
	case "limit":
		if len(args) < 4 {
			fmt.Println("Usage: limit <METHOD> <path> rate=<rate> burst=<burst>")
			return
		}
		method := args[1]
		path := args[2]
		rate := ""
		burst := 0
		for _, part := range args[3:] {
			if strings.HasPrefix(part, "rate=") {
				rate = strings.TrimPrefix(part, "rate=")
			}
			if strings.HasPrefix(part, "burst=") {
				fmt.Sscanf(part, "burst=%d", &burst)
			}
		}
		rules.AddRule(Rule{
			Method: method,
			Path:   path,
			Action: ActionLimit,
			Rate:   rate,
			Burst:  burst,
		})
	case "show":
		if len(args) >= 2 && args[1] == "rules" {
			ruleList := rules.GetRules()
			if len(ruleList) == 0 {
				fmt.Println("📭 No active rules.")
				return
			}
			for _, rule := range ruleList {
				fmt.Printf("🔹 %+v\n", rule)
			}
		}
	case "help":
		fmt.Println(`Available commands:
  block <path>                     - Block all methods for path
  unblock <path>                   - Remove block rule for path
  limit <METHOD> <path> rate=R burst=B - Rate limit method/path
  show rules                       - List active rules
  help                             - Show this help`)
	default:
		fmt.Println("❓ Unknown command. Type 'help' for help.")
	}
}
