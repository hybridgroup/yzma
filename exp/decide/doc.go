// Package decide runs System One models, which answer a typed question about
// a state with calibrated probabilities instead of text. Three model families
// are loaded with a config file, and [Open] loads models that need none.
//
// Jev-Style models, loaded with [New], use the macjev-render-v1 layout. Each
// segment is tokenized on its own, with no BOS or EOS and with special tokens
// kept as plain text.
//
//	State:\n<state>\n\n
//	Question [<type>]: <question>\nOptions:\n
//	- <option 1>\n ... - <option K>\n
//	Judge each option:\n
//	<option 1> ->\n ... <option K> ->\n
//
// One decode reads the logits at each " ->" slot. The score of option k is
// logit(" yes") minus logit(" no") at slot k, and the probabilities are
// softmax(scores / T). The token ids and the temperatures T come from the
// readout_config.json that ships with the model, see [LoadConfig].
//
// JevK5 models, loaded with [NewJevK5], take a chat prompt with the state,
// question and lettered options as JSON. One decode reads the logits of the
// letters A to P at the last token, and the probabilities are softmax(logits / T).
// A question with more than 16 options is read in several passes and combined.
// T comes from jevk5_config.json, see [LoadJevK5Config].
//
// Decider models, loaded with [NewDeciderModel], take the state and one
// question with lettered options.
//
//	Context:\n<state>
//	\n\nQuestion: <question>\nOptions:\n(A) <option 1> ... \n(J) <option 10>
//	\nAnswer: (
//
// One decode reads the logits of the option labels at the last token, and the
// probabilities are softmax(logits / T). A question with more than 10 options
// uses A to Z and then two letter labels, up to 255 options. With
// isolated_levels, each level of a score question is its own yes or no
// question, and the level probabilities are the yes probabilities normalized.
// T comes from decider_config.json, see [LoadDeciderConfig].
//
// [Open] loads a model whose GGUF holds its decision type, prompt template
// and temperatures, as converted for the llama-server /v1/systemone API. The
// prompt is the GGUF's "systemone" template. OpenJev reads the logits of the
// letters A to Z and a to z at the last token. Lev reads the labels A to Z and
// then AA to ZZ, reads a choice question in both orders and averages them, and
// reads a noul question as a rating from 0 to 8.
//
// Laya and Julia-1 are encoders with a decision head in the llama.cpp graph.
// The score of an option is the output at its mask token, in the column of
// the question type, and each question is decoded on its own. Kev scores an
// option by the dot product of the outputs of the last token and of the token
// that ends the option. Both need a llama.cpp build that has these models.
//
// [Decider.DecideMany] asks several questions about one state and decodes
// the state once. See [ManyMode] for the exact and batched modes.
//
// [ParseRequest] and [Decider.Answer] use the JSON of the TypeSafe
// /v1/systemone API, the same one llama-server serves.
//
// This package is experimental and its API can change.
package decide
