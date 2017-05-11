package bigot

import (
	"fmt"
	"testing"

	"github.com/rexlv/bigot/provider/file"
)

func Test_ReadFile(t *testing.T) {
	p := file.NewProvider("/tmp/config.yaml")
	b := New(p)
	b.ReadInConfig()
	a := b.GetInt("wlog.formatters.f1.level")
	fmt.Printf("a: %#v\n", a)
}
