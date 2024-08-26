package database

const sqlFileDelimiter = "/*SQLEND*/"

var statementsInstance *Statements

func NewStatements() *Statements {
	if statementsInstance == nil {
		statementsInstance = &Statements{}

		statementsInstance.primarySchemaTemplates = make(map[string]string)
	}

	return statementsInstance
}

func (statements *Statements) AddPrimarySchemaTemplate(key, template string) error {
	if statements != statementsInstance {
		return ErrStatementsInvalid
	}

	if statements == nil {
		return ErrStatementsUninitialized
	}

	statements.primarySchemaTemplates[key] = template

	return nil
}

func (statements *Statements) AddPrimarySchemaTemplateMap(templates map[string]string) error {
	if statements != statementsInstance {
		return ErrStatementsInvalid
	}

	if statements == nil {
		return ErrStatementsUninitialized
	}

	for key, template := range templates {
		err := statements.AddPrimarySchemaTemplate(key, template)
		if err != nil {
			return err
		}

	}

	return nil
}
