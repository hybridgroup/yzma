package llama

import (
	"errors"
	"runtime"
	"strconv"
)

// Threads gives a good number of threads for inference on the CPU of this
// machine. It counts the cores that do the arithmetic well, which is one
// thread for each physical core, and only the performance cores of a machine
// that has two kinds. The count comes from the operating system, and a machine
// that tells nothing gets half of its logical CPUs.
//
// llama.cpp asks for four threads unless a caller changes it, which is slow on
// a machine with many cores. Thus [ContextDefaultParams] sends this value for
// a batch. For one token, [InitFromModel] uses [ModelThreads].
func Threads() int32 {
	if n := mathCores(); n > 0 {
		return int32(n)
	}
	return int32(defaultThreads())
}

// ErrNoPerformanceCPUs says that the system does not tell which CPUs belong
// to the performance cores.
var ErrNoPerformanceCPUs = errors.New("the system does not name the performance CPUs")

// PerformanceCPUs gives one CPU for each core that does the arithmetic well,
// which is one CPU of each performance core. It gives nothing when the system
// says nothing. Use it to hold the threads of a pool to a core, see
// [NewPerformanceThreadpool].
func PerformanceCPUs() []int32 {
	return mathCPUs()
}

// defaultThreads is the count that llama.cpp uses when it can read nothing
// about the cores. Half of the logical CPUs is one thread for each physical
// core of a machine with SMT.
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

// bytesPerThread is the weight size that one thread reads for each token. More
// threads than this make the generation slower, because each has too little work.
const bytesPerThread = 80 << 20

// minModelThreads is the smallest count that [ModelThreads] gives.
const minModelThreads = 4

// ModelThreads gives a good number of threads to generate one token with the
// model. The generation waits on memory, thus the count comes from the bytes
// that one token reads. A MoE model reads only the experts that it uses. The
// count is at least 4 and no more than [Threads].
func ModelThreads(model Model) int32 {
	return int32(threadsForSize(activeBytes(model), int(Threads())))
}

// threadsForSize gives the thread count for a model that reads size bytes for
// each token, on a machine with limit good cores.
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

// activeBytes gives the bytes of weights that one token reads.
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
	// A Mixtral model gives the expert size as the size of the dense layer.
	if moe.ffn == 0 {
		moe.ffn = key("feed_forward_length")
	}
	return moe.activeBytes(size, ModelNParams(model))
}

// activeBytes gives the part of size that one token reads. Each expert has
// three matrices of embd by ffn in each layer.
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
