package generator

import (
	"testing"
	"time"

	_ "github.com/tonnerre/golang-pretty"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type propertySuite struct{}

var _ = Suite(&propertySuite{})

type ExampleJSONBasic struct {
	Omitted    string  `json:"-,omitempty"`
	Bool       bool    `json:",omitempty"`
	Integer    int     `json:",omitempty"`
	Integer8   int8    `json:",omitempty"`
	Integer16  int16   `json:",omitempty"`
	Integer32  int32   `json:",omitempty"`
	Integer64  int64   `json:",omitempty"`
	UInteger   uint    `json:",omitempty"`
	UInteger8  uint8   `json:",omitempty"`
	UInteger16 uint16  `json:",omitempty"`
	UInteger32 uint32  `json:",omitempty"`
	UInteger64 uint64  `json:",omitempty"`
	String     string  `json:",omitempty"`
	Bytes      []byte  `json:",omitempty"`
	Float32    float32 `json:",omitempty"`
	Float64    float64
	Interface  interface{} `required:"true"`
	Timestamp  time.Time   `json:",omitempty"`
}

func (self *propertySuite) TestLoad(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONBasic{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type:     "object",
			Required: []string{"Interface"},
			Properties: map[string]*property{
				"Bool":       &property{Type: "boolean"},
				"Integer":    &property{Type: "integer"},
				"Integer8":   &property{Type: "integer"},
				"Integer16":  &property{Type: "integer"},
				"Integer32":  &property{Type: "integer"},
				"Integer64":  &property{Type: "integer"},
				"UInteger":   &property{Type: "integer"},
				"UInteger8":  &property{Type: "integer"},
				"UInteger16": &property{Type: "integer"},
				"UInteger32": &property{Type: "integer"},
				"UInteger64": &property{Type: "integer"},
				"String":     &property{Type: "string"},
				"Bytes":      &property{Type: "string"},
				"Float32":    &property{Type: "number"},
				"Float64":    &property{Type: "number"},
				"Interface":  &property{},
				"Timestamp":  &property{Type: "string", Format: "date-time"},
			},
		},
	})
}

type ExampleJSONBasicWithTag struct {
	Bool         bool    `json:"test"`
	String       string  `json:"string" description:"blah" minLength:"3" maxLength:"10" pattern:"m{3,10}"`
	Const        string  `json:"const" const:"blah"`
	Float        float32 `json:"float" min:"1.5" max:"42"`
	Int          int64   `json:"int" exclusiveMin:"-10" exclusiveMax:"0"`
	AnswerToLife int     `json:"answer" const:"42"`
	Fruit        string  `json:"fruit" enum:"apple|banana|pear"`
}

func (self *propertySuite) TestLoadWithTag(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONBasicWithTag{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"test": &property{Type: "boolean"},
				"string": &property{
					Type:        "string",
					MinLength:   int64ptr(3),
					MaxLength:   int64ptr(10),
					Pattern:     "m{3,10}",
					Description: "blah",
				},
				"const": &property{
					Type:  "string",
					Const: "blah",
				},
				"float": &property{
					Type:    "number",
					Minimum: float64ptr(1.5),
					Maximum: float64ptr(42),
				},
				"int": &property{
					Type:             "integer",
					ExclusiveMinimum: float64ptr(-10),
					ExclusiveMaximum: float64ptr(0),
				},
				"answer": &property{
					Type:  "integer",
					Const: int64(42),
				},
				"fruit": &property{
					Type: "string",
					Enum: []string{"apple", "banana", "pear"},
				},
			},
		},
	})
}

type ExampleJSONBasicSlices struct {
	Slice            []string      `json:",foo,omitempty"`
	SliceOfInterface []interface{} `json:",foo" required:"true"`
}

func (self *propertySuite) TestLoadSliceAndContains(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONBasicSlices{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Slice": &property{
					Type:  "array",
					Items: &property{Type: "string"},
				},
				"SliceOfInterface": &property{
					Type: "array",
				},
			},

			Required: []string{"SliceOfInterface"},
		},
	})
}

type ExampleJSONNestedStruct struct {
	Struct struct {
		Foo string `required:"true"`
	}
}

func (self *propertySuite) TestLoadNested(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONNestedStruct{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Struct": &property{
					Type: "object",
					Properties: map[string]*property{
						"Foo": &property{Type: "string"},
					},
					Required: []string{"Foo"},
				},
			},
		},
	})
}

type ExampleJSONBasicMaps struct {
	Maps           map[string]string `json:",omitempty"`
	MapOfInterface map[string]interface{}
}

func (self *propertySuite) TestLoadMap(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONBasicMaps{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Maps": &property{
					Type: "object",
					Properties: map[string]*property{
						".*": &property{Type: "string"},
					},
					AdditionalProperties: false,
				},
				"MapOfInterface": &property{
					Type:                 "object",
					AdditionalProperties: true,
				},
			},
		},
	})
}

func (self *propertySuite) TestLoadNonStruct(c *C) {
	j := &Document{}
	j.Read([]string{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type:  "array",
			Items: &property{Type: "string"},
		},
	})
}

func (self *propertySuite) TestString(c *C) {
	j := &Document{}
	j.Read(true)

	expected := "{\n" +
		"    \"$schema\": \"http://json-schema.org/schema#\",\n" +
		"    \"type\": \"boolean\"\n" +
		"}"

	c.Assert(j.String(), Equals, expected)
}

func (self *propertySuite) TestMarshal(c *C) {
	j := &Document{}
	j.Read(10)

	expected := "{\n" +
		"    \"$schema\": \"http://json-schema.org/schema#\",\n" +
		"    \"type\": \"integer\"\n" +
		"}"

	json := j.String()
	c.Assert(string(json), Equals, expected)
}

type ExampleJSONNestedSliceStruct struct {
	Struct  []ItemStruct
	Struct2 []*ItemStruct
}
type ItemStruct struct {
	Foo string `required:"true"`
}

func (self *propertySuite) TestLoadNestedSlice(c *C) {
	j := &Document{}
	j.Read(&ExampleJSONNestedSliceStruct{})

	c.Assert(*j, DeepEquals, Document{
		Schema: "http://json-schema.org/schema#",
		property: property{
			Type: "object",
			Properties: map[string]*property{
				"Struct": &property{
					Type: "array",
					Items: &property{
						Type: "object",
						Properties: map[string]*property{
							"Foo": &property{Type: "string"},
						},
						Required: []string{"Foo"},
					},
				},
				"Struct2": &property{
					Type: "array",
					Items: &property{
						Type: "object",
						Properties: map[string]*property{
							"Foo": &property{Type: "string"},
						},
						Required: []string{"Foo"},
					},
				},
			},
		},
	})
}
