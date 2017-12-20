package urn_test

import (
	"fmt"
	"testing"

	"github.com/olivoil/pkg/urn"
	"github.com/stretchr/testify/assert"
)

func TestURN(t *testing.T) {
	t.Run("IsNil", func(t *testing.T) {
		t.Run("it returns true for a Nil value", func(t *testing.T) {
			assert.True(t, urn.Nil.IsNil())
		})

		t.Run("it returns false for other values", func(t *testing.T) {
			u := urn.New("myApp", "myService", "myResource", "myID")
			assert.False(t, u.IsNil())
		})
	})

	t.Run("Parse", func(t *testing.T) {
		testCases := []struct {
			String string
			URN    urn.URN
			Err    error
		}{
			{`myapp://myservice/myresource/myid`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "myid"}, nil},
			{`myapp://myservice/myresource/*`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "*"}, nil},
			{`myApp://myService/myResource/myID`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "myid"}, nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("it parses %s", tc.String), func(t *testing.T) {
				u, err := urn.Parse(tc.String)
				assert.Equal(t, err, tc.Err)
				assert.Equal(t, u, tc.URN)
			})
		}
	})

	t.Run("MarshalText", func(t *testing.T) {
		testCases := []struct {
			String string
			URN    urn.URN
			Err    error
		}{
			{`myapp://myservice/myresource/myid`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "myid"}, nil},
			{`myapp://myservice/myresource/*`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "*"}, nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("it marshals %s", tc.String), func(t *testing.T) {
				b, err := tc.URN.MarshalText()
				assert.Equal(t, err, tc.Err)
				assert.Equal(t, string(b), tc.String)
			})
		}
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		testCases := []struct {
			String string
			URN    urn.URN
			Err    error
		}{
			{`myapp://myservice/myresource/myid`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "myid"}, nil},
			{`myapp://myservice/myresource/*`, urn.URN{Domain: "myapp", Service: "myservice", Name: "myresource", ID: "*"}, nil},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("it unmarshals %v", tc.URN), func(t *testing.T) {
				u := urn.URN{}
				err := u.UnmarshalText([]byte(tc.String))
				assert.Equal(t, err, tc.Err)
				assert.Equal(t, u.String(), tc.String)
			})
		}
	})
}
