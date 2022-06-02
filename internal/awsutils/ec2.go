package awsutils

import (
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	logger "github.com/sirupsen/logrus"
)

type ProfileData struct {
	Region string
	Name   string
}

func PrintOutSList(dl []types.Snapshot) error {
	tpl, err := template.New("").Parse(stemplat)
	if err != nil {
		logger.Errorf("parse template failure %v", err)
		return err
	}
	err = tpl.Execute(os.Stdout, dl)
	if err != nil {
		logger.Errorf("render template failure %v", err)
		return err
	}
	return nil
}

var stemplat string = `
List:
{{ range $index,$item := . -}}
-----------{{$index}}----------
tags: 
---
{{ range $tag := $item.Tags -}}
{{$tag.Key}}={{$tag.Value}}
{{ end -}}
---
id: {{$item.SnapshotId}}
encrypted: {{$item.Encrypted}}
volume: {{$item.VolumeId}}
state: {{$item.State}}
{{ end -}}

`

func PrintOutVList(dl []types.Volume) error {
	tpl, err := template.New("").Parse(vtemplat)
	if err != nil {
		logger.Errorf("parse template failure %v", err)
		return err
	}
	err = tpl.Execute(os.Stdout, dl)
	if err != nil {
		logger.Errorf("render template failure %v", err)
		return err
	}
	return nil
}

var vtemplat string = `
List:
{{ range $index,$item := . -}}
-----------{{$index}}----------
tags: 
---
{{ range $tag := $item.Tags -}}
{{$tag.Key}}={{$tag.Value}}
{{ end -}}
---
id: {{$item.VolumeId}}
encrypted: {{$item.Encrypted}}
az: {{$item.AvailabilityZone}}
state: {{$item.State}}
type: {{$item.VolumeType}}
{{ end -}}

`
