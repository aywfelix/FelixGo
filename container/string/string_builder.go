package container

import (
	"fmt"
	"strconv"
	"strings"
)

type StringBuilder struct {
	strings.Builder
}

func (sb *StringBuilder) Append(obj interface{}) (n int, err error) {
	switch val := obj.(type) {
	case nil:
		return sb.WriteString("[nil]")
	case byte: // == uint8
		return 1, sb.WriteByte(val)
	case []byte:
		return sb.Write(val)
	case rune: // == int32
		return sb.WriteRune(val)
	case string:
		return sb.WriteString(val)
	case float32:
		return sb.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 64))
	case float64:
		return sb.WriteString(strconv.FormatFloat(val, 'f', -1, 64))
	case int:
		return sb.WriteString(strconv.Itoa(val))
	case int8:
		return sb.WriteString(strconv.FormatInt(int64(val), 10))
	case int16:
		return sb.WriteString(strconv.FormatInt(int64(val), 10))
	case int64:
		return sb.WriteString(strconv.FormatInt(val, 10))
	case uint:
		return sb.WriteString(strconv.Itoa(int(val)))
	case uint16:
		return sb.WriteString(strconv.FormatInt(int64(val), 10))
	case uint32:
		return sb.WriteString(strconv.FormatInt(int64(val), 10))
	case uint64:
		return sb.WriteString(strconv.FormatInt(int64(val), 10))
	default:
		return sb.WriteString(fmt.Sprint(val))
	}
}

func (sb *StringBuilder) AppendArgs(objs ...interface{}) (n int, err error) {
	for _, obj := range objs {
		if num, err := sb.Append(obj); err != nil {
			return n, err
		} else {
			n += num
		}
	}

	return n, nil
}

func (sb *StringBuilder) AppendLine(obj interface{}) (n int, err error) {
	if n, err = sb.Append(obj); err != nil {
		return n, err
	}

	if err = sb.WriteByte('\n'); err != nil {
		return n, err
	} else {
		n++
	}

	return n, nil
}

func (sb *StringBuilder) AppendFormat(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(sb, format, args...)
}

func (sb *StringBuilder) String() string {
	return sb.Builder.String()
}
