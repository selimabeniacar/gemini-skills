package styles

// Colors defines the official color palette
var Colors = map[string]struct {
	Fill   string
	Stroke string
	Text   string
}{
	"service":  {Fill: "#228be6", Stroke: "#1971c2", Text: "#ffffff"},
	"entry":    {Fill: "#40c057", Stroke: "#2f9e44", Text: "#ffffff"},
	"kafka":    {Fill: "#12b886", Stroke: "#099268", Text: "#ffffff"},
	"database": {Fill: "#fab005", Stroke: "#f59f00", Text: "#000000"},
	"cache":    {Fill: "#be4bdb", Stroke: "#9c36b5", Text: "#ffffff"},
	"external": {Fill: "#868e96", Stroke: "#495057", Text: "#ffffff"},
	"error":    {Fill: "#fa5252", Stroke: "#e03131", Text: "#ffffff"},
	"warning":  {Fill: "#fd7e14", Stroke: "#e8590c", Text: "#ffffff"},
}

// Shapes defines the mermaid shape syntax for each node type
var Shapes = map[string]struct {
	Open  string
	Close string
}{
	"rectangle":        {Open: "[", Close: "]"},
	"cylinder":         {Open: "[(", Close: ")]"},
	"stadium":          {Open: "([", Close: "])"},
	"double_rectangle": {Open: "[[", Close: "]]"},
	"diamond":          {Open: "{", Close: "}"},
	"rounded":          {Open: "(", Close: ")"},
	"circle":           {Open: "((", Close: "))"},
}

// ArrowTypes defines the arrow syntax for each connection type
var ArrowTypes = map[string]string{
	"sync":     "==>",
	"async":    "-.->",
	"internal": "-->",
}

// NodeTypeToShape maps node types to their expected shapes
var NodeTypeToShape = map[string]string{
	"service":        "rectangle",
	"handler":        "rectangle",
	"database":       "cylinder",
	"kafka_topic":    "cylinder",
	"consumer_group": "stadium",
	"external":       "double_rectangle",
	"cache":          "rounded",
	"decision":       "diamond",
}

// NodeTypeToClass maps node types to their CSS class
var NodeTypeToClass = map[string]string{
	"service":        "service",
	"handler":        "service",
	"database":       "database",
	"kafka_topic":    "kafka",
	"consumer_group": "kafka",
	"external":       "external",
	"cache":          "cache",
	"entry":          "entry",
}
