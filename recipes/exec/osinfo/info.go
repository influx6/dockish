// Package osinfo is taking from the github.com/docker/machine project.
package osinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/influx6/box/recipes/exec"
	"github.com/influx6/faux/context"
)

// OSInfo retrieves the OSRelease details related to the operating system.
func OSInfo(ctx context.CancelContext) (*OsRelease, error) {
	var outs bytes.Buffer
	lsCmd := exec.New(exec.Command("cat /etc/hosts"), exec.Sync(), exec.Output(&outs))

	if err := lsCmd.Exec(ctx); err != nil {
		return nil, err
	}

	return NewOsRelease(outs.Bytes())
}

// The /etc/os-release file contains operating system identification data
// See http://www.freedesktop.org/software/systemd/man/os-release.html for more details

// OsRelease reflects values in /etc/os-release
// Values in this struct must always be string
// or the reflection will not work properly.
type OsRelease struct {
	AnsiColor    string `osr:"ANSI_COLOR"`
	Name         string `osr:"NAME"`
	Version      string `osr:"VERSION"`
	Variant      string `osr:"VARIANT"`
	VariantID    string `osr:"VARIANT_ID"`
	ID           string `osr:"ID"`
	IDLike       string `osr:"ID_LIKE"`
	PrettyName   string `osr:"PRETTY_NAME"`
	VersionID    string `osr:"VERSION_ID"`
	HomeURL      string `osr:"HOME_URL"`
	SupportURL   string `osr:"SUPPORT_URL"`
	BugReportURL string `osr:"BUG_REPORT_URL"`
}

// NewOsRelease returns a giving OSRelease instance from the provided content.
func NewOsRelease(contents []byte) (*OsRelease, error) {
	osr := &OsRelease{}
	if err := osr.ParseOsRelease(contents); err != nil {
		return nil, err
	}
	return osr, nil
}

// ParseOsRelease attempts to parse the provided data.
func (osr *OsRelease) ParseOsRelease(osReleaseContents []byte) error {
	r := bytes.NewReader(osReleaseContents)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		key, val, err := parseLine(scanner.Text())
		if err != nil {
			continue
		}
		if err := osr.setIfPossible(key, val); err != nil {
			return err
		}
	}
	return nil
}

func stripQuotes(val string) string {
	if len(val) > 0 && val[0] == '"' {
		return val[1 : len(val)-1]
	}
	return val
}

func (osr *OsRelease) setIfPossible(key, val string) error {
	v := reflect.ValueOf(osr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)
		originalName := fieldType.Tag.Get("osr")
		if key == originalName && fieldValue.Kind() == reflect.String {
			fieldValue.SetString(val)
			return nil
		}
	}
	return fmt.Errorf("Couldn't set key %s, no corresponding struct field found", key)
}

func parseLine(osrLine string) (string, string, error) {
	if osrLine == "" {
		return "", "", nil
	}

	vals := strings.Split(osrLine, "=")
	if len(vals) != 2 {
		return "", "", fmt.Errorf("Expected %s to split by '=' char into two strings, instead got %d strings", osrLine, len(vals))
	}
	key := vals[0]
	val := stripQuotes(vals[1])
	return key, val, nil
}