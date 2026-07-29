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
// See INSTALL.md for the longer version.
package download
