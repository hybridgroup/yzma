package download

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	getter "github.com/hashicorp/go-getter"
)

// Target identifies the machine a llama.cpp build is for.
type Target struct {
	Arch      Arch
	OS        OS
	Processor Processor

	// Version is the llama.cpp release tag, e.g. "b7974" or "v0.3.0". "" takes
	// [DefaultVersion], or the newest release if that is empty. "latest" always
	// resolves to the newest release.
	Version string

	// UpstreamVersion is the nightly build tag that has the llama.cpp binaries for
	// Version. A tagged release has no binaries of its own, so [Install] fills this
	// in with [LlamaNightlyTag]. Empty means use Version.
	UpstreamVersion string

	// ManifestSHA256 is the expected digest of the raw digest manifest of Version,
	// in hexadecimal. Empty means the digest is not pinned, which is the usual case.
	//
	// A pin makes verification mandatory. The manifest bytes must agree, the
	// manifest must name every resolved asset, and the bytes of each asset must
	// agree with it. [Install] also takes the pin as a suffix on Version, in the
	// form "b10785@sha256:<digest>", and moves it here.
	ManifestSHA256 string
}

// Resolver reports the release assets to install for a Target, as URLs downloaded in
// the order returned. Implement it to reach builds the built-in table does not name.
// Implementations must not download.
type Resolver interface {
	Resolve(target Target) (urls []string, err error)
}

// ResolverFunc adapts an ordinary function to [Resolver].
type ResolverFunc func(target Target) ([]string, error)

// Resolve calls f.
func (f ResolverFunc) Resolve(target Target) ([]string, error) { return f(target) }

// DefaultResolver resolves the assets published on the llama.cpp and llama-cpp-builder
// release pages. [Install] uses it when no resolver is given. It satisfies
// [AssetResolver] as well, so it reports the expected digest of each asset.
var DefaultResolver Resolver = defaultResolver{}

// defaultResolver is the built-in resolver. It reads the digest manifest that
// llama-cpp-builder publishes for each release tag.
type defaultResolver struct{}

// Resolve reports the assets to install as URLs.
func (defaultResolver) Resolve(target Target) ([]string, error) {
	return defaultResolve(target)
}

// ResolveAssets reports the assets to install with their expected digests. A manifest
// that cannot be read gives assets with no digest, which [VerifyIfAvailable] permits
// and [VerifyRequired] refuses.
//
// A [Target.ManifestSHA256] that is set makes the manifest mandatory. The bytes must
// have that digest, and a manifest that cannot be read is an error rather than an
// install with no check.
func (r defaultResolver) ResolveAssets(target Target) ([]Asset, error) {
	return r.resolveAssets(context.Background(), target)
}

// resolveAssets does the work of [defaultResolver.ResolveAssets] under a context.
// [AssetResolver] takes no context, so [Install] calls this to pass its own.
func (r defaultResolver) resolveAssets(ctx context.Context, target Target) ([]Asset, error) {
	urls, err := defaultResolve(target)
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, len(urls))
	for i, url := range urls {
		assets[i] = Asset{URL: url}
	}

	m, err := fetchManifest(ctx, target.Version, target.ManifestSHA256)
	if err != nil {
		if target.ManifestSHA256 != "" {
			return nil, err
		}
		return assets, nil
	}

	for i := range assets {
		assets[i].SHA256 = m.digestFor(assets[i].URL)
	}

	return assets, nil
}

// llama.cpp renamed its ROCm assets at these two builds.
const (
	rocmRenameBuild = 10356
	rocm10Build     = 10767
)

// rocmVersionNames holds the ROCm asset names for one build range.
type rocmVersionNames struct {
	linux   string
	windows string
}

// rocmNames reports the ROCm asset names that tag published. A tag that is not a
// nightly build takes the newest names.
func rocmNames(tag string) rocmVersionNames {
	current := rocmVersionNames{
		linux:   "llama-%s-bin-ubuntu-rocm-10.0-x64.tar.gz",
		windows: "llama-%s-bin-win-rocm-10.0-x64.zip",
	}
	if !nightlyPattern.MatchString(tag) {
		return current
	}
	build, err := strconv.Atoi(tag[1:])
	if err != nil {
		return current
	}

	switch {
	case build < rocmRenameBuild:
		return rocmVersionNames{
			linux:   "llama-%s-bin-ubuntu-rocm-7.2-x64.tar.gz",
			windows: "llama-%s-bin-win-hip-radeon-x64.zip",
		}
	case build < rocm10Build:
		return rocmVersionNames{
			linux:   "llama-%s-bin-ubuntu-rocm-7.14-x64.tar.gz",
			windows: "llama-%s-bin-win-rocm-7.14-x64.zip",
		}
	default:
		return current
	}
}

// defaultResolve is the built-in platform table.
func defaultResolve(target Target) ([]string, error) {
	arch, os, prcssr, version := target.Arch, target.OS, target.Processor, target.Version

	// The llama.cpp releases hold the binaries under the nightly build tag, while the
	// llama-cpp-builder releases use the requested tag.
	upstream := target.UpstreamVersion
	if upstream == "" {
		upstream = version
	}
	builderLocation := fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)

	var extra []string
	var filename string
	location := fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", upstream)
	tag := upstream

	switch os {
	case Linux:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cpu-arm64.tar.gz", tag)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.tar.gz", tag)
		case CUDA:
			location, tag = builderLocation, version
			if arch == ARM64 {
				// defaults to CUDA 12 assuming that is running Jetson Orin.
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-arm64.tar.gz", tag)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-13-x64.tar.gz", tag)
			}
		case Vulkan:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-arm64.tar.gz", tag)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.tar.gz", tag)
		case ROCm:
			if arch != AMD64 {
				return nil, errors.New("precompiled binaries for Linux ARM64 ROCm are not available")
			}
			filename = fmt.Sprintf(rocmNames(tag).linux, tag)
		default:
			return nil, ErrUnknownProcessor
		}

	case Bookworm:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cpu-arm64.tar.gz", tag)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		case CUDA:
			location, tag = builderLocation, version
			if arch == ARM64 {
				// Jetson Orin.
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-arm64.tar.gz", tag)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		case Vulkan:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-arm64.tar.gz", tag)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		default:
			return nil, ErrUnknownProcessor
		}

	case Trixie:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-trixie-cpu-arm64.tar.gz", tag)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.tar.gz", tag)
		case CUDA:
			location, tag = builderLocation, version
			if arch == ARM64 {
				// not yet
				return nil, ErrUnknownProcessor
			} else {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-13-x64.tar.gz", tag)
			}
		case Vulkan:
			if arch == ARM64 {
				location, tag = builderLocation, version
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-trixie-vulkan-arm64.tar.gz", tag)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.tar.gz", tag)
		default:
			return nil, ErrUnknownProcessor
		}

	case Darwin:
		switch prcssr {
		case Metal:
			if arch != ARM64 {
				return nil, errors.New("precompiled binaries for macOS non-ARM64 CPU/Metal are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.tar.gz", tag)
		case CPU:
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-macos-arm64.tar.gz", tag)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-macos-x64.tar.gz", tag)
			}
		default:
			return nil, ErrUnknownProcessor
		}

	case Windows:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-win-cpu-arm64.zip", tag)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-win-cpu-x64.zip", tag)
			}
		case CUDA:
			if arch == ARM64 {
				return nil, errors.New("precompiled binaries for Windows ARM64 CUDA are not available")
			}
			// also requires the CUDA RT files
			cudart := "cudart-llama-bin-win-cuda-13.3-x64.zip"
			extra = append(extra, fmt.Sprintf("%s/%s", location, cudart))
			filename = fmt.Sprintf("llama-%s-bin-win-cuda-13.3-x64.zip", tag)
		case Vulkan:
			if arch == ARM64 {
				return nil, errors.New("precompiled binaries for Windows ARM64 Vulkan are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip", tag)
		case ROCm:
			if arch != AMD64 {
				return nil, errors.New("precompiled binaries for Windows ARM64 ROCm are not available")
			}
			filename = fmt.Sprintf(rocmNames(tag).windows, tag)
		default:
			return nil, ErrUnknownProcessor
		}

	case Wasm:
		// Every build of the target comes down, whichever processor the caller
		// names, because the JavaScript glue chooses at run time and needs them
		// all: WebGPU where the browser has it, more than one thread where the
		// page is isolated, and one thread everywhere else.
		//
		// CUDA, Metal, ROCm and Vulkan have no meaning in a browser.
		if prcssr != CPU && prcssr != WebGPU {
			return nil, ErrUnknownProcessor
		}
		location, tag = builderLocation, version
		extra = append(extra,
			fmt.Sprintf("%s/llama-%s-bin-wasm-simd-mt.tar.gz", location, tag),
			fmt.Sprintf("%s/llama-%s-bin-wasm-webgpu.tar.gz", location, tag),
		)
		filename = fmt.Sprintf("llama-%s-bin-wasm-simd.tar.gz", tag)

	default:
		return nil, ErrUnknownOS
	}

	return append(extra, fmt.Sprintf("%s/%s", location, filename)), nil
}

// InstallOption changes what [Install] does.
type InstallOption func(*installOptions)

// installOptions holds the settings that an [InstallOption] changes.
type installOptions struct {
	verify VerifyPolicy
}

// WithVerify sets what [Install] does about the digest of an asset. The default is
// [VerifyIfAvailable].
func WithVerify(policy VerifyPolicy) InstallOption {
	return func(o *installOptions) { o.verify = policy }
}

// Install downloads the llama.cpp binaries for target into dest. A nil resolver means
// [DefaultResolver]. An empty [Target.Version] takes [DefaultVersion].
//
// Install checks the digest of each asset that has one, and stops before it writes
// anything if the bytes do not agree. Use [WithVerify] to change that.
//
// [Target.Version] may carry the expected digest of the digest manifest of the
// release, in the form "b10785@sha256:<digest>". Install moves it to
// [Target.ManifestSHA256] and uses only the tag for URLs, for the resolver, for the
// install record, and for the version it reports. A pin makes verification mandatory,
// so it cannot be given with [VerifyOff].
func Install(ctx context.Context, target Target, dest string, progress getter.ProgressTracker, resolver Resolver, opts ...InstallOption) error {
	if resolver == nil {
		resolver = DefaultResolver
	}

	var options installOptions
	for _, opt := range opts {
		opt(&options)
	}

	// The digest comes off the version before anything validates the version or
	// builds a URL from it.
	tag, digest, err := ParsePinnedVersion(target.Version)
	if err != nil {
		return err
	}
	target.Version = tag
	switch {
	case digest == "":
		// Nothing to say. The version carried no digest.
	case target.ManifestSHA256 == "":
		target.ManifestSHA256 = digest
	case !strings.EqualFold(target.ManifestSHA256, digest):
		return fmt.Errorf("%w: the version pins %s and ManifestSHA256 is %s", ErrInvalidDigest, digest, target.ManifestSHA256)
	}

	// A pin asks for a check, so it does not agree with a policy that checks nothing,
	// and it makes an asset with no digest an error.
	if target.ManifestSHA256 != "" {
		if options.verify == VerifyOff {
			return ErrVerifyDisabled
		}
		options.verify = VerifyRequired
	}

	// An empty version takes the release pinned by this yzma release. "latest" always
	// asks for the newest build, so it skips the pin.
	if target.Version == "" && DefaultVersion != "" {
		target.Version = DefaultVersion
	}

	autoVersion := target.Version == "" || target.Version == "latest"
	if autoVersion {
		latest, err := LlamaLatestVersion()
		if err != nil {
			return err
		}
		target.Version = latest
	}
	if err := VersionIsValid(target.Version); err != nil {
		return ErrInvalidVersion
	}

	// Only a tagged release needs a lookup on the llama.cpp release page. A nightly
	// tag names its own assets, so it stays on the llama-cpp-builder site, which does
	// not rate limit like the GitHub API.
	if target.UpstreamVersion == "" && IsTaggedRelease(target.Version) {
		upstream, err := LlamaNightlyTag(target.Version)
		if err != nil {
			return err
		}
		target.UpstreamVersion = upstream
	}

	assets, err := installAssets(ctx, target, dest, progress, resolver, options)
	if err == nil {
		return recordInstall(dest, target, assets)
	}

	// The newest release may still be building for this platform.
	if autoVersion && errors.Is(err, ErrFileNotFound) {
		previous, prevErr := LlamaPreviousVersion()
		if prevErr != nil {
			return err
		}
		target.Version = previous
		target.UpstreamVersion = ""
		if IsTaggedRelease(previous) {
			target.UpstreamVersion, prevErr = LlamaNightlyTag(previous)
			if prevErr != nil {
				return err
			}
		}
		assets, err := installAssets(ctx, target, dest, progress, resolver, options)
		if err != nil {
			return err
		}
		return recordInstall(dest, target, assets)
	}

	return err
}

// recordInstall leaves a record of what was installed, so [VerifyInstall] can check
// the files later.
func recordInstall(dest string, target Target, assets []Asset) error {
	return WriteInstallRecord(dest, InstallRecord{
		Tag:         target.Version,
		UpstreamTag: target.UpstreamVersion,
		Arch:        target.Arch.String(),
		OS:          target.OS.String(),
		Processor:   target.Processor.String(),
		Installed:   time.Now().UTC(),
		Assets:      assets,
	})
}

func installAssets(ctx context.Context, target Target, dest string, progress getter.ProgressTracker, resolver Resolver, options installOptions) ([]Asset, error) {
	assets, err := resolveAssets(ctx, target, resolver, options.verify)
	if err != nil {
		return nil, err
	}

	for _, asset := range assets {
		switch {
		case options.verify == VerifyOff, asset.SHA256 != "":
			// Nothing to say. A digest that is there is checked as it downloads.
		case options.verify == VerifyRequired:
			return nil, fmt.Errorf("%w: %s", ErrDigestMissing, asset.URL)
		case VerifyWarning != nil:
			VerifyWarning(asset.URL)
		}

		if err := getFunc(ctx, asset, dest, progress); err != nil {
			return nil, err
		}
	}

	return assets, nil
}

// resolveAssets asks a resolver for the assets to install. A resolver that reports
// digests is preferred, so a resolver that only has [Resolver] keeps working.
// [VerifyOff] takes the plain [Resolver], because a digest that nothing reads is not
// worth the fetch of a manifest.
func resolveAssets(ctx context.Context, target Target, resolver Resolver, verify VerifyPolicy) ([]Asset, error) {
	// A pinned manifest is the authority for every asset, whichever resolver named
	// them, so it is read here instead of in the resolver.
	if target.ManifestSHA256 != "" {
		return pinnedAssets(ctx, target, resolver)
	}

	if verify != VerifyOff {
		// The built-in resolver has a form that carries the context of the install.
		if d, ok := resolver.(defaultResolver); ok {
			return d.resolveAssets(ctx, target)
		}
		if assetResolver, ok := resolver.(AssetResolver); ok {
			return assetResolver.ResolveAssets(target)
		}
	}

	urls, err := resolver.Resolve(target)
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, len(urls))
	for i, url := range urls {
		assets[i] = Asset{URL: url}
	}

	return assets, nil
}

// pinnedAssets gives the assets to install when [Target.ManifestSHA256] pins the
// digest manifest. The resolver names the assets and the pinned manifest gives the
// digest of each one. An asset that the manifest does not name is an error.
func pinnedAssets(ctx context.Context, target Target, resolver Resolver) ([]Asset, error) {
	urls, err := resolver.Resolve(target)
	if err != nil {
		return nil, err
	}

	m, err := fetchManifest(ctx, target.Version, target.ManifestSHA256)
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, len(urls))
	for i, url := range urls {
		digest := m.digestFor(url)
		if digest == "" {
			return nil, fmt.Errorf("%w: %s is not in the pinned digests of %s", ErrDigestMissing, url, target.Version)
		}
		assets[i] = Asset{URL: url, SHA256: digest}
	}

	return assets, nil
}
