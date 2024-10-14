package parsemake

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
)

// Makefile and it's corresponded parts matches signature of checkmake/parser
type Makefile struct {
	FileName  string
	Rules     RuleList
	Variables VariableList
}

type Rule struct {
	Target       StringNoAlloc
	Dependencies []StringNoAlloc
	Body         []StringNoAlloc
	LineNumber   int
}

type RuleList []Rule

type Variable struct {
	Name            StringNoAlloc
	SimplyExpanded  bool
	Assignment      StringNoAlloc
	SpecialVariable bool
	LineNumber      int
}

type VariableList []Variable

// TODO delete below

var (
	ErrParse                    = errors.New("parse error")
	reFindRule                  = regexp.MustCompile("^([a-zA-Z]+):(.*)")
	reFindRuleBody              = regexp.MustCompile("^\t+(.*)")
	reFindSimpleVariable        = regexp.MustCompile("^([a-zA-Z]+) ?:=(.*)")
	reFindExpandedVariable      = regexp.MustCompile("^([a-zA-Z]+) ?=(.*)")
	reFindSpecialVariable       = regexp.MustCompile("^\\.([a-zA-Z_]+):(.*)")
	prefixComment          byte = '#'
	prefixSpecial          byte = '.'
)

func Parse(filepath string) (ret *Makefile, err error) {
	return ParseLog(filepath, slog.LevelInfo)
}

// ParseLog is the main function to parse a Makefile from a file path string to a
// Makefile struct. This function should be kept fairly small and ideally most
// of the heavy lifting will live in the specific parsing functions below that
// know how to deal with individual lines.
func ParseLog(filepath string, logLvl slog.Level) (ret *Makefile, err error) {
	var (
		scanner *MakefileScanner
		line    []byte
	)
	slog.SetLogLoggerLevel(logLvl)
	ret = new(Makefile)
	ret.FileName = filepath
	scanner, err = NewMakefileScanner(filepath)
	if err != nil {
		return ret, err
	}
	defer scanner.Close()
	scanner.Scan()
	for {
		line = scanner.Bytes()
		slog.Debug(fmt.Sprintf("value of bytes '%s'", line))

		if scanner.Finished {
			return
		}
		if len(line) == 0 {
			scanner.Scan()
			continue
		}
		switch line[0] {
		case prefixComment:
			scanner.Scan()
		case prefixSpecial:
			if matches := reFindSpecialVariable.FindSubmatch(line); matches != nil {
				specialVar := Variable{
					Name:            bytes.TrimSpace(matches[1]),
					Assignment:      bytes.TrimSpace(matches[2]),
					SpecialVariable: true,
					LineNumber:      scanner.LineNumber,
				}
				ret.Variables = append(ret.Variables, specialVar)
			}
			scanner.Scan()
		default:
			err = parseRuleOrVariable(scanner, ret, line)
			if err != nil {
				return ret, err
			}
		}
		if scanner.Finished == true {
			return
		}
	}
}

func parseRuleOrVariable(scanner *MakefileScanner, m *Makefile, line []byte) (err error) {
	var ruleOrVariable any
	ruleOrVariable, err = scanner.parseRuleOrVariable(line)
	if err != nil {
		return err
	}
	switch ruleOrVariable.(type) {
	case Rule:
		rule, found := ruleOrVariable.(Rule)
		if found != true {
			return ErrParse
		}
		m.Rules = append(m.Rules, rule)
	case Variable:
		variable, found := ruleOrVariable.(Variable)
		if found != true {
			return ErrParse
		}
		m.Variables = append(m.Variables, variable)
	}
	return nil
}

// parseRuleOrVariable gets the parsing scanner in a state where it resides on
// a line that could be a variable or a rule. The function parses the line and
// subsequent lines if there is a rule body to parse and returns an interface
// that is either a Variable or a Rule struct and leaves the scanner in a
// state where it resides on the first line after the content parsed into the
// returned struct. The parsing of line details is done via regexing for now
// since it seems ok as a first pass but will likely have to change later into
// a proper lexer/parser setup.
func (s *MakefileScanner) parseRuleOrVariable(line []byte) (ret any, err error) {

	if matches := reFindRule.FindSubmatch(line); matches != nil {
		// we found a rule so we need to advance the scanner to figure out if
		// there is a body
		beginLineNumber := s.LineNumber - 1
		s.Scan()
		bodyMatches := reFindRuleBody.FindSubmatch(s.Bytes())
		ruleBody := make(ArrStringNoAlloc, 0, 20)
		for bodyMatches != nil {
			ruleBody = append(ruleBody, bytes.TrimSpace(bodyMatches[1]))

			// done parsing the rule body line, advance the scanner and potentially
			// go into the next loop iteration
			s.Scan()
			bodyMatches = reFindRuleBody.FindSubmatch(s.Bytes())
		}
		// trim whitespace from all dependencies
		deps := bytes.Split(matches[2], spaceAsBytes)
		filteredDeps := make(ArrStringNoAlloc, 0, cap(deps))

		for idx := range deps {
			item := bytes.TrimSpace(deps[idx])
			if len(item) > 0 {
				filteredDeps = append(filteredDeps, item)
			}
		}
		ret = Rule{
			Target:       bytes.TrimSpace(matches[1]),
			Dependencies: filteredDeps,
			Body:         ruleBody,
			LineNumber:   beginLineNumber,
		}
	} else if matches = reFindSimpleVariable.FindSubmatch(line); matches != nil {
		ret = Variable{
			Name:           bytes.TrimSpace(matches[1]),
			Assignment:     bytes.TrimSpace(matches[2]),
			SimplyExpanded: true,
			LineNumber:     s.LineNumber}
		s.Scan()
	} else if matches = reFindExpandedVariable.FindSubmatch(line); matches != nil {
		ret = Variable{
			Name:           bytes.TrimSpace(matches[1]),
			Assignment:     bytes.TrimSpace(matches[2]),
			SimplyExpanded: false,
			LineNumber:     s.LineNumber}
		s.Scan()
	} else {
		slog.Debug(fmt.Sprintf("Unable to match line '%s' to a Rule or Variable", line))
		s.Scan()
	}

	return
}

// MakefileScanner is a wrapping struct around bufio.Scanner which provides
// extra functionality like the current line number
type MakefileScanner struct {
	Scanner    *bufio.Scanner
	LineNumber int
	FileHandle *os.File
	Finished   bool
}

// Scan is a thin wrapper around the bufio.Scanner Scan() function
func (s *MakefileScanner) Scan() bool {
	s.LineNumber++
	scanResult := s.Scanner.Scan()
	if scanResult == false && s.Scanner.Err() == nil {
		s.Finished = true
	}
	return scanResult
}

// Close closes all open handles the scanner has
func (s *MakefileScanner) Close() error {
	return s.FileHandle.Close()
}

// Bytes is a thin wrapper around bufio.Scanner Bytes()
func (s *MakefileScanner) Bytes() []byte {
	return s.Scanner.Bytes()
}

// NewMakefileScanner returns a MakefileScanner struct for parsing a Makefile
func NewMakefileScanner(filepath string) (ret *MakefileScanner, err error) {
	ret = new(MakefileScanner)
	ret.FileHandle, err = os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening the provided filepath '%s %w'", filepath, err)
	}
	ret.Scanner = bufio.NewScanner(ret.FileHandle)
	ret.Scanner.Split(bufio.ScanLines)
	ret.LineNumber = 1

	return ret, nil
}
