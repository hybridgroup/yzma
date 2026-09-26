// Package decide runs System One models, which answer a typed question about
// a state with calibrated probabilities instead of text. Two model families
// are supported.
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
// [Decider.DecideMany] asks several questions about one state and decodes
// the state once. See [ManyMode] for the exact and batched modes.
//
// This package is experimental and its API can change.
package decide
