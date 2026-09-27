// Package download provides utilities for downloading both the llama.cpp
// libraries and also model files.
//
// [Get] and its variants install the build that the built-in table picks for a platform.
// [Install] takes a [Target] plus an optional [Resolver], so an application can install
// builds the table does not name — an internal mirror, a local file, or its own
// llama.cpp build:
//
//	resolver := download.ResolverFunc(func(t download.Target) ([]string, error) {
//		if t.OS == download.Linux && t.Processor == download.CUDA {
//			return []string{mirrorURL(t.Version)}, nil
//		}
//		return download.DefaultResolver.Resolve(t)
//	})
//
//	err := download.Install(ctx, target, libPath, download.ProgressTracker, resolver)
//
// An empty [Target.Version] takes [DefaultVersion], the llama.cpp release this yzma
// release was tested with. "latest" always gets the most recent nightly build.
//
// A Linux CUDA install uses [Target.CUDAVersion]. [WithCUDAVersion] sets it for callers
// that pass strings, and [HasCUDA] reports it for the machine. An empty value means
// CUDA 12 on ARM64 and CUDA 13 on AMD64. The processors [CUDA12] and [CUDA13] select a
// specific CUDA release instead.
//
// # Checking downloads
//
// Before writing anything, [Install] checks the SHA-256 of each asset against the
// digest manifest that llama-cpp-builder publishes for the release. If an asset does
// not match, the install stops with [ErrDigestMismatch] and nothing is extracted.
//
// The default policy is [VerifyIfAvailable]. An asset with no digest still installs
// and [VerifyWarning] reports it. A deployment that must know what it loads can
// require more:
//
//	err := download.Install(ctx, target, libPath, download.ProgressTracker, nil,
//		download.WithVerify(download.VerifyRequired))
//
// A [Resolver] reports no digests, so [VerifyRequired] rejects it. Implement
// [AssetResolver] to provide a digest for each asset.
//
// # Pinning the digests
//
// Every digest above comes from the same site that serves the asset. Anyone who can
// replace an asset can also replace the manifest that lists its digest. So the check
// catches a corrupted download, but not a deliberately replaced one.
//
// Keep the expected value where the release host cannot change it, by pinning the
// manifest digest in the release of your program that uses yzma. The version takes
// the digest as a suffix:
//
//	target := download.Target{Version: "b10785@sha256:" + wantManifest}
//	err := download.Install(ctx, target, libPath, download.ProgressTracker, nil)
//
// [Target.ManifestSHA256] holds the same value for callers that do not want to build
// the string. [Get] and its variants accept the suffix form in their version argument,
// and so does "yzma install --version".
//
// # Where the manifest digest comes from
//
// llama-cpp-builder publishes the manifest for a release as an asset of that release,
// named "<tag>.json". GitHub records the SHA-256 of every asset it stores, so that
// value is the manifest digest. The release notes for the tag print the full pin,
// and https://hybridgroup.github.io/llama-cpp-builder/version.json has it for the
// newest build.
//
// The manifest digest is not the digest of a platform archive. The manifest holds
// those, one per asset, and the pin covers the manifest itself.
//
// [ManifestDigest] and [PinnedVersion] look up the value for a tag, so a program can
// record the pin it is about to use:
//
//	pin, err := download.PinnedVersion(ctx, tag)
//
// [DefaultVersion] already holds the full pin for the release this yzma release
// installs.
//
// The pin covers the manifest, not a single archive. A version selects different
// assets for each target, and some targets need more than one, so no single archive
// digest covers them all. The chain goes from the pin, to the manifest bytes, to the
// digest of every asset.
//
// A pin makes verification mandatory. The manifest bytes are checked before they are
// decoded. A manifest that cannot be read is an error, not an unchecked install. An
// asset missing from the manifest stops the install with [ErrDigestMissing]. A pin
// also works with a plain [Resolver], because [Install] reads the manifest itself,
// but not with a resolver that returns assets the release does not publish. A pin
// with [VerifyOff] returns [ErrVerifyDisabled], because the two ask for opposite things.
//
// "latest" and an empty version name whichever release is newest at the time, so
// neither can carry a digest.
//
// A version with no digest still gets checked. It installs as it always has, using
// the digests in the manifest under the policy above. It just has no value from
// outside the release host to check the manifest against.
//
// # The manifest of an install
//
// [Install] saves the manifest it read in the library directory as yzma-manifest.json,
// and stores its digest in the install record. [VerifyInstall] reads that copy and
// checks it against the pin the operator gives, or against the recorded digest when
// there is no pin, so checking an installed release needs no network. If the manifest
// is missing or does not match, the check fetches it as before and saves the result.
//

// See https://yzma.ai/docs/guides/programmatic-install/ for the longer version.
package download
