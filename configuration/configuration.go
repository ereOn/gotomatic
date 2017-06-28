// Package configuration provides functions to create and load configurations.
package configuration

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/intelux/gotomatic/conditional"
	"github.com/intelux/gotomatic/trigger"
)

// Configuration represents a configuration.
type Configuration interface {
	// GetCondition returns a named condition from the configuration, if it
	// finds it.
	//
	// Any attempt to close the returned the condition is without effect.
	GetCondition(name string) conditional.Condition

	// AddCondition adds a named condition to the configuration.
	//
	// The caller should never use the passed-in condition directly ever again.
	AddCondition(name string, condition conditional.Condition) error

	// Clear the configuration, freeing any resource associated with it.
	//
	// The configuration can be reused again after a call to Clear.
	Clear()

	// Close the configuration.
	//
	// Must be called when the configuration is no longer needed, to free all
	// its resources.
	Close()
}

// New creates a new empty configuration.
func New() Configuration {
	return newConfigurationImpl()
}

// Load a configuration from the specified reader, as a YAML stream.
func Load(r io.Reader) (Configuration, error) {
	b, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var data interface{}
	err = yaml.Unmarshal(b, &data)

	if err != nil {
		return nil, err
	}

	return Decode(data)
}

// ConditionTrigger represents a trigger linked to a conditional.Condition.
type ConditionTrigger struct {
	trigger.Trigger
	Condition conditional.Condition
}

// Decode a configuration.
func Decode(data interface{}) (Configuration, error) {
	configuration := newConfigurationImpl()

	var decl struct {
		Conditions []conditional.Condition
		Triggers   []ConditionTrigger
	}

	if err := configuration.decode(data, &decl); err != nil {
		return nil, err
	}

	configuration.triggers = make([]trigger.Trigger, len(decl.Triggers))

	for i, trigger := range decl.Triggers {
		configuration.triggers[i] = trigger.Trigger
	}

	return configuration, nil
}

type configurationImpl struct {
	namedConditions map[string]conditional.Condition
	triggers        []trigger.Trigger
}

func newConfigurationImpl() *configurationImpl {
	return &configurationImpl{
		namedConditions: make(map[string]conditional.Condition),
	}
}

func (c *configurationImpl) GetCondition(name string) conditional.Condition {
	condition := c.namedConditions[name]

	if condition != nil {
		return conditional.Dereference(condition)
	}

	return nil
}

func (c *configurationImpl) AddCondition(name string, condition conditional.Condition) error {
	if _, ok := c.namedConditions[name]; ok {
		return fmt.Errorf("a condition named \"%s\" already exists", name)
	}

	c.namedConditions[name] = condition

	return nil
}

func (c *configurationImpl) Clear() {
	for _, condition := range c.namedConditions {
		condition.Close()
	}

	c.namedConditions = nil
}

func (c *configurationImpl) Close() { c.Clear() }
