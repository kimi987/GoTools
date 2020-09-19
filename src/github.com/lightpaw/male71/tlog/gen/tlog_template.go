package gen

const tmpl = `
package tlog

import (
	"bytes"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/shared_proto"
)

{{range .Structs}}
// {{.Desc}}
{{if .HasHeroField}}
func (s *TlogService) Tlog{{hump .Name}}ById(heroId int64{{.DefineFieldString}}) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.{{hump .Name}}ById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.Tlog{{hump .Name}}(hero{{.PassFieldString}})
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.{{hump .Name}}ById hero not found")
	}
}
{{end}}

func (s *TlogService) Tlog{{hump .Name}}({{.TlogPrefix}}{{.DefineFieldString}}) {
	s.WriteLog(s.build{{hump .Name}}({{.TlogPrefixPass}}{{.PassFieldString}}))
}

func (s *TlogService) build{{hump .Name}}({{.TlogPrefix}}{{.DefineFieldString}}) string {
	return s.{{.BuildMethod}}"{{hump .Name}}", func({{.BuildPrefix}}) string {
		return s.Build{{hump .Name}}({{.BuildPrefixPass}}{{.PassFieldString}})
	})
}

func (s *TlogService) Build{{hump .Name}}({{.BuildPrefix}}{{.DefineFieldString}}) string {

	buf := &bytes.Buffer{}
	buf.WriteString("{{hump .Name}}")
    {{.BuildFields}}

	str := buf.String()
	return str
}

{{end}}
`
