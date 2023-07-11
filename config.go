package main

import (
	"encoding/json"

	xds "github.com/cncf/xds/go/xds/type/v3"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/api"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	http.RegisterHttpFilterConfigFactory("basic-auth", configFactory)
	http.RegisterHttpFilterConfigParser(&parser{})
}

type config struct {
	users []User
}

type User struct {
	Username string
	Password string
}

type parser struct {
}

func (p *parser) Parse(any *anypb.Any) (interface{}, error) {
	configStruct := &xds.TypedStruct{}
	if err := any.UnmarshalTo(configStruct); err != nil {
		return nil, err
	}

	v := configStruct.Value
	conf := &config{}
	var users []User

	if userMap, ok := v.AsMap()["users"]; ok {
		data, err := json.Marshal(userMap)
		if err != nil {
			return conf, err
		}

		err = json.Unmarshal(data, &users)
		if err != nil {
			return conf, err
		}

		conf.users = users
	}

	return conf, nil
}

func (p *parser) Merge(parent interface{}, child interface{}) interface{} {
	parentConfig := parent.(*config)
	childConfig := child.(*config)

	newConfig := *parentConfig
	if childConfig.users != nil {
		newConfig.users = childConfig.users
	}
	return &newConfig
}

func configFactory(c interface{}) api.StreamFilterFactory {
	conf, ok := c.(*config)
	if !ok {
		panic("unexpected config type")
	}
	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &filter{
			callbacks: callbacks,
			config:    conf,
		}
	}
}
