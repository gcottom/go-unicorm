package unicorm

import "strings"

// subjects
const GET = "Get"
const FIND = "Find"
const SEARCH = "Search"
const QUERY = "Query"
const READ = "Read"
const SAVE = "Save"
const UPDATE = "Update"
const DELETE = "Delete"
const REMOVE = "Remove"
const COUNT = "Count"
const DISTINCT = "Distinct"
const EXISTS = "Exists"
const FIRST = "First"
const TOP = "Top"

// predicates
const AND = "And"
const OR = "Or"
const AFTER = "After"
const BEFORE = "Before"
const CONTAINING = "Containing"
const BETWEEN = "Between"
const ENDING_WITH = "EndingWith"
const FALSE = "False"
const GREATER_THAN = "GreaterThan"
const GREATER_THAN_EQUALS = "GreaterThanEquals"
const IN = "In"
const IS = "Is"
const IS_EMPTY = "IsEmpty"
const IS_NOT_EMPTY = "IsNotEmpty"
const IS_NOT_NULL = "IsNotNull"
const IS_NULL = "IsNull"
const LESS_THAN = "LessThan"
const LESS_THAN_EQUAL = "LessThanEqual"
const LIKE = "Like"
const NEAR = "Near"
const NOT = "Not"
const NOT_IN = "NotIn"
const NOT_LIKE = "NotLike"
const REGEX = "Regex"
const STARTING_WITH = "StartingWIth"
const TRUE = "True"
const WITHIN = "Within"

const QUERY_JOINER_BY = "By"
const QUERY_JOINER_WHERE = "Where"

const QUERY_SHORT_ALL = "All"

var subjectList []string
var predicateList []string
var queryJoiners []string
var subjectMap map[string]string
var predicateMap map[string]string
var querySubjects []string
var statementSubjects []string

func init() {
	initMaps()
}

func initMaps() {
	subjectList = []string{GET, FIND, SEARCH, QUERY, READ, SAVE, UPDATE, DELETE, REMOVE, COUNT, DISTINCT, EXISTS, FIRST, TOP}
	predicateList = []string{AND, OR, AFTER, BEFORE, CONTAINING, BETWEEN, ENDING_WITH, FALSE, GREATER_THAN, GREATER_THAN_EQUALS, IN, IS, IS_EMPTY, IS_NOT_EMPTY, IS_NOT_NULL, IS_NULL, LESS_THAN, LESS_THAN_EQUAL, LIKE, NEAR, NOT, NOT_IN, NOT_LIKE, REGEX, STARTING_WITH, TRUE, WITHIN}
	querySubjects = []string{GET, FIND, SEARCH, QUERY, READ, COUNT, DISTINCT, EXISTS, FIRST, TOP}
	queryJoiners = []string{QUERY_JOINER_BY, QUERY_JOINER_WHERE}
	statementSubjects = []string{SAVE, UPDATE, DELETE, REMOVE}
	subjectMap = make(map[string]string)
	predicateMap = make(map[string]string)
	for _, s := range subjectList {
		subjectMap[s] = strings.ToUpper(s)
	}
	for _, p := range predicateList {
		predicateMap[p] = strings.ToUpper(p)
	}

}
