//go:build js && wasm

package llamawasm

// errBadHandle is what a function of the shim returns for a bad handle if its
// normal result can also be negative. It is the same value as
// YZMA_ERR_BAD_HANDLE in wasm/yzma_wasm.cpp.
const errBadHandle = -1000000

// Tokenize turns text into tokens.
func Tokenize(vocab Vocab, text string, addSpecial bool, parseSpecial bool) []Token {
	if !Loaded() {
		return nil
	}

	textPtr, err := textScratch.reserve(len(text) + 1)
	if err != nil {
		return nil
	}
	writeString(textPtr, text)

	// One token holds at least one byte of text, so the number of bytes plus
	// room for the special tokens is always enough.
	max := len(text) + 8

	for {
		tokenPtr, err := tokenScratch.reserve(max * 4)
		if err != nil {
			return nil
		}

		n := call("_yzma_tokenize", int(vocab), textPtr, len(text), tokenPtr, max,
			boolToInt(addSpecial), boolToInt(parseSpecial))
		switch {
		case n <= errBadHandle:
			return nil
		case n < 0:
			// The shim gives the negative of the number of tokens that it
			// needs.
			max = int(-n)
			continue
		default:
			return readTokens(tokenPtr, int(n))
		}
	}
}

// Detokenize turns tokens back into text.
//
// The shim has no call for this, so the text comes from TokenToPiece of each
// token, which is what a generation loop does.
func Detokenize(vocab Vocab, tokens []Token, removeSpecial bool, unparseSpecial bool) string {
	if !Loaded() {
		return ""
	}

	buf := make([]byte, 64)
	out := make([]byte, 0, len(tokens)*4)
	for _, token := range tokens {
		n := TokenToPiece(vocab, token, buf, 0, unparseSpecial)
		if n <= 0 {
			continue
		}
		out = append(out, buf[:n]...)
	}
	return string(out)
}

// TokenToPiece writes the text of one token into buf and gives the number of
// bytes that it wrote. A negative result is the negative of the number of
// bytes that buf needs.
func TokenToPiece(vocab Vocab, token Token, buf []byte, lstrip int32, special bool) int32 {
	if !Loaded() {
		return 0
	}

	ptr, err := pieceScratch.reserve(len(buf))
	if err != nil {
		return 0
	}

	n := call("_yzma_token_to_piece", int(vocab), int(token), ptr, len(buf),
		int(lstrip), boolToInt(special))
	switch {
	case n <= errBadHandle:
		return 0
	case n <= 0:
		return n
	}

	copy(buf, readBytes(ptr, int(n)))
	return n
}

// VocabIsEOG tells if a token ends the generation.
func VocabIsEOG(vocab Vocab, token Token) bool {
	if !Loaded() {
		return false
	}
	return call("_yzma_vocab_is_eog", int(vocab), int(token)) == 1
}

// VocabBOS gives the token that starts a sequence.
func VocabBOS(vocab Vocab) Token {
	return vocabToken("_yzma_vocab_bos", vocab)
}

// VocabEOS gives the token that ends a sequence.
func VocabEOS(vocab Vocab) Token {
	return vocabToken("_yzma_vocab_eos", vocab)
}

// VocabNTokens gives the number of tokens in the vocabulary.
func VocabNTokens(vocab Vocab) int32 {
	if !Loaded() {
		return 0
	}
	n := call("_yzma_vocab_n_tokens", int(vocab))
	if n <= errBadHandle {
		return 0
	}
	return n
}

// VocabGetAddBOS tells if the vocabulary adds the token that starts a sequence.
func VocabGetAddBOS(vocab Vocab) bool {
	if !Loaded() {
		return false
	}
	return call("_yzma_vocab_get_add_bos", int(vocab)) == 1
}

func vocabToken(name string, vocab Vocab) Token {
	if !Loaded() {
		return -1
	}
	t := call(name, int(vocab))
	if t <= errBadHandle {
		return -1
	}
	return Token(t)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
