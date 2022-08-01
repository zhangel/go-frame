package combined

import (
	"time"

	"github.com/zhangel/go-frame.git/config"
)

type CombinedConfigView struct {
	views []config.View
}

func (s *CombinedConfigView) Bool(key string) bool {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Bool(key)
		}
	}

	return false
}

func (s *CombinedConfigView) Bytes(key string) []byte {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Bytes(key)
		}
	}

	return []byte{}
}

func (s *CombinedConfigView) Duration(key string) time.Duration {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Duration(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) Float64(key string) float64 {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Float64(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) GetByPrefix(prefix string) map[string]string {
	result := map[string]string{}
	for _, view := range s.views {
		for k, v := range view.GetByPrefix(prefix) {
			result[k] = v
		}
	}

	return result
}

func (s *CombinedConfigView) Has(key string) bool {
	for _, view := range s.views {
		if view.Has(key) {
			return true
		}
	}

	return false
}

func (s *CombinedConfigView) Int(key string) int {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Int(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) Int64(key string) int64 {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Int64(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) Int64List(key string) []int64 {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Int64List(key)
		}
	}

	return []int64{}
}

func (s *CombinedConfigView) IntList(key string) []int {
	for _, view := range s.views {
		if view.Has(key) {
			return view.IntList(key)
		}
	}

	return []int{}
}

func (s *CombinedConfigView) String(key string) string {
	for _, view := range s.views {
		if view.Has(key) {
			return view.String(key)
		}
	}

	return ""
}

func (s *CombinedConfigView) StringList(key string) []string {
	for _, view := range s.views {
		if view.Has(key) {
			return view.StringList(key)
		}
	}

	return []string{}
}

func (s *CombinedConfigView) StringMap(key string) map[string]string {
	for _, view := range s.views {
		if view.Has(key) {
			return view.StringMap(key)
		}
	}

	return map[string]string{}
}

func (s *CombinedConfigView) Uint(key string) uint {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Uint(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) Uint64(key string) uint64 {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Uint64(key)
		}
	}

	return 0
}

func (s *CombinedConfigView) Uint64List(key string) []uint64 {
	for _, view := range s.views {
		if view.Has(key) {
			return view.Uint64List(key)
		}
	}

	return []uint64{}
}

func (s *CombinedConfigView) UintList(key string) []uint {
	for _, view := range s.views {
		if view.Has(key) {
			return view.UintList(key)
		}
	}

	return []uint{}
}

func NewCombinedConfigView(views ...config.View) *CombinedConfigView {
	return &CombinedConfigView{views}
}
