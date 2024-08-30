package database

const sqlFileDelimiter = "/*SQLEND*/"

var statementsInstance *Statements

func NewStatements() *Statements {
	if statementsInstance == nil {
		statementsInstance = &Statements{}
		statementsInstance.schemaTemplates = make(map[string]string)
	}

	return statementsInstance
}

func (statements *Statements) AddSchemaTemplateMap(templates map[string]string) error {
	if statements == nil || statements != statementsInstance {
		return ErrStatementsUninitialized
	}

	for key, template := range templates {
		statements.schemaTemplates[key] = template
	}

	return nil
}

func (statements *Statements) GetSchemaTemplate(key string) string {
	if statements == nil || statements != statementsInstance {
		return ""
	}

	return statements.schemaTemplates[key]
}
