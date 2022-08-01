package di

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"
	"sync"
)

type _DepInjector struct {
	providers []*_Provider
	mu        sync.RWMutex
}

type _Provider struct {
	invoker       interface{}
	provideType   reflect.Type
	withError     bool
	singleton     bool
	mapSingletonV map[uint64]interface{}
	mu            sync.RWMutex
}

type DepInjector interface {
	Provide(resolver interface{}, singleton bool) error
	Invoke(invoker interface{}, parameters ...interface{}) error
	Create(ptr interface{}, parameters ...interface{}) error
	Clear()
}

var GlobalDepInjector = NewDepInjector()

func NewDepInjector() DepInjector {
	return &_DepInjector{}
}

func (s *_DepInjector) Provide(provider interface{}, singleton bool) error {
	providerType := reflect.TypeOf(provider)
	providerOut, err := s.checkProvider(providerType)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers = append(s.providers, &_Provider{
		invoker:       provider,
		provideType:   providerOut,
		withError:     providerType.NumOut() == 2,
		singleton:     singleton,
		mapSingletonV: make(map[uint64]interface{}),
	})

	return nil
}

func (s *_DepInjector) checkProvider(providerType reflect.Type) (reflect.Type, error) {
	if providerType == nil || providerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("the provider must be a function object")
	}

	switch providerType.NumOut() {
	case 1:
	case 2:
		if providerType.Out(1).Kind() != reflect.Interface || !providerType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return nil, fmt.Errorf("the provider must returns an instance or instance with an error")
		}
	default:
		return nil, fmt.Errorf("the provider must returns an instance or instance with an error")
	}

	typeSet := map[reflect.Type]struct{}{}
	for i := 0; i < providerType.NumIn(); i++ {
		if _, ok := typeSet[providerType.In(i)]; ok {
			return nil, fmt.Errorf("the parameters of provider can not have duplicate type %v", providerType.In(i))
		}
		typeSet[providerType.In(i)] = struct{}{}
	}

	return providerType.Out(0), nil
}

func (s *_DepInjector) Invoke(invoker interface{}, parameters ...interface{}) error {
	invokerType := reflect.TypeOf(invoker)
	if invokerType == nil || invokerType.Kind() != reflect.Func {
		return fmt.Errorf("cannot detect type of the invoker, make sure your are passing function object")
	}

	switch invokerType.NumOut() {
	case 0:
		arguments, err := s.arguments(invoker, parameters, map[*_Provider]struct{}{}, []*_Provider{})
		if err != nil {
			return err
		}

		reflect.ValueOf(invoker).Call(arguments)
		return nil
	case 1:
		if invokerType.Out(0).Kind() != reflect.Interface || !invokerType.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("the invoker should returns error or nothing")
		}

		arguments, err := s.arguments(invoker, parameters, map[*_Provider]struct{}{}, []*_Provider{})
		if err != nil {
			return err
		}

		resp := reflect.ValueOf(invoker).Call(arguments)[0].Interface()
		if resp == nil {
			return nil
		} else {
			return resp.(error)
		}
	default:
		return fmt.Errorf("the invoker should returns error or nothing")
	}
}

func (s *_DepInjector) Create(ptr interface{}, parameters ...interface{}) error {
	ptrType := reflect.TypeOf(ptr)
	if ptrType == nil || ptrType.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot detect type of the object-ptr, make sure your are passing reference of the object")
	}

	objType := ptrType.Elem()

	var lastErr error

	s.mu.RLock()
	providers := s.providers
	s.mu.RUnlock()

	for _, provider := range providers {
		if provider.provideType == objType || (objType.Kind() == reflect.Interface && provider.provideType.Implements(objType)) {
			instance, err := s.invokeProvider(provider, parameters, map[*_Provider]struct{}{}, []*_Provider{})
			if err == nil {
				if instance == nil {
					reflect.ValueOf(ptr).Elem().Set(reflect.Zero(objType))
				} else {
					reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(instance))
				}
				return nil
			} else {
				lastErr = err
			}
		}
	}
	if lastErr != nil {
		return lastErr
	}

	return fmt.Errorf("no provider found for the type: " + objType.String())
}

func (s *_DepInjector) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers = nil
}

func (s *_DepInjector) arguments(function interface{}, parameters []interface{}, providerSet map[*_Provider]struct{}, providerList []*_Provider) ([]reflect.Value, error) {
	argumentsCount := reflect.TypeOf(function).NumIn()
	arguments := make([]reflect.Value, argumentsCount)

Out:
	for i := 0; i < argumentsCount; i++ {
		objType := reflect.TypeOf(function).In(i)
		for _, parameter := range parameters {
			if objType == reflect.TypeOf(parameter) || (objType.Kind() == reflect.Interface && reflect.TypeOf(parameter).Implements(objType)) {
				arguments[i] = reflect.ValueOf(parameter)
				continue Out
			}
		}

		var lastErr error

		s.mu.RLock()
		providers := s.providers
		s.mu.RUnlock()

		for _, provider := range providers {
			if provider.provideType == objType || (objType.Kind() == reflect.Interface && provider.provideType.Implements(objType)) {
				instance, err := s.invokeProvider(provider, parameters, providerSet, providerList)
				if err == nil {
					if instance == nil {
						arguments[i] = reflect.New(objType).Elem()
					} else {
						arguments[i] = reflect.ValueOf(instance)
					}
					continue Out
				} else {
					lastErr = err
				}
			}
		}
		if lastErr != nil {
			return nil, lastErr
		}

		return nil, fmt.Errorf("no provider found for the type: " + objType.String())
	}

	return arguments, nil
}

func (s *_DepInjector) invokeProvider(provider *_Provider, parameters []interface{}, providerSet map[*_Provider]struct{}, providerList []*_Provider) (v interface{}, err error) {
	if _, ok := providerSet[provider]; ok {
		var chain []string
		for _, p := range providerList {
			chain = append(chain, p.provideType.String())
		}
		chain = append(chain, provider.provideType.String())

		return nil, fmt.Errorf("detect circular dependency of type %v, dependency chain: %s", provider.provideType, strings.Join(chain, " -> "))
	}

	providerSet[provider] = struct{}{}
	providerList = append(providerList, provider)
	defer func() {
		delete(providerSet, provider)
		providerList = providerList[:len(providerList)-1]
	}()

	var hash uint64
	var hashOk bool
	var hashT reflect.Type
	if provider.singleton {
		hash, hashOk, hashT = parameterHash(parameters)
		if !hashOk {
			return nil, fmt.Errorf("can not calculate the hash of type %v which appears in parameters for singleton provider", hashT)
		}
	}

	defer func() {
		if provider.singleton && err == nil && hashOk {
			provider.mu.Lock()
			provider.mapSingletonV[hash] = v
			provider.mu.Unlock()
		}
	}()

	if provider.singleton && hashOk {
		provider.mu.RLock()
		if singleton, ok := provider.mapSingletonV[hash]; ok {
			provider.mu.RUnlock()
			return singleton, nil
		}
		provider.mu.RUnlock()
	}

	args, err := s.arguments(provider.invoker, parameters, providerSet, providerList)
	if err != nil {
		return nil, err
	}

	resp := reflect.ValueOf(provider.invoker).Call(args)
	if !provider.withError || resp[1].Interface() == nil {
		return resp[0].Interface(), nil
	} else {
		return resp[0].Interface(), resp[1].Interface().(error)
	}
}

func parameterHash(parameters []interface{}) (uint64, bool, reflect.Type) {
	if len(parameters) == 0 {
		return 0, true, nil
	}

	hasher := fnv.New64()
	for _, parameter := range parameters {
		v := reflect.ValueOf(parameter)

		for {
			if v.Kind() == reflect.Interface {
				v = v.Elem()
				continue
			}

			break
		}

		if !v.IsValid() {
			v = reflect.Zero(reflect.TypeOf(0))
		}

		switch v.Kind() {
		case reflect.Int:
			v = reflect.ValueOf(v.Int())
		case reflect.Uint:
			v = reflect.ValueOf(v.Uint())
		case reflect.Bool:
			var tmp int8
			if v.Bool() {
				tmp = 1
			}
			v = reflect.ValueOf(tmp)
		}

		k := v.Kind()

		if k >= reflect.Int && k <= reflect.Complex64 {
			if err := binary.Write(hasher, binary.LittleEndian, v.Interface()); err == nil {
				continue
			}
		} else if k == reflect.String {
			if _, err := hasher.Write([]byte(v.String())); err == nil {
				continue
			}
		} else if k == reflect.Ptr {
			if err := binary.Write(hasher, binary.LittleEndian, uint64(v.Pointer())); err == nil {
				continue
			}
		}

		return 0, false, v.Type()
	}

	return hasher.Sum64(), true, nil
}
