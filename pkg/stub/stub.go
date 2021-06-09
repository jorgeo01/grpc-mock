package stub

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	log "github.com/golang/glog"
	"github.com/monlabs/grpc-mock/pkg/models"
)

type Stub struct {
	Service  string  `json:"service"`
	Method   string  `json:"method"`
	Metadata *Input  `json:"metadata"`
	In       *Input  `json:"in"`
	Out      *Output `json:"out"`
}

type Input struct {
	Equals   map[string]interface{} `json:"equals"`
	Contains map[string]interface{} `json:"contains"`
	Has      map[string]interface{} `json:"has"`
	Matches  map[string]interface{} `json:"matches"`
}

type Output struct {
	Data    map[string]interface{} `json:"data"`
	Code    int32                  `json:"code"`
	Message string                 `json:"message"`
}

func (s *Stub) Validate() error {
	if s.Service == "" {
		return errors.New("missing service")
	}
	if s.Method == "" {
		return errors.New("missing method")
	}

	var metadata bool
	if s.Metadata != nil && len(s.Metadata.Equals)+len(s.Metadata.Contains)+len(s.Metadata.Has)+len(s.Metadata.Matches) > 0 {
		metadata = true
	}
	if !metadata && s.In == nil {
		return errors.New("missing input")
	}
	if !metadata && len(s.In.Equals)+len(s.In.Contains)+len(s.In.Matches)+len(s.In.Has) == 0 {
		return errors.New("require at least one of equals, contains or matches")
	}
	if s.Out == nil {
		return errors.New("missing output")
	}
	return nil
}

func (s *Stub) Match(req models.Request) bool {
	mdOk := s.match(s.Metadata, req.Metadata)
	fmt.Println("mdOK: ", mdOk)
	if !mdOk && s.Metadata != nil {
		fmt.Println("Exiting now")
		return false
	}

	if s.In != nil {
		res := s.match(s.In, req.Data)
		fmt.Println("Checked in:", res)
		return res
	}

	return mdOk
}

func (s *Stub) match(input *Input, in map[string]interface{}) bool {
	if input == nil {
		return false
	}
	if input.Equals != nil {
		return equals(input.Equals, in)
	}
	if input.Contains != nil {
		return contains(input.Contains, in)
	}
	if input.Has != nil {
		return contains(in, input.Has)
	}
	if input.Matches != nil {
		return matches(input.Matches, in)
	}
	return false
}

func equals(pattern, in map[string]interface{}) bool {
	return reflect.DeepEqual(pattern, in)
}

func contains(pattern, in map[string]interface{}) bool {
	for k, v := range in {
		p := pattern[k]
		if p == nil || !reflect.DeepEqual(p, v) {
			return false
		}
	}
	return true
}

func matches(pattern, in map[string]interface{}) bool {
	for k, v := range in {
		valStr, ok := v.(string)
		if !ok {
			return false
		}

		pStr, ok := pattern[k].(string)
		if !ok {
			return false
		}

		match, err := regexp.Match(pStr, []byte(valStr))
		if err != nil {
			log.Errorf("match regexp '%s' with '%s' failed: %v", pStr, valStr, err)
		}

		if !match {
			return false
		}
	}
	return true
}
