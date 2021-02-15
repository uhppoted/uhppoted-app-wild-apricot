package acl

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/logger"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/sirupsen/logrus"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type Rules struct {
	hash    []byte
	library *ast.KnowledgeLibrary
}

func NewRules(ruleset []byte, debug bool) (*Rules, error) {
	if debug {
		logger.SetLogLevel(logrus.TraceLevel)
	} else {
		logger.SetLogLevel(logrus.ErrorLevel)
	}

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

	startDate := startOfYear()
	endDate := endOfYear().AddDate(0, 1, 0)

	for _, m := range members.Members {
		if m.CardNumber != nil {
			r := record{
				Name:       m.Name,
				CardNumber: uint32(*m.CardNumber),
				StartDate:  startDate,
				EndDate:    endDate,
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			}

			if err := rules.eval(m, &r); err != nil {
				return nil, err
			}

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
	kb := rules.library.NewKnowledgeBaseInstance("acl", "0.0.0")

	return enjin.Execute(context, kb)
}
