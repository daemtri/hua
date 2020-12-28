package huarpc

import (
	"testing"

	"github.com/duanqy/hua/example/api"
)

func TestParse(t *testing.T) {
	parseService(&api.CalcService{})
}
