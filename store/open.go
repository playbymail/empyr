// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
	_ "modernc.org/sqlite"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

var (
	//go:embed migrations/*.sql
	migrations embed.FS

	//go:embed sqlc/schema.sql
	schemaDDL string
)

// Create creates a new store at the given path.
// It returns an error if the path does not exist or store already exists.
func Create(path string) error {
	scripts, err := loadMigrations()
	if err != nil {
		return err
	}

	// input must be an absolute path and not already exist
	if !filepath.IsAbs(path) {
		return ErrInvalidPath
	} else if stdlib.IsExists(path) {
		return ErrAlreadyExists
	}

	// create the store
	log.Printf("store: create: %s\n", path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("store: create: %v\n", err)
		return err
	}
	defer db.Close()

	// confirm that the store has foreign keys enabled
	if rows, err := db.Exec("PRAGMA" + " foreign_keys = ON"); err != nil {
		log.Printf("store: create: foreign keys are disabled\n")
		log.Printf("store: create: %v\n", err)
		return ErrForeignKeysDisabled
	} else if rows == nil {
		log.Printf("store: create: foreign keys pragma failed\n")
		return ErrPragmaReturnedNil
	}

	// run the schema DDL
	if _, err := db.Exec(schemaDDL); err != nil {
		log.Printf("store: create: %v\n", err)
		return err
	}

	// run the initialization scripts
	for _, script := range scripts {
		if _, err := db.Exec(script.script); err != nil {
			log.Printf("store: create: %d: %v\n", script.version, err)
			return errors.Join(ErrCreateSchema, err)
		}
		log.Printf("store: create: %d: %s\n", script.version, script.comment)
	}

	log.Printf("store: created %q\n", path)
	return nil
}

// Open opens an existing store.
// It returns an error if the path is invalid or the store does not exist.
// Caller must call Close() when done.
func Open(path string, ctx context.Context) (*Store, error) {
	// it is an error if the store isn't an absolute path or doesn't already exist
	if !filepath.IsAbs(path) {
		return nil, ErrInvalidPath
	} else if !stdlib.IsFileExists(path) {
		return nil, ErrInvalidPath
	}

	log.Printf("store: open: %s\n", path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// confirm that the database has foreign keys enabled
	if rows, err := db.Exec("PRAGMA" + " foreign_keys = ON"); err != nil {
		_ = db.Close()
		log.Printf("store: open: foreign keys are disabled\n")
		log.Printf("store: create: %v\n", err)
		return nil, ErrForeignKeysDisabled
	} else if rows == nil {
		_ = db.Close()
		log.Printf("store: open: foreign keys pragma failed\n")
		return nil, ErrPragmaReturnedNil
	}

	// return the store.
	return &Store{path: path, db: db, ctx: ctx, q: sqlc.New(db)}, nil
}

func (s *Store) Close() error {
	var err error
	if s != nil {
		if s.db != nil {
			err = s.db.Close()
			s.db = nil
		}
	}
	return err
}

type migrationScript struct {
	version int
	comment string
	script  string
}

func loadMigrations() (scripts []migrationScript, err error) {
	re, err := regexp.Compile(`^(\d{12})_([a-zA-Z0-9_])+\.sql$`)
	if err != nil {
		return nil, err
	}

	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// verify pattern: YYYYMMDDHHMM_comment.sql
		match := re.FindSubmatch([]byte(name))
		if match == nil {
			continue
		}
		version, err := strconv.Atoi(string(match[1]))
		if err != nil {
			return nil, err
		}
		comment := string(match[2])
		content, err := migrations.ReadFile("migrations/" + name)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, migrationScript{version: version, comment: comment, script: string(content)})
	}

	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].version < scripts[j].version
	})

	return scripts, nil
}
