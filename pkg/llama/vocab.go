package llama

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

var (
	// LLAMA_API const struct llama_vocab * llama_model_get_vocab(const struct llama_model * model);
	modelGetVocabFunc ffi.Fun

	// LLAMA_API llama_token llama_vocab_bos(const struct llama_vocab * vocab); // beginning-of-sentence
	vocabBOSFunc ffi.Fun

	// LLAMA_API llama_token llama_vocab_eos(const struct llama_vocab * vocab); // end-of-sentence
	vocabEOSFunc ffi.Fun

	// LLAMA_API bool llama_vocab_is_eog(const struct llama_vocab * vocab, llama_token token);
	vocabIsEOGFunc ffi.Fun

	// LLAMA_API bool llama_vocab_is_control(const struct llama_vocab * vocab, llama_token token);
	vocabIsControlFunc ffi.Fun

	// LLAMA_API int32_t llama_vocab_n_tokens(const struct llama_vocab * vocab);
	vocabNTokensFunc ffi.Fun

	// LLAMA_API int32_t llama_token_to_piece(
	// 		const struct llama_vocab * vocab,
	// 					llama_token   token,
	// 							char * buf,
	// 						int32_t   length,
	// 						int32_t   lstrip,
	// 							bool   special);
	tokenToPieceFunc ffi.Fun

	// LLAMA_API int32_t llama_tokenize(
	//     const struct llama_vocab * vocab,
	//                   const char * text,
	//                      int32_t   text_len,
	//                  llama_token * tokens,
	//                      int32_t   n_tokens_max,
	//                         bool   add_special,
	//                         bool   parse_special);
	tokenizeFunc ffi.Fun
)

func loadVocabFuncs(lib ffi.Lib) error {
	var err error
	if modelGetVocabFunc, err = lib.Prep("llama_model_get_vocab", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if vocabBOSFunc, err = lib.Prep("llama_vocab_bos", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if vocabEOSFunc, err = lib.Prep("llama_vocab_eos", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if vocabIsEOGFunc, err = lib.Prep("llama_vocab_is_eog", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if vocabIsControlFunc, err = lib.Prep("llama_vocab_is_control", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if vocabNTokensFunc, err = lib.Prep("llama_vocab_n_tokens", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if tokenToPieceFunc, err = lib.Prep("llama_token_to_piece", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeUint8); err != nil {
		return err
	}

	if tokenizeFunc, err = lib.Prep("llama_tokenize", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8, &ffi.TypeUint8); err != nil {
		return err
	}

	return nil

}

func ModelGetVocab(model Model) Vocab {
	var vocab Vocab
	modelGetVocabFunc.Call(unsafe.Pointer(&vocab), unsafe.Pointer(&model))

	return vocab
}
func VocabBOS(vocab Vocab) Token {
	var token ffi.Arg
	vocabBOSFunc.Call(unsafe.Pointer(&token), unsafe.Pointer(&vocab))

	return Token(token)
}
func VocabEOS(vocab Vocab) Token {
	var token ffi.Arg
	vocabEOSFunc.Call(unsafe.Pointer(&token), unsafe.Pointer(&vocab))

	return Token(token)
}

func VocabIsEOG(vocab Vocab, token Token) bool {
	var result ffi.Arg
	vocabIsEOGFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&vocab), unsafe.Pointer(&token))

	return result.Bool()
}

func VocabIsControl(vocab Vocab, token Token) bool {
	var result ffi.Arg
	vocabIsControlFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&vocab), unsafe.Pointer(&token))

	return result.Bool()
}

func VocabNTokens(vocab Vocab) int32 {
	var result ffi.Arg
	vocabNTokensFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&vocab))

	return int32(result)
}

func TokenToPiece(vocab Vocab, token Token, buf []byte, lstrip int32, special bool) int32 {
	piece := make([]byte, len(buf))
	b := unsafe.SliceData(piece)
	bLen := int32(len(piece))

	var result ffi.Arg
	tokenToPieceFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&vocab), &token, unsafe.Pointer(&b),
		&bLen, &lstrip, &special)

	resLen := int32(result)
	if resLen > 0 && resLen <= int32(len(piece)) {
		copy(buf, piece[:resLen])
	}

	return resLen
}

func Tokenize(vocab Vocab, text string, tokens []Token, addSpecial bool, parseSpecial bool) int32 {
	txt, _ := utils.BytePtrFromString(text)
	txtLen := int32(len(text))

	toks := unsafe.SliceData(tokens)
	nTokensMax := len(tokens)

	var result ffi.Arg
	tokenizeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&vocab), unsafe.Pointer(&txt), &txtLen,
		unsafe.Pointer(&toks), &nTokensMax, &addSpecial, &parseSpecial)

	return int32(result)
}
