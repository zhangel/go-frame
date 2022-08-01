package test

import (
	"fmt"
	"testing"

	"github.com/zhangel/go-frame.git/config"
	"github.com/zhangel/go-frame.git/config_plugins"
)

type WatcherMock struct {
	tag string
}

func (s *WatcherMock) OnUpdate(v map[string]string) {
	fmt.Printf("%s: OnUpdate, v = %+v\n", s.tag, v)
}

func (s *WatcherMock) OnDelete(v []string) {
	fmt.Printf("%s: OnDelete, v = %+v\n", s.tag, v)
}

func (s *WatcherMock) OnSync(v map[string]string) {
	fmt.Printf("%s: OnSync, v = %+v\n", s.tag, v)
}

func TestMerge(t *testing.T) {
	memorySourceLow := config_plugins.NewMemoryConfigSource(map[string]interface{}{
		"key":    "val",
		"ns.key": "ns.val",
	})

	memorySourceHigh := config_plugins.NewMemoryConfigSource(map[string]interface{}{
		"key": "high.val",
	})

	config, err := config.NewConfig([]string{"ns"}, memorySourceLow, memorySourceHigh)
	if err != nil {
		t.Fatal(err)
	}

	val := config.String("key")
	if val != "high.val" {
		t.Fatalf("val == %s, 'high.val' expected", val)
	}
}

func TestMemoryConfig(t *testing.T) {
	memorySourceLow := config_plugins.NewMemoryConfigSource(map[string]interface{}{
		"K1": "L_V_1",
		"K2": "L_V_2",
		"K3": "L_V_3",
		"K4": "L_V_4",
		"K5": "L_V_5",
	})

	memorySourceHigh := config_plugins.NewMemoryConfigSource(map[string]interface{}{
		"K1":             "H_V_1",
		"K2":             "H_V_2",
		"prefix.K3":      "H_V_3",
		"K4":             "H_V_4",
		"K5":             "H_V_5",
		"high.K1":        "H_H_V_1",
		"high.K2":        "H_H_V_2",
		"high.prefix.K3": "H_H_V_3",
		"high.K4":        "H_H_V_4",
		"high.K5":        "H_H_V_5",
	})

	config, err := config.NewConfig([]string{"", "high"}, memorySourceLow, memorySourceHigh)
	if err != nil {
		t.Fatal(err)
	}

	prefixConfig := config.WithPrefix("prefix")

	config.Watch(&WatcherMock{"config"})
	prefixConfig.Watch(&WatcherMock{"prefixConfig"})

	if config.String("K1") != "H_H_V_1" {
		t.Fatal()
	}

	memorySourceHigh.Del("high.K1")

	if config.String("K1") != "H_V_1" {
		t.Fatal()
	}

	memorySourceLow.Put("K1", "Low")
	memorySourceHigh.Put("high.K1", "End")

	if config.String("K1") != "End" {
		t.Fatal()
	}

	memorySourceHigh.Del("high.K1")
	if config.String("K1") != "H_V_1" {
		t.Fatal()
	}

	memorySourceHigh.Del("K1")
	if config.String("K1") != "Low" {
		t.Fatal()
	}

	memorySourceLow.Put("k1", "L_V_1")
	if config.String("K1") != "L_V_1" {
		t.Fatal()
	}

	memorySourceLow.Del("K1")
	if config.String("K1") != "" {
		t.Fatal()
	}

	if prefixConfig.String("K3") != "H_H_V_3" {
		t.Fatal()
	}

	memorySourceHigh.Del("high.prefix.K3")

	if prefixConfig.String("K3") != "H_V_3" {
		t.Fatal()
	}

	memorySourceHigh.Del("prefix.K3")
	if prefixConfig.String("K3") != "" {
		t.Fatal()
	}
}

func BenchmarkMemoryConfig(b *testing.B) {
	memorySource := config_plugins.NewMemoryConfigSource(map[string]interface{}{
		"n1.k1": "v1",
	})

	config, err := config.NewConfig([]string{""}, memorySource)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if config.String("n1.k1") != "v1" {
			b.Fatal("n1.k1 != v1")
		}
	}
	b.StopTimer()
}
