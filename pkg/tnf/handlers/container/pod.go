package container

import (
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/testcases"
	"regexp"
	"strings"
	"time"
)

// Pod that is under test.
type Pod struct {
	result  int
	timeout time.Duration
	args    []string
	// Command, Is the command that is executed on the test subject
	Command string
	// Name, is the name of the test subject,e.g. CNF name
	Name string
	// Namespace, is the name of the namespace the CNF is deployed
	Namespace string
	// ExpectStatus, Is the is list of  expected output from the command
	ExpectStatus []string
	// Action, Specifies the action to ba taken on the return result
	Action testcases.TestAction
	// ResultType, Informs the type of result expected
	ResultType testcases.TestResultType
	// facts, stores the result to be used when collecting test output as facts for passing it on to the  next subject
	facts string
}

// Args returns the command line args for the test.
func (p *Pod) Args() []string {
	return p.args
}

// Timeout return the timeout for the test.
func (p *Pod) Timeout() time.Duration {
	return p.timeout
}

// Result returns the test result.
func (p *Pod) Result() int {
	return p.result
}

// ReelFirst returns a step which expects an pod status for the given pod.
func (p *Pod) ReelFirst() *reel.Step {
	return &reel.Step{
		Expect:  []string{testcases.GetOutRegExp(testcases.AllowAll)},
		Timeout: p.timeout,
	}
}

func contains(arr []string, str string) (found bool) {
	found = false
	for _, a := range arr {
		if a == str {
			found = true
			break
		}
	}
	return
}

// ReelMatch parses the status output and set the test result on match.
// Returns no step; the test is complete.
func (p *Pod) ReelMatch(_ string, _ string, match string) *reel.Step {
	// for type: array ,should match for any expected status or fail on any expected status
	// based on the action type allow (default)|deny
	p.facts = match
	if p.ResultType == testcases.ArrayType {
		re := regexp.MustCompile(testcases.GetOutRegExp(testcases.NullFalse)) // Single value matching null or false is considered postive
		matched := re.MatchString(match)
		if matched {
			p.result = tnf.SUCCESS
			return nil
		}
		replacer := strings.NewReplacer(`[`, ``, "\"", ``, `]`, ``, `, `, `,`)
		match = replacer.Replace(match)
		f := func(c rune) bool {
			return c == ','
		}
		matchSlice := strings.FieldsFunc(match, f)
		for _, status := range matchSlice {
			if contains(p.ExpectStatus, status) {
				if p.Action == testcases.Deny { // Single deny match is failure.
					return nil
				}
			} else if p.Action == testcases.Allow {
				return nil // should be in allowed list
			}
		}
	} else {
		for _, status := range p.ExpectStatus {
			re := regexp.MustCompile(testcases.GetOutRegExp(testcases.RegExType(status)))
			matched := re.MatchString(match)
			if !matched {
				return nil
			}
		}
	}

	p.result = tnf.SUCCESS
	return nil
}

// ReelTimeout does nothing;
func (p *Pod) ReelTimeout() *reel.Step {
	return nil
}

// ReelEOF does nothing.
func (p *Pod) ReelEOF() {
}

// Facts collects facts of the container
func (p *Pod) Facts() string {
	return p.facts
}

// NewPod creates a `Container` test  on the configured test cases.
func NewPod(args []string, name, namespace string, expectedStatus []string, resultType testcases.TestResultType, action testcases.TestAction, timeout time.Duration) *Pod {
	return &Pod{
		Name:         name,
		Namespace:    namespace,
		ExpectStatus: expectedStatus,
		Action:       action,
		ResultType:   resultType,
		result:       tnf.ERROR,
		timeout:      timeout,
		args:         args,
	}
}
