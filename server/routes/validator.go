

type (


	 Validtable interface { 
		 validate()error
	}
	
	
	 Validator struct {
	
	}

	RuleField struct {
		message string
		Rule
	}
	
	 Rule interface {
		validate(field interface{}) error	
	}
	
	 FieldRules struct {
		field interface{}
		rules []Rule
	}

	
)

func Required (message string) RuleField {

	return &RuleField{
		message:message,
		validate: func(){
			return 
		}
	}
} 

func(v Validator) Field(field *interface{},...Rule)  {

}

func ValidateStruct(obj interface{},validators ...FieldRules)  {



	return &FieldRules{
		field:obj,
		rules:FieldRules
	}
}

func (v Vli) validate(obj interface{},) error{

	return error("Something went wrong")
}