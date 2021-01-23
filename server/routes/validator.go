package routes

import (
	"errors"
)

type (
	Validtable interface {
		validate() error
	}
	Validator struct{}
	RuleField struct {
		message string
		Rule
	}
	Rule interface {
		validate(value interface{}) error
	}
	FieldRules struct {
		field interface{}
		rules []Rule
	}

	required struct{}
)

func (required) validate(value interface{}) error {
	var err error
	if err = nil; value == nil {
		if str, ok := value.(string); ok {
			err = errors.New(str + ": is required.")
		} else {
			err = errors.New("This is required")
		}

	}

	return err
}

func Required(message string) *RuleField {

	return &RuleField{
		message: message,
		Rule:    required{},
	}
}

func (v Validator) Field(field *interface{}, rules ...Rule) []error {

	var errors = []error{}

	for _, rl := range rules {
		er := rl.validate(field)
		if er != nil {
			errors = append(errors, er)
		}
	}

	return errors

}

func ValidateStruct(obj interface{}, validators ...FieldRules) {

}
