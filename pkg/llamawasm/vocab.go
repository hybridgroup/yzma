//go:build js && wasm

package llamawasm

// errBadHandle is what a shim function returns for a bad handle when its
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

	// Each token holds at least one byte of text. So the byte count plus room
	// for the special tokens is always enough.
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
			// The shim returns the negated number of tokens it needs.
			max = int(-n)
			continue
		default:
			return readTokens(tokenPtr, int(n))
		}
	}
}

// Detokenize turns tokens back into text.
//
// The shim has no call for this. So the text comes from TokenToPiece on each
// token, which is what a generation loop uses.
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

// TokenToPiece writes the text of one token into buf and returns the number of
// bytes written. A negative result is the negated number of bytes that buf
// needs.
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

// VocabIsEOG reports whether a token ends generation.
func VocabIsEOG(vocab Vocab, token Token) bool {
	if !Loaded() {
		return false
	}
	return call("_yzma_vocab_is_eog", int(vocab), int(token)) == 1
}

// VocabBOS returns the beginning of sequence token.
func VocabBOS(vocab Vocab) Token {
	return vocabToken("_yzma_vocab_bos", vocab)
}

// VocabEOS returns the end of sequence token.
func VocabEOS(vocab Vocab) Token {
	return vocabToken("_yzma_vocab_eos", vocab)
}

// VocabNTokens returns the number of tokens in the vocabulary.
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

// VocabGetAddBOS reports whether the vocabulary adds a start of sequence token.
func VocabGetAddBOS(vocab Vocab) bool {
	if !Loaded() {
		return false
	}
	return call("_yzma_vocab_get_add_bos", int(vocab)) == 1
}

// VocabEOT returns the end of turn token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabEOT(vocab Vocab) Token {
	if !has("_yzma_vocab_eot") {
		return -1
	}
	return vocabToken("_yzma_vocab_eot", vocab)
}

// VocabSEP returns the sentence separator token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabSEP(vocab Vocab) Token {
	if !has("_yzma_vocab_sep") {
		return -1
	}
	return vocabToken("_yzma_vocab_sep", vocab)
}

// VocabNL returns the newline token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabNL(vocab Vocab) Token {
	if !has("_yzma_vocab_nl") {
		return -1
	}
	return vocabToken("_yzma_vocab_nl", vocab)
}

// VocabPAD returns the padding token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabPAD(vocab Vocab) Token {
	if !has("_yzma_vocab_pad") {
		return -1
	}
	return vocabToken("_yzma_vocab_pad", vocab)
}

// VocabMASK returns the mask token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabMASK(vocab Vocab) Token {
	if !has("_yzma_vocab_mask") {
		return -1
	}
	return vocabToken("_yzma_vocab_mask", vocab)
}

// VocabFIMPre returns the fill in the middle prefix token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMPre(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_pre") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_pre", vocab)
}

// VocabFIMSuf returns the fill in the middle suffix token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMSuf(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_suf") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_suf", vocab)
}

// VocabFIMMid returns the fill in the middle middle token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMMid(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_mid") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_mid", vocab)
}

// VocabFIMPad returns the fill in the middle padding token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMPad(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_pad") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_pad", vocab)
}

// VocabFIMRep returns the fill in the middle repository token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMRep(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_rep") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_rep", vocab)
}

// VocabFIMSep returns the fill in the middle file separator token.
//
// It returns -1 if the vocabulary has no such token, or if the module predates
// ABI version 5.
func VocabFIMSep(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_sep") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_sep", vocab)
}

// VocabGetAddEOS reports whether the vocabulary adds an end of sequence token.
func VocabGetAddEOS(vocab Vocab) bool {
	return vocabFlag("_yzma_vocab_get_add_eos", vocab)
}

// VocabGetAddSEP reports whether the vocabulary adds a sentence separator
// token.
func VocabGetAddSEP(vocab Vocab) bool {
	return vocabFlag("_yzma_vocab_get_add_sep", vocab)
}

// VocabIsControl reports whether a token is a control token and not text.
func VocabIsControl(vocab Vocab, token Token) bool {
	if !has("_yzma_vocab_is_control") {
		return false
	}
	return call("_yzma_vocab_is_control", int(vocab), int(token)) == 1
}

// VocabGetAttr returns the attributes of a token.
func VocabGetAttr(vocab Vocab, token Token) TokenAttr {
	if !has("_yzma_vocab_get_attr") {
		return TokenAttrUndefined
	}

	attr := call("_yzma_vocab_get_attr", int(vocab), int(token))
	if attr < 0 {
		return TokenAttrUndefined
	}
	return TokenAttr(attr)
}

// GetVocabType returns the tokenizer type of the vocabulary.
func GetVocabType(vocab Vocab) VocabType {
	if !has("_yzma_vocab_type") {
		return VocabTypeNone
	}

	t := call("_yzma_vocab_type", int(vocab))
	if t < 0 {
		return VocabTypeNone
	}
	return VocabType(t)
}

// VocabGetScore returns the score of a token, or 0 if the vocabulary has none.
func VocabGetScore(vocab Vocab, token Token) float32 {
	if !has("_yzma_vocab_get_score") {
		return 0
	}
	v, _ := callNumber("_yzma_vocab_get_score", int(vocab), int(token))
	return float32(v)
}

// VocabGetText returns the raw token text from the vocabulary, including
// tokenizer markers. Use TokenToPiece for text that a program prints.
func VocabGetText(vocab Vocab, token Token) string {
	if !has("_yzma_vocab_get_text") {
		return ""
	}

	size := 64
	for {
		ptr, err := pieceScratch.reserve(size)
		if err != nil {
			return ""
		}

		n := call("_yzma_vocab_get_text", int(vocab), int(token), ptr, size)
		switch {
		case n == errTooSmall:
			// The shim does not report how much it needs, so ask for more.
			size *= 4
			if size > 1<<20 {
				return ""
			}
		case n <= 0:
			return ""
		default:
			return string(readBytes(ptr, int(n)))
		}
	}
}

// VocabGetSuppressTokens returns the tokens the model keeps the sampler from
// picking, or nil if there are none.
func VocabGetSuppressTokens(vocab Vocab) []Token {
	if !has("_yzma_vocab_get_suppress_tokens") {
		return nil
	}

	max := 64
	for {
		ptr, err := tokenScratch.reserve(max * 4)
		if err != nil {
			return nil
		}

		n := call("_yzma_vocab_get_suppress_tokens", int(vocab), ptr, max)
		switch {
		case n <= errBadHandle:
			return nil
		case n < 0:
			// The shim returns the negated number of tokens it needs.
			max = int(-n)
		case n == 0:
			return nil
		default:
			return readTokens(ptr, int(n))
		}
	}
}

// vocabFlag reads a shim call that returns 1 or 0.
func vocabFlag(name string, vocab Vocab) bool {
	if !has(name) {
		return false
	}
	return call(name, int(vocab)) == 1
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
