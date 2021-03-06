package sqb

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var _ SelectBuilder = (*builder)(nil)
var _ FromBuilder = (*builder)(nil)
var _ WhereBuilder = (*builder)(nil)
var _ OrderByBuilder = (*builder)(nil)

type step int

const (
	stepInit step = iota
	stepSelect
	stepDelete
	stepFrom
	stepWhere
	stepOrderBy
	stepLimit
)

type builder struct {
	buf      bytes.Buffer
	errCount int

	currentStep step

	flagSelect  bool
	flagFrom    bool
	flagWhere   bool
	flagOrderBy bool
}

func (b *builder) writeError(msg string) {
	b.errCount++
	b.buf.WriteString(" !ERR")
	b.buf.WriteString(strconv.Itoa(b.errCount))
	b.buf.WriteString(":`")
	b.buf.WriteString(msg)
	b.buf.WriteString("`!")
}

func (b *builder) Build() (string, error) {
	if b.errCount != 0 {
		return "", fmt.Errorf("%d error(s) occured! %s", b.errCount, b.buf.String())
	}

	return b.buf.String(), nil
}

func (b *builder) Select() SelectBuilder {
	if b.currentStep == stepInit {
		b.currentStep = stepSelect
		b.buf.WriteString("SELECT")
	} else if b.currentStep > stepSelect {
		b.writeError("unexpected SELECT keyword")
	}

	return b
}

func (b *builder) Delete() DeleteBuilder {
	if b.currentStep == stepInit {
		b.currentStep = stepDelete
		b.buf.WriteString("DELETE")
	} else if b.currentStep > stepDelete {
		b.writeError("unexpected DELETE keyword")
	}

	return b
}

func (b *builder) AsStruct() SelectBuilder {
	if !b.flagSelect {
		b.buf.WriteString(" AS STRUCT")
	} else {
		b.writeError("unexpected AS STRUCT keyword")
	}

	return b
}

func (b *builder) C(name string, at ...string) SelectBuilder {
	if !b.flagSelect {
		b.flagSelect = true
	} else {
		b.buf.WriteString(",")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(name)
	if len(at) != 0 {
		b.buf.WriteString(" ")
		b.buf.WriteString(at[0])
	}

	if 2 <= len(at) {
		b.writeError(fmt.Sprintf("too many arguments: %v", at))
	}

	return b
}

func (b *builder) CS(names ...string) SelectBuilder {
	if len(names) == 0 {
		return b
	}

	if !b.flagSelect {
		b.flagSelect = true
	} else {
		b.buf.WriteString(",")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(strings.Join(names, ", "))

	return b
}

func (b *builder) From() FromBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepDelete {
		b.currentStep = stepFrom
		b.buf.WriteString(" FROM")
	} else if b.currentStep != stepFrom {
		b.writeError("unexpected FROM keyword")
	}

	return b
}

func (b *builder) Name(tableName string, at ...string) FromBuilder {
	if !b.flagFrom {
		b.flagFrom = true
	} else {
		b.buf.WriteString(",")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(tableName)
	if len(at) != 0 {
		b.buf.WriteString(" ")
		b.buf.WriteString(at[0])
	}

	if 2 <= len(at) {
		b.writeError(fmt.Sprintf("too many arguments: %v", at))
	}

	return b
}

func (b *builder) Where() WhereBuilder {
	if b.currentStep == stepFrom {
		b.currentStep = stepWhere
		b.buf.WriteString(" WHERE")
	} else if b.currentStep != stepWhere {
		b.writeError("unexpected WHERE keyword")
	}

	return b
}

func (b *builder) E(token ...string) WhereBuilder {
	if len(token) == 0 {
		return b
	}
	if b.flagWhere {
		b.buf.WriteString(" AND")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(strings.Join(token, " "))

	// TODO AND OR ????????????????????????????????????????????????????????????
	b.flagWhere = true

	return b
}

func (b *builder) OrderBy() OrderByBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepFrom || b.currentStep == stepWhere {
		b.currentStep = stepOrderBy
		b.buf.WriteString(" ORDER BY")
	} else if b.currentStep != stepOrderBy {
		b.writeError("unexpected ORDER BY keyword")
	}

	return b
}

func (b *builder) O(token ...string) OrderByBuilder {
	if !b.flagOrderBy {
		b.flagOrderBy = true
	} else {
		b.buf.WriteString(",")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(strings.Join(token, " "))

	return b
}

func (b *builder) Limit(limit string) VoidBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepFrom || b.currentStep == stepWhere || b.currentStep == stepOrderBy {
		b.currentStep = stepLimit
		b.buf.WriteString(" LIMIT")
	} else {
		b.writeError("unexpected LIMIT keyword")
	}

	b.buf.WriteString(" ")
	b.buf.WriteString(limit)

	return b
}
