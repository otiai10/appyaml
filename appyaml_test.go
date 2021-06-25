package appyaml

import (
	"testing"

	. "github.com/otiai10/mint"
)

func TestLoad(t *testing.T) {
	app, err := Load("./testdata/app.yaml")
	Expect(t, err).ToBe(nil)
	Expect(t, app).Not().ToBe(nil)
}
