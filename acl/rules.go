package acl

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/logger"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	core "github.com/uhppoted/uhppote-core/types"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type Rules struct {
	hash    []byte
	library *ast.KnowledgeLibrary
}

func NewRules(ruleset []byte, debug bool) (*Rules, error) {
	if debug {
		logger.SetLogLevel(logger.TraceLevel)
	} else {
		logger.SetLogLevel(logger.ErrorLevel)
	}

	// ... check for header and footer
	// Ref. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/2

	first := ""
	last := ""
	b := bytes.NewBuffer(ruleset)
	scanner := bufio.NewScanner(b)

	for scanner.Scan() {
		line := scanner.Text()
		if first == "" {
			first = line
		}

		if strings.TrimSpace(line) != "" {
			last = line
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !strings.Contains(first, "** GRULES **") || !strings.Contains(last, "*** END GRULES ***") {
		return nil, fmt.Errorf("invalid 'grules' file - missing start/end markers")
	}

	// ... parse
	hash := sha256.Sum256(ruleset)
	kb := ast.NewKnowledgeLibrary()
	if err := builder.NewRuleBuilder(kb).BuildRuleFromResource("acl", "0.0.0", pkg.NewBytesResource(ruleset)); err != nil {
		return nil, err
	}

	return &Rules{
		hash:    hash[:],
		library: kb,
	}, nil
}

func (rules *Rules) Updated(hash string) bool {
	if hash != "" && hash == hex.EncodeToString(rules.hash) {
		return false
	}

	return true
}

func (rules *Rules) MakeACL(members types.Members, doors []string) (*ACL, error) {
	acl := ACL{
		doors:   doors,
		records: []record{},
	}

	for _, m := range members.Members {
		r := record{
			Name:      m.Name,
			StartDate: startOfYear(),
			EndDate:   plusOneDay(endOfYear()),
			Granted:   map[string]interface{}{},
			Revoked:   map[string]struct{}{},
		}

		if m.CardNumber != nil {
			r.CardNumber = uint32(*m.CardNumber)
		}

		if err := rules.eval(m, &r); err != nil {
			return nil, err
		}

		if r.CardNumber > 0 {
			acl.records = append(acl.records, r)
		}
	}

	sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

	return &acl, nil
}

func (rules *Rules) MakeACLWithPIN(members types.Members, doors []string) (*ACL, error) {
	acl := ACL{
		doors:   doors,
		records: []record{},
	}

	for _, m := range members.Members {
		r := record{
			Name:      m.Name,
			PIN:       m.PIN,
			StartDate: startOfYear(),
			EndDate:   plusOneDay(endOfYear()),
			Granted:   map[string]interface{}{},
			Revoked:   map[string]struct{}{},
		}

		if m.CardNumber != nil {
			r.CardNumber = uint32(*m.CardNumber)
		}

		if err := rules.eval(m, &r); err != nil {
			return nil, err
		}

		if r.CardNumber > 0 {
			acl.records = append(acl.records, r)
		}
	}

	sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

	return &acl, nil
}

func (rules *Rules) Hash() string {
	if rules != nil {
		return hex.EncodeToString(rules.hash)
	}

	return ""
}

func (rules *Rules) eval(m types.Member, r *record) error {
	context := ast.NewDataContext()

	if err := context.Add("member", &m); err != nil {
		return err
	}

	if err := context.Add("permissions", r); err != nil {
		return err
	}

	enjin := engine.NewGruleEngine()

	if kb, err := rules.library.NewKnowledgeBaseInstance("acl", "0.0.0"); err != nil {
		return err
	} else {
		return enjin.Execute(context, kb)
	}
}

func plusOneDay(date core.Date) core.Date {
	d := time.Time(date).AddDate(0, 1, 0)
	year := d.Year()
	month := d.Month()
	day := d.Day()

	return core.ToDate(year, month, day)
}
