package parser

import (
	"fmt"
	"strings"

	"github.com/templatesflow/cli/internal/token"
)

type Parser struct {
	tokens  []token.Token
	pos     int
	nodes   []Node
	errors  []error
	current token.Token
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    -1,
	}

	p.next()
	return p
}

func (p *Parser) Parse() ([]Node, []error) {
	for p.pos < len(p.tokens) {
		node := p.parseNode()
		if node != nil {
			p.nodes = append(p.nodes, node)
		} else {
			// Avoid infinite loop on errors, consume the token
			p.next()
		}
	}
	return p.nodes, p.errors
}

func (p *Parser) errorf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	p.errors = append(p.errors, err)
}

func (p *Parser) next() {
	p.pos++
	p.current = p.getCurrent()
}

func (p *Parser) consumeWhitespaces() string {
	var builder strings.Builder
	for p.current.Typ == token.WS {
		builder.WriteString(p.current.Val)
		p.next()
	}
	return builder.String()
}

func (p *Parser) getCurrent() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}

	return token.Token{Typ: token.EOF}
}

func (p *Parser) parseNode() Node {
	switch p.current.Typ {
	case token.TEXT:
		return p.parseText()
	case token.LEXPR:
		return p.parseExprBlock()
	case token.EOF:
		return nil // End of input
	default:
		p.errorf("unexpected token: %v", p.current)
		return nil
	}
}

func (p *Parser) parseText() Node {
	text := Text{
		Pos: p.current.Pos,
		Val: p.current.Val,
	}
	p.next()
	return text
}

func (p *Parser) parseExprBlock() Node {
	exprBlock := ExprBlock{
		LBrace: p.current.Pos,
	}
	p.next() // Consume LEXPR
	exprBlock.PostLWS = p.consumeWhitespaces()

	exprBlock.Body = p.parseExpr()

	if p.current.Typ != token.REXPR {
		p.errorf("expected REXPR, got %v", p.current)
		return exprBlock // Still return the partial ExprBlock
	}
	exprBlock.RBrace = p.current.Pos
	p.next() // Consume REXPR
	return exprBlock
}

func (p *Parser) parseExpr() Node {
	return p.parseBinaryExpr(1)
}

// parseBinaryExpr parses expressions with operator precedence.
func (p *Parser) parseBinaryExpr(minPrecedence int) Node {
	left := p.parsePrimary()

	for {
		opPrecedence, isRightAssoc := getPrecedence(p.current.Typ)
		if opPrecedence < minPrecedence {
			break
		}

		op := p.current
		p.next()

		ws := p.consumeWhitespaces()

		nextMinPrecedence := opPrecedence
		if !isRightAssoc {
			nextMinPrecedence++ // Left-associative operators require higher precedence for right operand
		}

		right := p.parseBinaryExpr(nextMinPrecedence)

		left = &BinaryExpr{
			X:        left,
			OpPos:    op.Pos,
			PostOpWS: ws,
			Op:       op.Typ,
			Y:        right,
		}
	}
	return left
}

// parsePrimary handles literals, identifiers, and parenthesized expressions.
func (p *Parser) parsePrimary() Node {
	switch p.current.Typ {
	case token.IDENT:
		ident := Ident{
			Pos:  p.current.Pos,
			Name: p.current.Val,
		}
		p.next()
		ident.PostWS = p.consumeWhitespaces()
		return ident
	case token.INT, token.FLOAT:
		lit := Lit{
			Pos: p.current.Pos,
			Val: p.current.Val,
			Typ: p.current.Typ,
		}
		p.next()
		lit.PostWS = p.consumeWhitespaces()
		return lit
	case token.LPAREN:
		p.next() // Consume '('
		expr := p.parseExpr()
		if p.current.Typ != token.RPAREN {
			p.errorf("expected closing ')', got %v", p.current)
			return expr // Return partial expression
		}
		p.next() // Consume ')'
		p.consumeWhitespaces()
		return expr
	default:
		p.errorf("expected identifier, literal, or '(', got %v", p.current)
		return nil
	}
}

// getPrecedence returns the precedence and associativity of an operator.
func getPrecedence(op token.Type) (int, bool) {
	switch op {
	case token.ADD, token.SUB:
		return 1, false // Left-associative
	case token.MUL, token.DIV:
		return 2, false // Left-associative
	// case token.POW:
	// 	return 3, true
	default:
		return 0, false
	}
}
