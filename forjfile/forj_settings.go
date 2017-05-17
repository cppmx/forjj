package forjfile

import "github.com/forj-oss/goforjj"

// forj/settings: Collection of key/value pair
type ForjSettingsStruct struct {
	is_template bool
	forge *ForgeYaml
	Organization string
	ForjSettingsStructTmpl `yaml:",inline"`
}

type ForjSettingsStructTmpl struct {
	Default DefaultSettingsStruct
	More map[string]string `yaml:",inline"`
}

type DefaultSettingsStruct struct {
	forge *ForgeYaml
	UpstreamInstance string `yaml:"upstream-instance"`
	Flow string
	More map[string]string `yaml:",inline"`
}

func (f *ForjSettingsStruct) MarshalYAML() (interface{}, error) {
	return f.ForjSettingsStructTmpl, nil
}

func (s *ForjSettingsStruct) Get(instance, key string) (value *goforjj.ValueStruct, _ bool) {
	if instance == "default" {
		return s.Default.Get(key)
	}
	switch key {
	case "organization":
		return value.Set(s.Organization, (s.Organization != ""))
	default:
		v, f := s.More[key]
		return value.Set(v, f)
	}
}

func (r *ForjSettingsStruct)SetHandler(instance string, from func(field string)(string, bool), keys...string) {
	for _, key := range keys {
		if v, found := from(key) ; found {
			r.Set(instance, key, v)
		}
	}
}

func (s *ForjSettingsStruct) Set(instance, key string, value string) {
	if instance == "default" {
		s.Default.Set(key, value)
		return
	}
	switch key {
	case "organization":
		s.Organization = value
		s.forge.dirty()
		return
	default:
		if v, found := s.More[key] ; found && v != value {
			s.forge.dirty()
			s.More[key] = value
		}
	}
}

func (g *ForjSettingsStruct) set_forge(f *ForgeYaml) {
	g.forge = f
	g.Default.set_forge(f)
}

func (s *DefaultSettingsStruct) Get(key string) (value *goforjj.ValueStruct, found bool) {
	switch key {
	case "upstream-instance":
		return value.Set(s.UpstreamInstance, (s.UpstreamInstance != ""))
	case "flow":
		return value.Set(s.Flow, (s.Flow != ""))
	default:
		v, f := s.More[key]
		return value.Set(v, f)
	}
}



func (s *DefaultSettingsStruct) Set(key string, value string) {
	switch key {
	case "upstream-instance":
		s.UpstreamInstance = value
		s.forge.dirty()
	case "flow":
		s.Flow = value
		s.forge.dirty()
		return
	default:
		if v, found := s.More[key] ; found && v != value {
			s.forge.dirty()
			s.More[key] = value
		}
	}
}

func (d *DefaultSettingsStruct) set_forge(f *ForgeYaml) {
	d.forge = f
}
