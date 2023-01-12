package log

import (
	api "github.com/ppmasa8/proglog/api/v1"
	"google.golang.org/protobuf/proto"
)

type segment struct {
	store *store
	index *index
	baseOffset, nextOffset uint64
	config Config
}