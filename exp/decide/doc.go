// Package decide runs Jev-style System One models, which answer a typed
// question about a state with calibrated probabilities instead of text.
//
// The input is rendered in the macjev-render-v1 layout. Each segment is
// tokenized on its own, with no BOS or EOS and with special tokens kept as
// plain text.
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
// [Decider.DecideMany] asks several questions about one state. It decodes
// the whole ubatches of the state once, so long states gain the most.
//
// This package is experimental and its API can change.
package decide
