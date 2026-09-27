package llama

import (
	"errors"
	"runtime"
	"strconv"
)

// Threads returns a good thread count for CPU inference on this machine. It
// counts one thread per physical core, and only the performance cores on a
// hybrid machine. The count comes from the operating system. If the system
// reports nothing, it uses half of the logical CPUs.
//
// llama.cpp defaults to four threads, which is slow on a machine with many
// cores. So [ContextDefaultParams] uses this value for batch processing. For one token, [InitFromModel] uses [ModelThreads].
func Threads() int32 {
	if n := mathCores(); n > 0 {
		return int32(n)
	}
	return int32(defaultThreads())
}

// ErrNoPerformanceCPUs means the system does not report which CPUs belong to
// the performance cores.
var ErrNoPerformanceCPUs = errors.New("the system does not name the performance CPUs")

// PerformanceCPUs returns one CPU for each performance core, or nil when the
// system does not report them. Use it to pin pool threads to cores, see
// [NewPerformanceThreadpool].
func PerformanceCPUs() []int32 {
	return mathCPUs()
}

// defaultThreads is the count llama.cpp uses when it cannot read the core
// topology. Half the logical CPUs is one thread per physical core with SMT.
func defaultThreads() int {
	n := runtime.NumCPU()
	if n > 4 {
		n /= 2
	}
	if n < 1 {
		n = 1
	}
	return n
}

// bytesPerThread is the weight size one thread reads per token. More threads
// than this slow generation down, because each one has too little work.
const bytesPerThread = 80 << 20

// minModelThreads is the smallest count that [ModelThreads] returns.
const minModelThreads = 4

// ModelThreads returns a good thread count for generating one token with the
// model. Generation is memory bound, so the count comes from the bytes each
// token reads. A MoE model reads only the experts it uses. The count is at
// least 4 and at most [Threads].
func ModelThreads(model Model) int32 {
	return int32(threadsForSize(activeBytes(model), int(Threads())))
}

// threadsForSize returns the thread count for a model that reads size bytes
// per token, on a machine with limit usable cores.
func threadsForSize(size uint64, limit int) int {
	n := min(int(size/bytesPerThread), limit)
	n = max(n, min(minModelThreads, limit))
	return max(n, 1)
}

// moeParams holds the expert metadata of a MoE model.
type moeParams struct {
	experts, used, shared uint64
	ffn, embd, layers     uint64
}

// activeBytes returns the weight bytes that one token reads.
func activeBytes(model Model) uint64 {
	size := ModelSize(model)
	arch, ok := ModelMetaValStr(model, "general.architecture")
	if !ok {
		return size
	}
	key := func(name string) uint64 {
		v, ok := ModelMetaValStr(model, arch+"."+name)
		if !ok {
			return 0
		}
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		return n
	}
	moe := moeParams{
		experts: key("expert_count"),
		used:    key("expert_used_count"),
		shared:  key("expert_shared_count"),
		ffn:     key("expert_feed_forward_length"),
		embd:    key("embedding_length"),
		layers:  key("block_count"),
	}
	// Mixtral stores the expert size as the dense layer size.
	if moe.ffn == 0 {
		moe.ffn = key("feed_forward_length")
	}
	return moe.activeBytes(size, ModelNParams(model))
}

// activeBytes returns the part of size that one token reads. Each expert has
// three embd by ffn matrices per layer.
func (m moeParams) activeBytes(size, params uint64) uint64 {
	if m.experts == 0 || m.used == 0 || m.used >= m.experts {
		return size
	}
	perExpert := m.layers * 3 * m.embd * m.ffn
	all := m.experts * perExpert
	if perExpert == 0 || params == 0 || all >= params {
		return uint64(float64(size) * float64(m.used) / float64(m.experts))
	}
	active := params - all + min(m.used+m.shared, m.experts)*perExpert
	return uint64(float64(size) * float64(active) / float64(params))
}
