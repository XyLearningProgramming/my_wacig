package my_parser

import (
	"fmt"
	"monkey/my_ast"
	token "monkey/my_token"
	"strconv"
)

type PrecedenceLevel int

const (
	_ PrecedenceLevel = iota
	LOWEST
	REASSIGN    // let a =1; a=2;
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PREFIX      // -X or !X
	PRODUCT     // *
	CALL        // myFunction(X)
	INDEX       // []
	INDEXCOLON  // :
)

var InfixOperatorToPrecedences = map[my_ast.InfixOperator]PrecedenceLevel{
	my_ast.INOP_MINUS:      SUM,
	my_ast.INOP_PLUS:       SUM,
	my_ast.INOP_ASTERISK:   PRODUCT,
	my_ast.INOP_SLASH:      PRODUCT,
	my_ast.INOP_LT:         LESSGREATER,
	my_ast.INOP_GT:         LESSGREATER,
	my_ast.INOP_EQ:         EQUALS,
	my_ast.INOP_NOT_EQ:     EQUALS,
	my_ast.INOP_CALL:       CALL,
	my_ast.INOP_INDEX:      INDEX,
	my_ast.INOP_INDEXCOLON: INDEXCOLON,
	my_ast.INOP_GTE:        LESSGREATER,
	my_ast.INOP_LTE:        LESSGREATER,
	my_ast.INOP_REASSIGN:   REASSIGN,
}

func tokenPrecedenceLevel(t *token.Token) PrecedenceLevel {
	if pl, pok :=
		InfixOperatorToPrecedences[my_ast.InfixOperator(t.Type)]; pok {
		return pl
	}
	return LOWEST
}

type (
	prefixParseFn func() my_ast.Expression
	infixParseFn  func(my_ast.Expression) my_ast.Expression
)

func (p *Parser) parseExpression(precedence PrecedenceLevel) my_ast.Expression {
	prefixExpr, pok := p.prefixParseFns[p.curToken.Type]
	if !pok {
		p.appendExprFuncError(p.curToken, true)
		return nil
	}
	leftExpr := prefixExpr()

	// NOTE: consume to semicolon or EOF
	// or when meet a higher precedence with current token
	for p.peekToken.Type != token.SEMICOLON &&
		p.peekToken.Type != token.EOF &&
		precedence < tokenPrecedenceLevel(&p.peekToken) {
		infixFn, iok := p.infixParseFns[p.peekToken.Type]
		if !iok {
			return leftExpr
		}
		p.nextToken()
		leftExpr = infixFn(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseIdentifier() my_ast.Expression {
	return &my_ast.Identifier{
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() my_ast.Expression {
	val, err := strconv.ParseUint(p.curToken.Literal, 10, 64)
	if err != nil {
		p.appendError(fmt.Sprintf("cannot parse %s as uint :%v", p.curToken.Literal, err))
		return nil
	}
	return &my_ast.Integer{Value: val}
}

func (p *Parser) parseFloatLiteral() my_ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.appendError(fmt.Sprintf("cannot parse %s as float: %v", p.curToken.Literal, err))
		return nil
	}
	return &my_ast.Float{Value: val}
}

func (p *Parser) parseBooleanLiteral() my_ast.Expression {
	if p.curToken.Type == token.TRUE {
		return &my_ast.Boolean{Value: true}
	} else if p.curToken.Type == token.FALSE {
		return &my_ast.Boolean{Value: false}
	} else {
		p.appendError(fmt.Sprintf(
			"cannot parse %s with type %s as boolean", p.curToken.Literal, p.curToken.Type,
		))
		return nil
	}
}

func (p *Parser) parsePrefixExpression() my_ast.Expression {
	expr := &my_ast.PrefixExpression{
		Operator: my_ast.PrefixOperator(p.curToken.Type),
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left my_ast.Expression) my_ast.Expression {
	exp := &my_ast.InfixExpression{
		Left:     left,
		Operator: my_ast.InfixOperator(p.curToken.Type),
	}
	precedence := tokenPrecedenceLevel(&p.curToken)
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseGroupedExpression() my_ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.isPeekToken(token.RPAREN) {
		p.appendTokenError(token.RPAREN, p.peekToken)
		p.nextToken()
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) parseIfExpression() my_ast.Expression {
	// parse if condition as expression
	p.nextToken()
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	// xx )
	ie := &my_ast.IfExpression{Condition: p.parseExpression(LOWEST)}
	p.nextToken()
	if !p.isCurToken(token.RPAREN) {
		p.appendTokenError(token.RPAREN, p.curToken)
		return nil
	}
	// parse if consequence as block statement
	p.nextToken()

	ie.Consequence = p.parseBlockStatement()
	if !p.isCurToken(token.RBRACE) {
		p.appendTokenError(token.RBRACE, p.curToken)
		return nil
	}
	// } else
	// no "else" token is legal, return immediately
	if !p.isPeekToken(token.ELSE) {
		return ie
	}
	// else { xx
	p.nextToken()
	// parse if alternative as block statement
	p.nextToken()
	// { xx
	ie.Alternative = p.parseBlockStatement()
	if !p.isCurToken(token.RBRACE) {
		p.appendTokenError(token.RBRACE, p.curToken)
		return nil
	}
	return ie
}

func (p *Parser) parseFunction() my_ast.Expression {
	p.nextToken()
	// TODO: no function name after fn?
	fe := &my_ast.Function{
		Parameters: p.parseFunctionParameters(),
	}
	if !p.isCurToken(token.RPAREN) {
		p.appendTokenError(token.RPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	fe.Body = p.parseBlockStatement()
	return fe
}

func (p *Parser) parseFunctionParameters() []*my_ast.Identifier {
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	params := []*my_ast.Identifier{}
	if p.isCurToken(token.RPAREN) {
		// p.nextToken()
		return params
	}
	params = append(params, &my_ast.Identifier{Value: p.curToken.Literal})
	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &my_ast.Identifier{Value: p.curToken.Literal})
	}
	p.nextToken()
	if !p.isCurToken(token.RPAREN) {
		p.appendTokenError(token.RPAREN, p.curToken)
		return nil
	}
	return params
}

func (p *Parser) parseCallExpression(leftFunc my_ast.Expression) my_ast.Expression {
	return &my_ast.CallExpression{Function: leftFunc, Arguments: p.parseCallArguments()}
}

func (p *Parser) parseCallArguments() []my_ast.Expression {
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	args := []my_ast.Expression{}
	if p.isCurToken(token.RPAREN) {
		return args
	}
	args = append(args, p.parseExpression(LOWEST))
	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.isPeekToken(token.RPAREN) {
		p.appendTokenError(token.RPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	return args
}

func (p *Parser) parseStringExpression() my_ast.Expression {
	return &my_ast.StringExpression{
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayExpression() my_ast.Expression {
	return &my_ast.ArrayExpression{Elements: p.parseExpressionList(token.RBRACKET)}
}

func (p *Parser) parseExpressionList(end token.TokenType) []my_ast.Expression {
	list := []my_ast.Expression{}
	if p.isPeekToken(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.isPeekToken(end) {
		p.appendTokenError(end, p.peekToken)
		return nil
	}
	p.nextToken()
	return list
}

func (p *Parser) parseIndexExpression(left my_ast.Expression) my_ast.Expression {
	exp := &my_ast.IndexExpression{
		Left:            left,
		StartIndex:      nil,
		IsSetStartIndex: false,
		EndIndex:        nil,
		IsSetEndIndex:   false,
		Stride:          nil,
		IsSetStride:     false,
	}
	p.nextToken()
	if p.isCurToken(token.RBRACKET) {
		return exp
	}
	// parse start index to : or ]
	if !p.isCurToken(token.COLON) {
		startIdx := p.parseExpression(LOWEST)
		if startIdx == nil {
			return nil
		}
		exp.StartIndex = startIdx
		exp.IsSetStartIndex = true
		p.nextToken()
	}
	if p.isCurToken(token.RBRACKET) {
		return exp
	}
	if !p.isCurToken(token.COLON) {
		p.appendError(fmt.Sprintf("Expected : or ], but got: %s", p.curToken.Literal))
		return nil
	}
	exp.IsSetStartIndex = true
	exp.IsSetEndIndex = true
	p.nextToken()
	if p.isCurToken(token.RBRACKET) {
		return exp
	}
	// parse end index to : or ]
	if !p.isCurToken(token.COLON) {
		endIdx := p.parseExpression(LOWEST)
		if endIdx == nil {
			return nil
		}
		exp.EndIndex = endIdx
		exp.IsSetEndIndex = true
		p.nextToken()
	}
	if p.isCurToken(token.RBRACKET) {
		return exp
	}
	if !p.isCurToken(token.COLON) {
		p.appendError(fmt.Sprintf("Expected : or ], but got: %s", p.curToken.Literal))
		return nil
	}

	exp.IsSetStride = true
	p.nextToken()
	if p.isCurToken(token.RBRACKET) {
		return exp
	}
	// parse stride to ]
	stride := p.parseExpression(LOWEST)
	if stride == nil {
		return nil
	}
	exp.Stride = stride
	p.nextToken()
	if !p.isCurToken(token.RBRACKET) {
		p.appendTokenError(token.RBRACKET, p.curToken)
		return nil
	}
	return exp
}

func (p *Parser) parseHashLiteral() my_ast.Expression {
	hash := &my_ast.HashExpression{
		Pairs: make(map[my_ast.Expression]my_ast.Expression),
		Keys:  make([]my_ast.Expression, 0),
	}
	p.nextToken()
	for !p.isCurToken(token.RBRACE) && !p.isCurToken(token.EOF) {
		key := p.parseExpression(LOWEST)
		if !p.isPeekToken(token.COLON) {
			p.appendTokenError(token.COLON, p.peekToken)
			return nil
		}
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		hash.Keys = append(hash.Keys, key)
		if !p.isPeekToken(token.RBRACE) && !p.isPeekToken(token.COMMA) {
			p.appendError(
				fmt.Sprintf("expecting token RBRACE or COMMA, but got %s with literal %s instead", string(p.peekToken.Type), p.peekToken.Literal),
			)
			return nil
		}
		if p.isPeekToken(token.COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}
	if !p.isCurToken(token.RBRACE) {
		p.appendTokenError(token.RBRACE, p.peekToken)
		return nil
	}
	return hash
}

func (p *Parser) parseForExpression() my_ast.Expression {
	p.nextToken()
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	forExpression := &my_ast.ForExpression{}
	if !p.isCurToken(token.SEMICOLON) {
		forExpression.InitStmt = p.parseStatement()
		if forExpression.InitStmt == nil {
			return nil
		}
	} else {
		forExpression.InitStmt = nil
	}
	p.nextToken()
	if !p.isCurToken(token.SEMICOLON) {
		forExpression.TestExpr = p.parseExpression(LOWEST)
		if forExpression.TestExpr == nil {
			return nil
		}
		p.nextToken()
	} else {
		forExpression.TestExpr = nil
	}
	p.nextToken()
	if !p.isCurToken(token.RPAREN) {
		forExpression.UpdateStmt = p.parseStatement()
		if forExpression.UpdateStmt == nil {
			return nil
		}
		p.nextToken()
	} else {
		forExpression.UpdateStmt = nil
	}
	p.nextToken()
	forExpression.Body = p.parseBlockStatement()
	return forExpression
}

func (p *Parser) parseWhileExression() my_ast.Expression {
	p.nextToken()
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	whileExpr := &my_ast.WhileExpression{}
	if !p.isCurToken(token.RPAREN) {
		whileExpr.TestExpr = p.parseExpression(LOWEST)
		if whileExpr.TestExpr == nil {
			return nil
		}
		p.nextToken()
	} else {
		whileExpr.TestExpr = nil
	}
	p.nextToken()
	if !p.isCurToken(token.LBRACE) {
		p.appendTokenError(token.LBRACE, p.curToken)
		return nil
	}
	whileExpr.Body = p.parseBlockStatement()
	if whileExpr.Body == nil {
		return nil
	}
	return whileExpr
}

func (p *Parser) parseDoWhileExpression() my_ast.Expression {
	p.nextToken()
	if !p.isCurToken(token.LBRACE) {
		p.appendTokenError(token.LBRACE, p.curToken)
		return nil
	}
	doWhileExpr := &my_ast.DoWhileExpression{
		Body: p.parseBlockStatement(),
	}
	if doWhileExpr.Body == nil {
		return nil
	}
	if !p.isCurToken(token.RBRACE) {
		p.appendTokenError(token.RBRACE, p.curToken)
		return nil
	}
	p.nextToken()
	if !p.isCurToken(token.WHILE) {
		p.appendTokenError(token.WHILE, p.curToken)
		return nil
	}
	p.nextToken()
	if !p.isCurToken(token.LPAREN) {
		p.appendTokenError(token.LPAREN, p.curToken)
		return nil
	}
	p.nextToken()
	if !p.isCurToken(token.RPAREN) {
		doWhileExpr.TestExpr = p.parseExpression(LOWEST)
		if doWhileExpr.TestExpr == nil {
			return nil
		}
		p.nextToken()
	} else {
		doWhileExpr.TestExpr = nil
	}
	return doWhileExpr
}

func (p *Parser) parseNullLiteral() my_ast.Expression {
	return my_ast.NULL
}
