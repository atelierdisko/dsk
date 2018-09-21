// Copyright 2018 Atelier Disko. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func NewRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)

	return &Repository{
		Repository: repo,
		lookup:     make(map[string]time.Time, 0),
		path:       path,
		ticker:     time.NewTicker(5 * time.Second),
		done:       make(chan bool),
	}, err
}

type Repository struct {
	sync.RWMutex
	*git.Repository

	// Lookup table, mapping file paths to modified times.
	lookup map[string]time.Time

	// Current head reference.
	head *plumbing.Reference

	// Root of the repository's worktree, containing the .git
	// directory.
	path string

	// Ticker which triggers a cache rebuild.
	ticker *time.Ticker

	// Quit channel, receiving true, when we are closed.
	done chan bool
}

func (r *Repository) StartCacher() {
	yellow := color.New(color.FgYellow)

	go func() {
		for {
			select {
			case <-r.ticker.C:
				if r.IsCacheStale() {
					if err := r.BuildCache(); err != nil {
						log.Print(yellow.Sprintf("Failed to rebuild repository cache: %s", err))
						continue
					}
				}
			case <-r.done:
				log.Print("Stopping repo cacher (received quit)...")
				return
			}
		}
	}()
}

func (r *Repository) StopCacher() {
	r.done <- true
}

func (r *Repository) Close() {
	r.ticker.Stop()
}

func (r *Repository) IsCacheStale() bool {
	r.RLock()
	defer r.RUnlock()

	if r.head == nil {
		return false
	}
	ref, _ := r.Head()
	return r.head.Hash() != ref.Hash()
}

// BuildCache will warm up the cache. So lookups for a file's modified time
// are speed up. Will cache modified time for all files and directories
// discovered in root, which is recursively walked.
//
// Implementation based upon snippet provided in:
// https://github.com/src-d/go-git/issues/604
//
// Also see:
// https://github.com/src-d/go-git/issues/417
// https://github.com/src-d/go-git/issues/826
func (r *Repository) BuildCache() error {
	r.Lock()
	defer r.Unlock()

	start := time.Now()
	pathsCached := make(map[string]bool, 0)

	ref, _ := r.Head()
	if ref == nil {
		log.Printf("No commits in repository %s, yet", r.path)
		return nil
	}
	r.head = ref

	err := filepath.Walk(r.path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			isRoot := filepath.Base(r.path) == f.Name()

			if IgnoreNodesRegexp.MatchString(f.Name()) && !isRoot {
				return filepath.SkipDir
			}
			return nil // Git only knows about files
		}
		rel, _ := filepath.Rel(r.path, path)
		pathsCached[rel] = false
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to walk directory tree %s: %s", r.path, err)
	}

	r.lookup = make(map[string]time.Time, 0)

	commits, err := r.Log(&git.LogOptions{From: r.head.Hash()})
	if err != nil {
		return err
	}
	defer commits.Close()

	var prevCommit *object.Commit
	var prevTree *object.Tree

Outer:
	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}
		currentTree, err := commit.Tree()
		if err != nil {
			return err
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		changes, err := currentTree.Diff(prevTree)
		if err != nil {
			return err
		}

		for _, c := range changes {
			if c.To.Name == "" {
				continue
			}
			if isCached, ok := pathsCached[c.To.Name]; !ok || isCached {
				// Not interested in this file.
				continue
			}
			r.lookup[c.To.Name] = prevCommit.Author.When
			pathsCached[c.To.Name] = true

			if len(r.lookup) >= len(pathsCached) {
				break Outer
			}
		}

		prevCommit = commit
		prevTree = currentTree
	}

	log.Printf("Built repository cache in %s", time.Since(start))
	return nil
}

// Modified considers any changes inside given path or its
// subdirectories as a change to the path.
func (r *Repository) Modified(path string) (time.Time, error) {
	// Fast path for files.
	if m, ok := r.lookup[path]; ok {
		return m, nil
	}

	var modified time.Time

	for p, m := range r.lookup {
		if !filepath.HasPrefix(p, path) {
			continue
		}
		if m.After(modified) {
			modified = m
		}
	}

	if !modified.IsZero() {
		return modified, nil
	}
	if r.head == nil {
		return modified, nil
	}
	// When there's only one commit no diffing has been taken place.
	// It can be assumed that this is an initial commit adding all
	// files.
	commit, err := r.CommitObject(r.head.Hash())
	if err != nil {
		return modified, err
	}
	return commit.Author.When, nil
}