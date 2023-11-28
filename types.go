package feature

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/xuuid"
)

// ================================================================
type Md5Identifier string

func (ms *Md5Identifier) Scan(value any) error {
	if value == nil {
		*ms = ""
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Md5Identifier.Scan: expected []byte, got %T", value)
	}

	*ms = Md5Identifier(hex.EncodeToString(b))
	return nil
}

func (ms Md5Identifier) Value() (driver.Value, error) {
	b, err := hex.DecodeString(string(ms))
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ================================================================
type AccessRules struct {
	Subsets    []string `json:"subsets"`
	Exceptions []string `json:"exceptions,omitempty"`
}

func (r *AccessRules) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, r)
}

func (r AccessRules) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *AccessRules) AddSubset(rule string) {
	r.Subsets = append(r.Subsets, rule)
}

func (r *AccessRules) AddException(rule string) {
	r.Exceptions = append(r.Exceptions, rule)
}

func (r *AccessRules) RemoveRedundant() {
	r.Subsets = removeRedundant(r.Subsets)
	r.Exceptions = removeRedundant(r.Exceptions)
}

func (r *AccessRules) Merge(rules *AccessRules) {
	r.Subsets = append(r.Subsets, rules.Subsets...)
	r.Exceptions = append(r.Exceptions, rules.Exceptions...)
	r.RemoveRedundant()
}

func (r AccessRules) CanAccess(subset string) bool {
	if isCovered(subset, r.Subsets) {
		if isCovered(subset, r.Exceptions) {
			return false
		}
		return true
	}
	return false
}

// ================================================================
type AccessRulesWithBehavior struct {
	Behavior           string        `json:"behavior" db:"-" binding:"required"`
	AffectedEndpointId Md5Identifier `json:"affectedEndpointId" db:"endpoint_id" binding:"required"`
	AccessRules        *AccessRules  `json:"accessRules" db:"access_rules" binding:"required"`
}

func (r *AccessRulesWithBehavior) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, r)
}

func (r AccessRulesWithBehavior) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// ================================================================
type SubsetString string

func (s SubsetString) ExtractOwnerId(i int) (xuuid.UUID, error) {
	if (i <= 1) || (i%2 != 0) {
		return xuuid.UUID(uuid.Nil), errors.New("Invalid subset segment index")
	}

	segs := strings.Split(string(s), "/")
	if i >= len(segs) {
		return xuuid.UUID(uuid.Nil), errors.New("Invalid subset segment index")
	}

	return xuuid.Parse(segs[i])
}

// ================================================================
func removeRedundant(rules []string) []string {
	var patterns []*regexp.Regexp
	result := []string{}
	resultMap := map[string]bool{}

	for _, key := range rules {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				continue
			}
			patterns = append(patterns, re)
		}
	}

	for _, key := range rules {
		if !isCoveredByMoreGeneralPattern(key, patterns) && !resultMap[key] {
			result = append(result, key)
			resultMap[key] = true
		}
	}

	return result
}

func isCoveredByMoreGeneralPattern(key string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(key) && pattern.String() != "^"+strings.ReplaceAll(key, "*", ".*")+"$" {
			return true
		}
	}
	return false
}

func isCovered(rule string, rules []string) bool {
	if rules == nil {
		return false
	}

	var patterns []*regexp.Regexp

	for _, key := range rules {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			if re, err := regexp.Compile("^" + pattern + "$"); err == nil {
				patterns = append(patterns, re)
			}
		} else if rule == key {
			return true
		}
	}

	return isCoveredBy(rule, patterns)
}

func isCoveredBy(key string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}
