package validate

import (
	"fmt"
	"strconv"
	"time"
)

// normalize between go values and validate-lang value
func normalize(i interface{}) (interface{}, error) {
	switch i.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
		return toInteger(i)
	case float32, float64:
		return toNumber(i)
	}

	return i, nil
}

// toString converts `i` to a string.
func toString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

// toNumber converts `i` to a float64.
func toNumber(i interface{}) (float64, error) {
	res, err := strconv.ParseFloat(toString(i), 64)
	if err != nil {
		res = 0.0
	}
	return res, err
}

// toNumbers converts `s` to a []float64.
func toNumbers(s []interface{}) (numbers []float64, err error) {
	for _, i := range s {
		n, e := toNumber(i)
		if e != nil {
			err = e
		}

		numbers = append(numbers, n)
	}

	return
}

// toInteger converts `i` to an integer.
func toInteger(i interface{}) (int64, error) {
	res, err := strconv.ParseInt(toString(i), 0, 64)
	if err != nil {
		res = 0
	}
	return res, err
}

// toIntegers converts `s` to a []int64.
func toIntegers(s []interface{}) (integers []int64, err error) {
	for _, i := range s {
		integer, e := toInteger(i)
		if e != nil {
			err = e
		}
		integers = append(integers, integer)
	}

	return
}

// toDuration converts `i` to a duration.
func toDuration(i interface{}) (time.Duration, error) {
	d, err := toInteger(i)
	return time.Duration(d), err
}

// toDurations converts `s` to a []duration.
func toDurations(s []interface{}) (durations []time.Duration, err error) {
	for _, i := range s {
		d, e := toDuration(i)
		if e != nil {
			err = e
		}
		durations = append(durations, d)
	}

	return
}
