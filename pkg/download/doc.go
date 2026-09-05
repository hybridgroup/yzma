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
// # Pinning the digests
//
// Every digest above comes from the same site that serves the asset. Anyone who can
// replace an asset can also replace the manifest that gives its digest. So the check
// finds a damaged download, but it does not find one that was put there on purpose.
//
// Keep the expected value where the release host cannot change it. Pin the digest of
// the manifest in the release that uses yzma. The version takes the digest as a
// suffix:
//
//	target := download.Target{Version: "b10785@sha256:" + wantManifest}
//	err := download.Install(ctx, target, libPath, download.ProgressTracker, nil)
//
// [Target.ManifestSHA256] holds the same value for a caller that does not want to
// build the string. [Get] and its variants take the suffix form in their version
// argument, and so does "yzma install --version".
//
// # Where the manifest digest comes from
//
// llama-cpp-builder publishes the manifest for a release as an asset of that release,
// named "<tag>.json". GitHub records the SHA-256 of every asset it stores, so that
// value is the manifest digest. The release notes for the tag print the complete pin,
// and https://hybridgroup.github.io/llama-cpp-builder/version.json carries it for the
// newest build.
//
// The digest of a platform archive is not the manifest digest. Those digests are what
// the manifest holds, one for each asset. The pin covers the file that names them.
//
// [ManifestDigest] and [PinnedVersion] read the value for a tag, so a program can
// record the pin it is about to use:
//
//	pin, err := download.PinnedVersion(ctx, tag)
//
// [DefaultVersion] already holds the complete pin for the release that this yzma
// release installs.
//
// What gets pinned is the manifest and not one archive. A version selects a different
// set of assets for each target, and some targets need more than one, so no single
// archive digest covers them all. The chain goes from the pin, to the manifest bytes,
// to the digest of every asset.
//
// A pin makes verification mandatory. The manifest bytes are checked before they are
// decoded. A manifest that cannot be read is an error, and not an install with no
// check. An asset that the manifest does not name stops the install with
// [ErrDigestMissing]. So a pin also works with a plain [Resolver], because [Install]
// reads the manifest itself. A resolver that names assets the release does not
// publish cannot be used with a pin. A pin that comes with [VerifyOff] gives
// [ErrVerifyDisabled], because the two ask for opposite things.
//
// "latest" and an empty version name whichever release is newest at the time, so
// neither can carry a digest.
//
// A version with no digest is not the same as no checking. It installs as it always
// has, with whatever digests the manifest gives, under the policy above. What it does
// not have is a value from outside the release host to check the manifest against.
//
// See INSTALL.md for the longer version.
package download
