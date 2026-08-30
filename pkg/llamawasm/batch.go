//go:build js && wasm

package llamawasm

// BatchGetOne makes a batch that holds the given tokens. The positions of the
// tokens continue from the state of the context that decodes the batch.
func BatchGetOne(tokens []Token) Batch {
	return Batch{
		NTokens: int32(len(tokens)),
		tokens:  tokens,
	}
}

// Tokens gives the tokens of the batch.
func (b Batch) Tokens() []Token {
	return b.tokens
}
