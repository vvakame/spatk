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
	buf         bytes.Buffer
	noIndent    bool
	indentLevel int
	lineHead    bool
	errCount    int

	currentStep step

	flagSelect  bool
	flagFrom    bool
	flagWhere   bool
	flagOrderBy bool
}

func (b *builder) writeToken(token string) {
	if b.lineHead {
		b.writeIndent()
	} else {
		b.buf.WriteString(" ")
	}
	b.buf.WriteString(token)
}

func (b *builder) writeTokenWithoutSpace(token string) {
	if b.lineHead {
		b.writeIndent()
	}
	b.buf.WriteString(token)
}

func (b *builder) incrementIndent() {
	b.indentLevel += 1
}
func (b *builder) decrementIndent() {
	if b.indentLevel > 0 {
		b.indentLevel -= 1
	}
}

func (b *builder) writeIndent() {
	if b.noIndent {
		return
	}
	b.buf.WriteString(strings.Repeat("  ", b.indentLevel))
	b.lineHead = false
}

func (b *builder) newLine() {
	if b.noIndent {
		b.buf.WriteString(" ")
		return
	}
	b.buf.WriteString("\n")
	b.lineHead = true
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
		b.writeToken("SELECT")
		b.incrementIndent()
	} else if b.currentStep > stepSelect {
		b.writeError("unexpected SELECT keyword")
	}

	return b
}

func (b *builder) Delete() DeleteBuilder {
	if b.currentStep == stepInit {
		b.currentStep = stepDelete
		b.writeToken("DELETE")
		b.incrementIndent()
	} else if b.currentStep > stepDelete {
		b.writeError("unexpected DELETE keyword")
	}

	return b
}

func (b *builder) AsStruct() SelectBuilder {
	if !b.flagSelect {
		b.writeToken("AS STRUCT")
	} else {
		b.writeError("unexpected AS STRUCT keyword")
	}

	return b
}

func (b *builder) C(name string, at ...string) SelectBuilder {
	if len(at) == 0 {
		b.CS(name)
	} else {
		b.CS(fmt.Sprintf("%s %s", name, at[0]))
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
		b.newLine()
	} else {
		b.writeTokenWithoutSpace(",")
		b.newLine()
	}

	b.writeToken(names[0])

	b.CS(names[1:]...)

	return b
}

func (b *builder) From() FromBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepDelete {
		b.currentStep = stepFrom
		b.newLine()
		b.decrementIndent()
		b.writeToken("FROM")
		b.incrementIndent()
	} else if b.currentStep != stepFrom {
		b.writeError("unexpected FROM keyword")
	}

	return b
}

func (b *builder) Name(tableName string, at ...string) FromBuilder {
	if !b.flagFrom {
		b.flagFrom = true
		b.newLine()
	} else {
		b.writeTokenWithoutSpace(",")
		b.newLine()
	}

	b.writeToken(tableName)
	if len(at) != 0 {
		b.writeToken(at[0])
	}

	if 2 <= len(at) {
		b.writeError(fmt.Sprintf("too many arguments: %v", at))
	}

	return b
}

func (b *builder) Where() WhereBuilder {
	if b.currentStep == stepFrom {
		b.currentStep = stepWhere
		b.newLine()
		b.decrementIndent()
		b.writeToken("WHERE")
		b.incrementIndent()
	} else if b.currentStep != stepWhere {
		b.writeError("unexpected WHERE keyword")
	}

	return b
}

func (b *builder) E(token ...string) WhereBuilder {
	if len(token) == 0 {
		return b
	}
	b.newLine()
	if b.flagWhere {
		b.writeToken("AND")
	}

	for _, t := range token {
		b.writeToken(t)
	}

	// TODO AND OR とかが末尾にないかチェックしたほうがよい
	b.flagWhere = true

	return b
}

func (b *builder) OrderBy() OrderByBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepFrom || b.currentStep == stepWhere {
		b.currentStep = stepOrderBy
		b.newLine()
		b.decrementIndent()
		b.writeToken("ORDER BY")
		b.incrementIndent()
	} else if b.currentStep != stepOrderBy {
		b.writeError("unexpected ORDER BY keyword")
	}

	return b
}

func (b *builder) O(token ...string) OrderByBuilder {
	if !b.flagOrderBy {
		b.flagOrderBy = true
		b.newLine()
	} else {
		b.writeTokenWithoutSpace(",")
		b.newLine()
	}

	for _, t := range token {
		b.writeToken(t)
	}

	return b
}

func (b *builder) Limit(limit string) VoidBuilder {
	if b.currentStep == stepSelect || b.currentStep == stepFrom || b.currentStep == stepWhere || b.currentStep == stepOrderBy {
		b.currentStep = stepLimit
		b.newLine()
		b.decrementIndent()
		b.writeToken("LIMIT")
	} else {
		b.writeError("unexpected LIMIT keyword")
	}

	b.writeToken(limit)

	return b
}
