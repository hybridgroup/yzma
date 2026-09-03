// Package download provides utilities for downloading both the llama.cpp
// libraries and also model files.
//
// [Get] and its variants install the build the built-in table picks for a platform.
// [Install] takes a [Target] plus an optional [Resolver], so an application can install
// builds the table does not name — an internal mirror, a local file, its own llama.cpp
// build, or another CUDA major version:
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
// # Checking what comes down
//
// [Install] checks the SHA-256 of each asset before it writes anything, against the
// digest manifest that llama-cpp-builder publishes for the release. An asset whose
// bytes do not agree stops the install with [ErrDigestMismatch], and nothing is
// extracted.
//
// The default policy is [VerifyIfAvailable]: an asset with no digest installs and
// [VerifyWarning] says so. A deployment that must know what it loads asks for more:
//
//	err := download.Install(ctx, target, libPath, download.ProgressTracker, nil,
//		download.WithVerify(download.VerifyRequired))
//
// A [Resolver] reports no digests, so [VerifyRequired] refuses one. Implement
// [AssetResolver] to give a digest for each asset.
//
// See INSTALL.md for the longer version.
package download
