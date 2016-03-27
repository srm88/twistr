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

// Logs must be referenced in this set, or we will not marshal them correctly.
var logTypes map[string]bool = map[string]bool{
	"CoupLog":    true,
	"RealignLog": true,
}

func Marshal(c interface{}) ([]byte, error) {
	// Indirect can always be used; if the value is not a pointer, it just
	// returns the value.
	cv := reflect.Indirect(reflect.ValueOf(c))
	buf := new(bytes.Buffer)
	var err error
	if isLog(cv.Type()) {
		err = marshalLog(cv, buf)
	} else {
		err = marshalValue(cv, buf)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func marshalLog(v reflect.Value, buf *bytes.Buffer) error {
	var field reflect.Value
	n := v.NumField()
	for i := 0; i < n; i++ {
		field = v.Field(i)
		if err := marshalValue(field, buf); err != nil {
			return err
		}
		if i < (n - 1) {
			buf.WriteString(" ")
		}
	}
	return nil
}

func marshalSlice(field reflect.Value, buf *bytes.Buffer) error {
	// Writes "[ el1 el2 el3 ]". No leading or trailing spaces.
	var marshalFn func(field reflect.Value) string
	switch valueKind(field.Type().Elem()) {
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
	if field.Type().Kind() == reflect.Slice {
		return marshalSlice(field, buf)
	}
	switch valueKind(field.Type()) {
	case "string":
		buf.WriteString(field.String())
	case "int":
		buf.WriteString(strconv.Itoa(int(field.Int())))
	case "country":
		buf.WriteString(marshalCountryPtr(field))
	case "card":
		buf.WriteString(marshalCard(field))
	case "aff":
		buf.WriteString(strings.ToLower(Aff(field.Int()).String()))
	case "playkind":
		buf.WriteString(strings.ToLower(PlayKind(field.Int()).String()))
	case "opskind":
		buf.WriteString(strings.ToLower(OpsKind(field.Int()).String()))
	default:
		return fmt.Errorf("Unknown field '%s'", field.Type().Name())
	}
	return nil
}

func Unmarshal(line string, c interface{}) (err error) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)
	// Value of c, dereferencing one pointer if necessary
	cv := reflect.Indirect(reflect.ValueOf(c))
	if isLog(cv.Type()) {
		err = unmarshalLog(scanner, cv)
	} else {
		err = unmarshalValue(scanner, cv)
	}
	return
}

func unmarshalLog(scanner *bufio.Scanner, cv reflect.Value) (err error) {
	for i := 0; i < cv.NumField(); i++ {
		if err = unmarshalValue(scanner, cv.Field(i)); err != nil {
			return
		}
	}
	return
}

func unmarshalSlice(scanner *bufio.Scanner, field reflect.Value) (err error) {
	var words []string
	if words, err = readSlice(scanner); err != nil {
		return err
	}
	switch valueKind(field.Type().Elem()) {
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

func unmarshalValue(scanner *bufio.Scanner, v reflect.Value) (err error) {
	if !scanner.Scan() {
		return fmt.Errorf("Not enough tokens for %s", v.Type().Name())
	}
	word := scanner.Text()
	if v.Type().Kind() == reflect.Slice {
		if word != "[" {
			return errors.New("Malformed list input. Expected '['")
		}
		if err = unmarshalSlice(scanner, v); err != nil {
			return
		}
	} else {
		if err = unmarshalWord(word, v); err != nil {
			return
		}
	}
	return
}

func unmarshalWord(word string, v reflect.Value) (err error) {
	switch valueKind(v.Type()) {
	case "string":
		v.SetString(word)
	case "int":
		var num int
		if num, err = strconv.Atoi(word); err != nil {
			return err
		}
		v.SetInt(int64(num))
	case "country":
		var country *Country
		if country, err = lookupCountry(word); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(country))
	case "card":
		var card Card
		if card, err = lookupCard(word); err != nil {
			return err
		}
		v.Set(reflect.ValueOf(card))
	case "aff":
		var aff Aff
		if aff, err = lookupAff(word); err != nil {
			return err
		}
		v.SetInt(int64(aff))
	case "playkind":
		var pk PlayKind
		if pk, err = lookupPlayKind(word); err != nil {
			return err
		}
		v.SetInt(int64(pk))
	case "opskind":
		var ok OpsKind
		if ok, err = lookupOpsKind(word); err != nil {
			return err
		}
		v.SetInt(int64(ok))
	}
	return
}

func isLog(t reflect.Type) bool {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return isLog(t.Elem())
	}
	return kind == reflect.Struct && logTypes[t.Name()]
}

func valueKind(vtype reflect.Type) string {
	if vtype.Kind() == reflect.Ptr {
		return valueKind(vtype.Elem())
	}
	switch vtype.Name() {
	case "string":
		return "string"
	case "int":
		return "int"
	case "Country":
		return "country"
	case "Card":
		return "card"
	case "Aff":
		return "aff"
	case "PlayKind":
		return "playkind"
	case "OpsKind":
		return "opskind"
	default:
		return "?"
	}
}
