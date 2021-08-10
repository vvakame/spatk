package scur

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
)

type cursorParameterType int

const (
	cursorParameterTypeUndefined cursorParameterType = iota
	cursorParameterTypeString
	cursorParameterTypeInt
	cursorParameterTypeInt64
	cursorParameterTypeBool
	cursorParameterTypeFloat64
	cursorParameterTypeTime
)

func encodeParameter(w io.StringWriter, v interface{}) error {
	// NOTE : をdelimiterに使うのでこれが出現しないようにしないようにやる必要がある
	switch v := v.(type) {
	case string:
		_, _ = w.WriteString(fmt.Sprintf("%02x", cursorParameterTypeString))
		_, _ = w.WriteString(base64.RawURLEncoding.EncodeToString([]byte(v)))
	case int:
		_, _ = w.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeInt, strconv.FormatInt(int64(v), 36)))
	case int64:
		_, _ = w.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeInt64, strconv.FormatInt(v, 36)))
	case bool:
		var nv int
		if v {
			nv = 1
		} else {
			nv = 0
		}
		_, _ = w.WriteString(fmt.Sprintf("%02x%d", cursorParameterTypeBool, nv))
	case float64:
		_, _ = w.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeFloat64, strconv.FormatUint(math.Float64bits(v), 36)))
	case time.Time:
		_, _ = w.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeTime, strconv.FormatInt(v.UnixNano(), 36)))
	case spanner.Encoder:
		v2, err := v.EncodeSpanner()
		if err != nil {
			return err
		}
		err = encodeParameter(w, v2)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported cursor value type: %T", v)
	}

	return nil
}

func decodeParameter(s string) (interface{}, error) {
	if len(s) < 2 {
		return nil, fmt.Errorf("unexpected encoded value: %s", s)
	}
	n, err := strconv.ParseInt(s[0:2], 16, 64)
	if err != nil {
		return nil, err
	}
	sv := s[2:]
	switch cursorParameterType(n) {
	case cursorParameterTypeString:
		b, err := base64.RawURLEncoding.DecodeString(sv)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case cursorParameterTypeInt:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return nil, fmt.Errorf("int value decode failed: %w", err)
		}
		return int(v), nil
	case cursorParameterTypeInt64:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return nil, fmt.Errorf("int64 value decode failed: %w", err)
		}
		return v, nil
	case cursorParameterTypeBool:
		v, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bool value decode failed: %w", err)
		}
		return v == 1, nil
	case cursorParameterTypeFloat64:
		v, err := strconv.ParseUint(sv, 36, 64)
		if err != nil {
			return nil, fmt.Errorf("float64 value decode failed: %w", err)
		}
		return math.Float64frombits(v), nil
	case cursorParameterTypeTime:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return nil, fmt.Errorf("time.Time value decode failed: %w", err)
		}
		return time.Unix(0, v), nil
	default:
		return nil, fmt.Errorf("unsupported cursor value type: %d", n)
	}
}
