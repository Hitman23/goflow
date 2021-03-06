package flows

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
)

// Result represents a result value in our flow run. Results have a name for which they are the result for,
// the value itself of the result, optional category and the date and node the result was collected on
type Result struct {
	Name              string    `json:"name"`
	Value             string    `json:"value"`
	Category          string    `json:"category,omitempty"`
	CategoryLocalized string    `json:"category_localized,omitempty"`
	NodeUUID          NodeUUID  `json:"node_uuid"`
	Input             string    `json:"input"`
	CreatedOn         time.Time `json:"created_on"`
}

// Resolve resolves the passed in key to a value. Result values have a name, value, category, node and created_on
func (r *Result) Resolve(key string) interface{} {
	switch key {
	case "name":
		return r.Name
	case "value":
		return r.Value
	case "category":
		return r.Category
	case "category_localized":
		if r.CategoryLocalized == "" {
			return r.Category
		}
		return r.CategoryLocalized
	case "created_on":
		return r.CreatedOn
	}

	return fmt.Errorf("no field '%s' on result", key)
}

// Default returns the default value for a result, which is our value
func (r *Result) Default() interface{} {
	return r.Value
}

// String returns the string representation of a result, which is our value
func (r *Result) String() string {
	return r.Value
}

var _ utils.VariableResolver = (*Result)(nil)

// Results is our wrapper around a map of snakified result names to result objects
type Results map[string]*Result

func NewResults() Results {
	return make(Results, 0)
}

func (r Results) Clone() Results {
	clone := make(Results, len(r))
	for k, v := range r {
		clone[k] = v
	}
	return clone
}

// Save saves a new result in our map. The key is saved in a snakified format
func (r Results) Save(name string, value string, category string, categoryLocalized string, nodeUUID NodeUUID, input string, createdOn time.Time) {
	r[utils.Snakify(name)] = &Result{
		Name:              name,
		Value:             value,
		Category:          category,
		CategoryLocalized: categoryLocalized,
		NodeUUID:          nodeUUID,
		Input:             input,
		CreatedOn:         createdOn,
	}
}

// Resolve resolves the passed in key, which is snakified before lookup
func (r Results) Resolve(key string) interface{} {
	key = utils.Snakify(key)
	result, ok := r[key]
	if !ok {
		return nil
	}

	return result
}

// Default returns the default value for our Results, which is the entire map
func (r Results) Default() interface{} {
	return r
}

// String returns the string representation of our Results, which is a key/value pairing of our fields
func (r Results) String() string {
	results := make([]string, 0, len(r))
	for _, v := range r {
		results = append(results, fmt.Sprintf("%s: %s", v.Name, v.Value))
	}
	return strings.Join(results, ", ")
}

var _ utils.VariableResolver = (*Results)(nil)
