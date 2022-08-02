package uri

import "testing"

func TestRegisterUri(t *testing.T) {
	parser := NewSimpleUriParser("db", "type").
		WithScheme("mongo").
		WithField(UriFieldUsername, "username", nil).
		WithField(UriFieldPassword, "password", nil).
		WithField(UriFieldHost, "addr", nil).
		WithField(UriFieldPath, "database_name", nil).
		WithQuery("test_param").
		WithQuery("test_param2").
		WithScheme("sqlite").
		WithField(UriFieldPath, "path", nil).
		WithQuery("test_param").
		WithQuery("test_param2")

	if err := RegisterUri(parser); err != nil {
		t.Fatal(err)
	}

	t.Log(parser)

	result, err := ParseUri("db", "sqlite:///c:\\test\\1.db?test_param=123&test_param2=1234&test_param=456")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	result, err = ParseUri("db", "sqlite:///tmp/1.db?test_param=123&test_param2=1234&test_param=456")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	result, err = ParseUri("db", "mongo://viecks:12345678@127.0.0.1:8001/test_db?test_param=123&test_param2=1234&test_param=456")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
