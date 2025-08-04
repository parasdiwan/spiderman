package publish

type ErrType string

const (
	ErrTypeNotFound ErrType = "not_found"
	ErrTypeNoAccess ErrType = "no_access"
	ErrTypeInternal ErrType = "internal_issue"
	ErrTypeUnknown  ErrType = "unknown"
)

type Error struct {
	name ErrType
	error
}
