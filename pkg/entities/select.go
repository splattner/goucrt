package entities

import "fmt"

type SelectEntityState EntityState
type SelectEntityAttribute EntityAttribute
type SelectEntityCommand EntityCommand

const (
	OnSelectEntityState SelectEntityState = "ON"
)

const (
	StateSelectEntityAttribute         SelectEntityAttribute = "state"
	CurrentOptionSelectEntityAttribute SelectEntityAttribute = "current_option"
	OptionsSelectEntityAttribute       SelectEntityAttribute = "options"
)

const (
	// SelectOptionSelectEntityCommand takes a required "option" param: the option to select,
	// which must be one of the values in the options attribute.
	SelectOptionSelectEntityCommand SelectEntityCommand = "select_option"
	SelectFirstSelectEntityCommand  SelectEntityCommand = "select_first"
	SelectLastSelectEntityCommand   SelectEntityCommand = "select_last"
	// SelectNextSelectEntityCommand and SelectPreviousSelectEntityCommand take an optional
	// "cycle" bool param (default true): whether to wrap around at the end/start of the list.
	SelectNextSelectEntityCommand     SelectEntityCommand = "select_next"
	SelectPreviousSelectEntityCommand SelectEntityCommand = "select_previous"
)

// SelectEntity offers a limited set of selectable options, defined by the integration driver as a
// static or dynamically-changing list. It has no features and no options (per the spec - not to be
// confused with the entity's own "options" attribute, the list of choices).
type SelectEntity struct {
	BaseEntity
	Commands map[SelectEntityCommand]func(SelectEntity, map[string]interface{}) int `json:"-"`
}

// NewSelectEntity creates a select entity with the given initial list of choices. The first option
// becomes the initial current_option, if any are given.
func NewSelectEntity(id string, name LanguageText, area string, options []string) *SelectEntity {

	selectEntity := SelectEntity{}
	selectEntity.Id = id
	selectEntity.Name = name
	selectEntity.Area = area

	selectEntity.Type = "select"

	selectEntity.Commands = make(map[SelectEntityCommand]func(SelectEntity, map[string]interface{}) int)
	selectEntity.Attributes = make(map[string]interface{})

	selectEntity.AddAttribute(string(StateSelectEntityAttribute), OnSelectEntityState)
	selectEntity.AddAttribute(string(OptionsSelectEntityAttribute), options)
	currentOption := ""
	if len(options) > 0 {
		currentOption = options[0]
	}
	selectEntity.AddAttribute(string(CurrentOptionSelectEntityAttribute), currentOption)

	return &selectEntity
}

func (e *SelectEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(SelectEntity)
	if !ok {
		return fmt.Errorf("cannot update SelectEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Commands = updated.Commands
	e.Attributes = updated.Attributes

	return nil
}

// Register a function for the Entity command
func (e *SelectEntity) AddCommand(command SelectEntityCommand, function func(SelectEntity, map[string]interface{}) int) {
	e.Commands[command] = function
}

// Map a SelectEntityCommand to a function call with params (select_option, select_next, select_previous)
func (e *SelectEntity) MapCommandWithParams(command SelectEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity SelectEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

// Map a SelectEntityCommand to a function call without params (select_first, select_last)
func (e *SelectEntity) MapCommand(command SelectEntityCommand, f func() error) {

	e.AddCommand(command, func(entity SelectEntity, params map[string]interface{}) int {

		if err := f(); err != nil {
			return 404
		}
		return 200
	})
}

// Call the registred function for this entity_command
func (e *SelectEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[SelectEntityCommand(cmd_id)] != nil {
		return e.Commands[SelectEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
