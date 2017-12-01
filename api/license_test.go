package api

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestGitHubAPINoKey(t *testing.T) {

	c := context.Background()
	l := &License{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}

	err := l.GetLicenses(c, nil, false)
	if err != nil {
		panic(err)
	}
	if l.Shortname != color.New(color.FgGreen).Sprintf("MIT") {
		t.Errorf("API did not return correct license or color.")
	}

}

func TestNonexistingLicense(t *testing.T) {

	c := context.Background()
	l := &License{
		URL:     "github.com/denysdovhan/wtfjs",
		Host:    "github.com",
		Author:  "denysdovhan",
		Project: "wtfjs",
	}

	err := l.GetLicenses(c, nil, false)
	if err != nil {
		panic(err)
	}
	if l.Shortname != color.New(color.FgYellow).Sprintf("wtfpl") {
		t.Errorf("API did not return correct license or color.")
	}

}

func TestGitHubAPIWithKey(t *testing.T) {

	c := context.Background()
	l := &License{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}

	v := map[string]string{
		"github.com": "apikey",
	}

	err := l.GetLicenses(c, v, false)
	if err == nil {
		t.Errorf("expected bad credentials error")
	}

}
func TestWriteToFile(t *testing.T) {

	dir, err := ioutil.TempDir("", "licenses")
	// check err
	defer os.RemoveAll(dir)
	l := &License{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}
	content := "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMTcgRW1p\nciBSaWJpYwpDb3B5cmlnaHQgKGMpIDIwMTYgQXN1a2EgU3V6dWtpCgpQZXJt\naXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8g\nYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mCnRoaXMgc29mdHdhcmUg\nYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0\nd2FyZSIpLCB0byBkZWFsIGluCnRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3Ry\naWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdo\ndHMgdG8KdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0\ncmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbCBjb3BpZXMgb2YKdGhl\nIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUg\nU29mdHdhcmUgaXMgZnVybmlzaGVkIHRvIGRvIHNvLApzdWJqZWN0IHRvIHRo\nZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQg\nbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGlu\nY2x1ZGVkIGluIGFsbApjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMg\nb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJB\nUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1Mg\nT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUg\nV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1MKRk9SIEEg\nUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5P\nIEVWRU5UIFNIQUxMIFRIRSBBVVRIT1JTIE9SCkNPUFlSSUdIVCBIT0xERVJT\nIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSIExJ\nQUJJTElUWSwgV0hFVEhFUgpJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRP\nUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTgpD\nT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RI\nRVIgREVBTElOR1MgSU4gVEhFIFNPRlRXQVJFLgo=\n"

	err = l.writeToFile(content, dir)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filepath.Join(dir, (l.Author + "-" + l.Project + "-license.MD"))); os.IsNotExist(err) {
		t.Errorf("License not returned")
	}

}

func TestInvalidFile(t *testing.T) {

	dir, err := ioutil.TempDir("", "licenses")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	var l *License
	content := "[]!"

	assert.Panics(t, func() { l.writeToFile(content, dir) }, "Panic - invalid base64 character")

}
