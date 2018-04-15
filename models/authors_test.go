// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

func testAuthors(t *testing.T) {
	t.Parallel()

	query := Authors(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testAuthorsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = author.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Authors(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAuthorsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{author}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testAuthorsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := AuthorExists(tx, author.ID)
	if err != nil {
		t.Errorf("Unable to check if Author exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AuthorExistsG to return true, but got false.")
	}
}
func testAuthorsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	authorFound, err := FindAuthor(tx, author.ID)
	if err != nil {
		t.Error(err)
	}

	if authorFound == nil {
		t.Error("want a record, got nil")
	}
}
func testAuthorsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Authors(tx).Bind(author); err != nil {
		t.Error(err)
	}
}

func testAuthorsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Authors(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAuthorsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Authors(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAuthorsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	authorOne := &Author{}
	authorTwo := &Author{}
	if err = randomize.Struct(seed, authorOne, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}
	if err = randomize.Struct(seed, authorTwo, authorDBTypes, false, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = authorOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = authorTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func authorBeforeInsertHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorAfterInsertHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorAfterSelectHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorBeforeUpdateHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorAfterUpdateHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorBeforeDeleteHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorAfterDeleteHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorBeforeUpsertHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func authorAfterUpsertHook(e boil.Executor, o *Author) error {
	*o = Author{}
	return nil
}

func testAuthorsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Author{}
	o := &Author{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, authorDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Author object: %s", err)
	}

	AddAuthorHook(boil.BeforeInsertHook, authorBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	authorBeforeInsertHooks = []AuthorHook{}

	AddAuthorHook(boil.AfterInsertHook, authorAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	authorAfterInsertHooks = []AuthorHook{}

	AddAuthorHook(boil.AfterSelectHook, authorAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	authorAfterSelectHooks = []AuthorHook{}

	AddAuthorHook(boil.BeforeUpdateHook, authorBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	authorBeforeUpdateHooks = []AuthorHook{}

	AddAuthorHook(boil.AfterUpdateHook, authorAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	authorAfterUpdateHooks = []AuthorHook{}

	AddAuthorHook(boil.BeforeDeleteHook, authorBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	authorBeforeDeleteHooks = []AuthorHook{}

	AddAuthorHook(boil.AfterDeleteHook, authorAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	authorAfterDeleteHooks = []AuthorHook{}

	AddAuthorHook(boil.BeforeUpsertHook, authorBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	authorBeforeUpsertHooks = []AuthorHook{}

	AddAuthorHook(boil.AfterUpsertHook, authorAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	authorAfterUpsertHooks = []AuthorHook{}
}
func testAuthorsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx, authorColumnsWithoutDefault...); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAuthorsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = author.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := AuthorSlice{author}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testAuthorsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Authors(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	authorDBTypes = map[string]string{`ID`: `bigint`, `Name`: `text`}
	_             = bytes.MinRead
)

func testAuthorsUpdate(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err = author.Update(tx); err != nil {
		t.Error(err)
	}
}

func testAuthorsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	author := &Author{}
	if err = randomize.Struct(seed, author, authorDBTypes, true, authorColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, author, authorDBTypes, true, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(authorColumns, authorPrimaryKeyColumns) {
		fields = authorColumns
	} else {
		fields = strmangle.SetComplement(
			authorColumns,
			authorPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(author))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := AuthorSlice{author}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testAuthorsUpsert(t *testing.T) {
	t.Parallel()

	if len(authorColumns) == len(authorPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	author := Author{}
	if err = randomize.Struct(seed, &author, authorDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = author.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err := Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &author, authorDBTypes, false, authorPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Author struct: %s", err)
	}

	if err = author.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Author: %s", err)
	}

	count, err = Authors(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
