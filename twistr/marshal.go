package twistr

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(c interface{}) ([]byte, error) {
	// Indirect can always be used; if the value is not a pointer, it just
	// returns the value.
	cv := reflect.Indirect(reflect.ValueOf(c))
	var field reflect.Value
	buf := new(bytes.Buffer)
	n := cv.NumField()
	for i := 0; i < n; i++ {
		field = cv.Field(i)
		if field.Type().Kind() == reflect.Slice {
			if err := marshalSlice(field, buf); err != nil {
				return []byte{}, err
			}
		} else {
			if err := marshalValue(field, buf); err != nil {
				return []byte{}, err
			}
		}
		if i < (n - 1) {
			buf.WriteString(" ")
		}
	}
	return buf.Bytes(), nil
}

func marshalSlice(field reflect.Value, buf *bytes.Buffer) error {
	// Writes "[ el1 el2 el3 ]". No leading or trailing spaces.
	var marshalFn func(field reflect.Value) string
	switch fieldKind(field.Type().Elem()) {
	case "country":
		marshalFn = marshalCountryPtr
	case "card":
		marshalFn = marshalCard
	default:
		return fmt.Errorf("Unsupported field '%s'", field.Type().Elem().Name())
	}
	buf.WriteString("[ ")
	n := field.Len()
	for i := 0; i < n; i++ {
		val := field.Index(i)
		buf.WriteString(marshalFn(val))
		buf.WriteString(" ")
	}
	buf.WriteString("]")
	return nil
}

func marshalCountryPtr(field reflect.Value) string {
	country := reflect.Indirect(field)
	return strings.ToLower(country.FieldByName("Name").String())
}

func marshalCard(field reflect.Value) string {
	card := reflect.Indirect(field)
	return strings.ToLower(card.FieldByName("Name").String())
}

func marshalValue(field reflect.Value, buf *bytes.Buffer) error {
	switch fieldKind(field.Type()) {
	case "int":
		buf.WriteString(strconv.Itoa(int(field.Int())))
	case "country":
		buf.WriteString(marshalCountryPtr(field))
	case "card":
		buf.WriteString(marshalCard(field))
	case "aff":
		buf.WriteString(strings.ToLower(Aff(field.Int()).String()))
	case "playkind":
		buf.WriteString(strconv.Itoa(int(field.Int())))
	case "opskind":
		buf.WriteString(strconv.Itoa(int(field.Int())))
	default:
		return fmt.Errorf("Unknown field '%s'", field.Type().Name())
	}
	return nil
}

func Unmarshal(line string, c interface{}) (err error) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	// Value of the struct that c points to
	cv := reflect.Indirect(reflect.ValueOf(c))
	var field reflect.Value
	var word string
	for i := 0; i < cv.NumField(); i++ {
		if !scanner.Scan() {
			return fmt.Errorf("Not enough tokens for %s", cv.Type().Name())
		}
		word = scanner.Text()
		field = cv.Field(i)
		if field.Type().Kind() == reflect.Slice {
			if word != "[" {
				return errors.New("Malformed list input. Expected '['")
			}
			if err = unmarshalSlice(scanner, field); err != nil {
				return
			}
		} else {
			if err = unmarshalValue(word, field); err != nil {
				return
			}
		}
	}
	return
}

func unmarshalSlice(scanner *bufio.Scanner, field reflect.Value) (err error) {
	var words []string
	if words, err = readSlice(scanner); err != nil {
		return err
	}
	switch fieldKind(field.Type().Elem()) {
	case "country":
		val := make([]*Country, len(words))
		for i, word := range words {
			if val[i], err = lookupCountry(word); err != nil {
				return
			}
		}
		field.Set(reflect.ValueOf(val))
	case "card":
		val := make([]Card, len(words))
		for i, word := range words {
			if val[i], err = lookupCard(word); err != nil {
				return
			}
		}
		field.Set(reflect.ValueOf(val))
	}
	return
}

func readSlice(scanner *bufio.Scanner) ([]string, error) {
	s := []string{}
	var word string
	for scanner.Scan() {
		word = scanner.Text()
		if word == "]" {
			return s, scanner.Err()
		}
		s = append(s, word)
	}
	return s, errors.New("Did not encounter ending ']' of list")
}

func unmarshalValue(word string, field reflect.Value) (err error) {
	switch fieldKind(field.Type()) {
	case "int":
		var num int
		if num, err = strconv.Atoi(word); err != nil {
			return err
		}
		field.SetInt(int64(num))
	case "country":
		var country *Country
		if country, err = lookupCountry(word); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(country))
	case "card":
		var card Card
		if card, err = lookupCard(word); err != nil {
			return err
		}
		field.Set(reflect.ValueOf(card))
	case "aff":
		var aff Aff
		if aff, err = lookupAff(word); err != nil {
			return err
		}
		field.SetInt(int64(aff))
	case "playkind":
		var pk PlayKind
		if pk, err = lookupPlayKind(word); err != nil {
			return err
		}
		field.SetInt(int64(pk))
	case "opskind":
		var ok OpsKind
		if ok, err = lookupOpsKind(word); err != nil {
			return err
		}
		field.SetInt(int64(ok))
	}
	return
}

func fieldKind(ftype reflect.Type) string {
	kind := ftype.Kind()
	name := ftype.Name()
	switch {
	case name == "string":
		return "string"
	case name == "int":
		return "int"
	case kind == reflect.Ptr && ftype.Elem().Name() == "Country":
		return "country"
	case name == "Card":
		return "card"
	case name == "Aff":
		return "aff"
	case name == "PlayKind":
		return "playkind"
	case name == "OpsKind":
		return "opskind"
	default:
		return "?"
	}
}
