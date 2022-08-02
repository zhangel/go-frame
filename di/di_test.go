package di

import (
	"fmt"
	"testing"
)

const (
	prefix  = "This is my prefix"
	postfix = "~~~~~"
)

type Config struct {
	prefix string
}

type ILogger interface {
	Sprint() string
}

type Logger struct {
	prefix string
}

func (s *Logger) Sprint() string {
	return s.prefix + postfix
}

func CreateLogger(name string) ILogger {
	return &Logger{name}
}

func TestDICreate(t *testing.T) {
	di := NewDepInjector()

	if err := di.Provide(func(prefix string) *Config {
		return &Config{prefix}
	}, false); err != nil {
		t.Fatal(err)
	}

	if err := di.Provide(func(cfg *Config) *Logger {
		return &Logger{cfg.prefix}
	}, false); err != nil {
		t.Fatal(err)
	}

	var logger *Logger
	if err := di.Create(&logger, prefix); err != nil {
		t.Fatal(err)
	}

	if logger.Sprint() != prefix+postfix {
		t.Fatal("logger output incorrect")
	}
	fmt.Println(logger.Sprint())
}

func TestDISingleton(t *testing.T) {
	di := NewDepInjector()

	var logger *Logger

	singleton := map[string]bool{}
	if err := di.Provide(func(s string) *Logger {
		if singleton[s] {
			t.Fatal("Duplicate")
		}
		fmt.Println("Create with string =", s)
		singleton[s] = true
		return &Logger{s}
	}, true); err != nil {
		t.Fatal(err)
	}

	if err := di.Create(&logger, "1"); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, "2"); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, "1"); err != nil {
		t.Fatal(err)
	}

	singleton2 := map[*Config]bool{}
	if err := di.Provide(func(config *Config) *Logger {
		if singleton2[config] {
			t.Fatal("Duplicate")
		}

		fmt.Printf("Create with config = %s, ptr = %p\n", config.prefix, config)
		singleton2[config] = true
		return &Logger{config.prefix}
	}, true); err != nil {
		t.Fatal(err)
	}

	config1 := &Config{"1"}
	config2 := &Config{"2"}
	if err := di.Create(&logger, config1); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, config2); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, config1); err != nil {
		t.Fatal(err)
	}

	singleton3 := map[int]bool{}
	if err := di.Provide(func(idx int) *Logger {
		if singleton3[idx] {
			t.Fatal("Duplicate")
		}

		fmt.Printf("Create with idx = %d, ptr = %p\n", idx, logger)
		singleton3[idx] = true
		return &Logger{fmt.Sprintf("%d", idx)}
	}, true); err != nil {
		t.Fatal(err)
	}

	if err := di.Create(&logger, 1); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, 2); err != nil {
		t.Fatal(err)
	}
	if err := di.Create(&logger, 1); err != nil {
		t.Fatal(err)
	}

	if err := di.Create(&logger, []string{"HAHA"}); err == nil {
		t.Fatal("[]string not support for singleton")
	} else {
		fmt.Println(err)
	}
}

func TestDICircularDependency(t *testing.T) {
	di := NewDepInjector()

	if err := di.Provide(func(int) int {
		return 1
	}, true); err != nil {
		t.Fatal(err)
	}

	var v int
	if err := di.Create(&v, 1); err != nil {
		t.Fatal(err)
	}

	if err := di.Create(&v); err == nil {
		t.Fatal(fmt.Errorf("should got circular dependency error"))
	} else {
		fmt.Println(err)
	}

	if err := di.Provide(func(*Config) *Logger {
		return &Logger{"Test"}
	}, false); err != nil {
		t.Fatal(err)
	}

	if err := di.Provide(func(*Logger) *Config {
		return &Config{"Test"}
	}, false); err != nil {
		t.Fatal(err)
	}

	var logger *Logger
	if err := di.Create(&logger); err == nil {
		t.Fatal(fmt.Errorf("should got circular dependency error"))
	} else {
		fmt.Println(err)
	}

	if err := di.Create(&logger, &Config{"test"}); err != nil {
		t.Fatal(err)
	}
}

func TestMultiParameters(t *testing.T) {
	di := NewDepInjector()

	if err := di.Provide(func() *Config {
		return &Config{"Test"}
	}, false); err != nil {
		t.Fatal(err)
	}

	if err := di.Provide(func(int, *Config, string) *Logger {
		return &Logger{"Test"}
	}, false); err != nil {
		t.Fatal(err)
	}

	var l *Logger
	if err := di.Create(&l); err == nil {
		t.Fatal("No enough parameters for provider, should fatal here")
	} else {
		fmt.Println(err)
	}

	if err := di.Create(&l, 1); err == nil {
		t.Fatal("No enough parameters for provider, should fatal here")
	} else {
		fmt.Println(err)
	}

	if err := di.Create(&l, 1, "hello"); err != nil {
		t.Fatal(err)
	}

	if err := di.Create(&l, "hello", 1); err != nil {
		t.Fatal(err)
	}
}

func TestDIDuplicateParameters(t *testing.T) {
	di := NewDepInjector()

	if err := di.Provide(func(int, int) string {
		return ""
	}, true); err == nil {
		t.Fatal(fmt.Errorf("the parameters of provider shouldn't have duplicate type"))
	} else {
		fmt.Println(err)
	}
}

func BenchmarkDICreate(b *testing.B) {
	di := NewDepInjector()

	if err := di.Provide(func(prefix string) *Config {
		return &Config{prefix}
	}, false); err != nil {
		b.Fatal(err)
	}

	if err := di.Provide(func(cfg *Config) *Logger {
		return &Logger{cfg.prefix}
	}, false); err != nil {
		b.Fatal(err)
	}

	var logger *Logger

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := di.Create(&logger, "This is my prefix"); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	if logger.Sprint() != prefix+postfix {
		b.Fatal("logger output incorrect")
	}
}

func BenchmarkNoDICreate(b *testing.B) {
	createConfig := func(prefix string) *Config {
		return &Config{prefix}
	}
	createLogger := func(cfg *Config) *Logger{
		return &Logger{cfg.prefix}
	}

	var logger *Logger

	b.StartTimer()
	for i:=0; i<b.N; i++ {
		logger = createLogger(createConfig("This is my prefix"))
	}
	b.StopTimer()

	if logger.Sprint() != prefix+postfix {
		b.Fatal("logger output incorrect")
	}
}

func BenchmarkDICreateSingleton(b *testing.B) {
	di := NewDepInjector()

	if err := di.Provide(func(config *Config) *Logger {
		return &Logger{config.prefix}
	}, true); err != nil {
		b.Fatal(err)
	}

	var logger *Logger
	config := &Config{"1"}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := di.Create(&logger, config); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
