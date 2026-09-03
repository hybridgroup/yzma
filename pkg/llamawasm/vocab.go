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

	// One token holds a minimum of one byte of text. Thus the number of bytes
	// with space for the special tokens is always sufficient.
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
// The shim has no call for this. Thus the text comes from TokenToPiece of each
// token, which is the method that a generation loop uses.
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

// VocabEOT gives the token that ends a turn.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabEOT(vocab Vocab) Token {
	if !has("_yzma_vocab_eot") {
		return -1
	}
	return vocabToken("_yzma_vocab_eot", vocab)
}

// VocabSEP gives the token that separates two sentences.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabSEP(vocab Vocab) Token {
	if !has("_yzma_vocab_sep") {
		return -1
	}
	return vocabToken("_yzma_vocab_sep", vocab)
}

// VocabNL gives the token of a new line.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabNL(vocab Vocab) Token {
	if !has("_yzma_vocab_nl") {
		return -1
	}
	return vocabToken("_yzma_vocab_nl", vocab)
}

// VocabPAD gives the token that fills a batch.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabPAD(vocab Vocab) Token {
	if !has("_yzma_vocab_pad") {
		return -1
	}
	return vocabToken("_yzma_vocab_pad", vocab)
}

// VocabMASK gives the token that hides a position.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabMASK(vocab Vocab) Token {
	if !has("_yzma_vocab_mask") {
		return -1
	}
	return vocabToken("_yzma_vocab_mask", vocab)
}

// VocabFIMPre gives the token before the text of a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMPre(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_pre") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_pre", vocab)
}

// VocabFIMSuf gives the token after the text of a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMSuf(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_suf") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_suf", vocab)
}

// VocabFIMMid gives the token of the middle of a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMMid(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_mid") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_mid", vocab)
}

// VocabFIMPad gives the token that fills a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMPad(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_pad") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_pad", vocab)
}

// VocabFIMRep gives the token of the repository of a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMRep(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_rep") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_rep", vocab)
}

// VocabFIMSep gives the token that separates the files of a fill in the middle prompt.
//
// The result is -1 if the vocabulary has no such token, or if the module is
// from a release before ABI version 5.
func VocabFIMSep(vocab Vocab) Token {
	if !has("_yzma_vocab_fim_sep") {
		return -1
	}
	return vocabToken("_yzma_vocab_fim_sep", vocab)
}

// VocabGetAddEOS tells if the vocabulary adds the token that ends a sequence.
func VocabGetAddEOS(vocab Vocab) bool {
	return vocabFlag("_yzma_vocab_get_add_eos", vocab)
}

// VocabGetAddSEP tells if the vocabulary adds the token that separates two
// sentences.
func VocabGetAddSEP(vocab Vocab) bool {
	return vocabFlag("_yzma_vocab_get_add_sep", vocab)
}

// VocabIsControl tells if a token controls the model instead of holding text.
func VocabIsControl(vocab Vocab, token Token) bool {
	if !has("_yzma_vocab_is_control") {
		return false
	}
	return call("_yzma_vocab_is_control", int(vocab), int(token)) == 1
}

// VocabGetAttr gives the attributes of a token.
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

// GetVocabType gives the kind of the tokenizer of the vocabulary.
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

// VocabGetScore gives the score of a token, or 0 if the vocabulary has none.
func VocabGetScore(vocab Vocab, token Token) float32 {
	if !has("_yzma_vocab_get_score") {
		return 0
	}
	return float32(callValue("_yzma_vocab_get_score", int(vocab), int(token)).Float())
}

// VocabGetText gives the text of a token as the vocabulary holds it, which
// keeps the marks of the tokenizer. Use TokenToPiece for the text that a
// program prints.
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
			// The shim does not say how much it needs, thus ask for more.
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

// VocabGetSuppressTokens gives the tokens that the model stops the sampler from
// taking, or nil if it names none.
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
			// The shim gives the negative of the number of tokens that it
			// needs.
			max = int(-n)
		case n == 0:
			return nil
		default:
			return readTokens(ptr, int(n))
		}
	}
}

// vocabFlag reads a call of the shim that gives 1 or 0.
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
