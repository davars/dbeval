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

func testArticles(t *testing.T) {
	t.Parallel()

	query := Articles(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testArticlesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = article.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testArticlesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Articles(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testArticlesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ArticleSlice{article}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testArticlesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ArticleExists(tx, article.ID)
	if err != nil {
		t.Errorf("Unable to check if Article exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ArticleExistsG to return true, but got false.")
	}
}
func testArticlesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	articleFound, err := FindArticle(tx, article.ID)
	if err != nil {
		t.Error(err)
	}

	if articleFound == nil {
		t.Error("want a record, got nil")
	}
}
func testArticlesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Articles(tx).Bind(article); err != nil {
		t.Error(err)
	}
}

func testArticlesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Articles(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testArticlesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	articleOne := &Article{}
	articleTwo := &Article{}
	if err = randomize.Struct(seed, articleOne, articleDBTypes, false, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}
	if err = randomize.Struct(seed, articleTwo, articleDBTypes, false, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = articleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = articleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Articles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testArticlesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	articleOne := &Article{}
	articleTwo := &Article{}
	if err = randomize.Struct(seed, articleOne, articleDBTypes, false, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}
	if err = randomize.Struct(seed, articleTwo, articleDBTypes, false, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = articleOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = articleTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func articleBeforeInsertHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleAfterInsertHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleAfterSelectHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleBeforeUpdateHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleAfterUpdateHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleBeforeDeleteHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleAfterDeleteHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleBeforeUpsertHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func articleAfterUpsertHook(e boil.Executor, o *Article) error {
	*o = Article{}
	return nil
}

func testArticlesHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Article{}
	o := &Article{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, articleDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Article object: %s", err)
	}

	AddArticleHook(boil.BeforeInsertHook, articleBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	articleBeforeInsertHooks = []ArticleHook{}

	AddArticleHook(boil.AfterInsertHook, articleAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	articleAfterInsertHooks = []ArticleHook{}

	AddArticleHook(boil.AfterSelectHook, articleAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	articleAfterSelectHooks = []ArticleHook{}

	AddArticleHook(boil.BeforeUpdateHook, articleBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	articleBeforeUpdateHooks = []ArticleHook{}

	AddArticleHook(boil.AfterUpdateHook, articleAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	articleAfterUpdateHooks = []ArticleHook{}

	AddArticleHook(boil.BeforeDeleteHook, articleBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	articleBeforeDeleteHooks = []ArticleHook{}

	AddArticleHook(boil.AfterDeleteHook, articleAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	articleAfterDeleteHooks = []ArticleHook{}

	AddArticleHook(boil.BeforeUpsertHook, articleBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	articleBeforeUpsertHooks = []ArticleHook{}

	AddArticleHook(boil.AfterUpsertHook, articleAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	articleAfterUpsertHooks = []ArticleHook{}
}
func testArticlesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testArticlesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx, articleColumnsWithoutDefault...); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testArticlesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = article.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testArticlesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ArticleSlice{article}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testArticlesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Articles(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	articleDBTypes = map[string]string{`Body`: `text`, `ID`: `bigint`, `PublishedAt`: `timestamp with time zone`, `Title`: `text`}
	_              = bytes.MinRead
)

func testArticlesUpdate(t *testing.T) {
	t.Parallel()

	if len(articleColumns) == len(articlePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	if err = article.Update(tx); err != nil {
		t.Error(err)
	}
}

func testArticlesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(articleColumns) == len(articlePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	article := &Article{}
	if err = randomize.Struct(seed, article, articleDBTypes, true, articleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, article, articleDBTypes, true, articlePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(articleColumns, articlePrimaryKeyColumns) {
		fields = articleColumns
	} else {
		fields = strmangle.SetComplement(
			articleColumns,
			articlePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(article))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ArticleSlice{article}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testArticlesUpsert(t *testing.T) {
	t.Parallel()

	if len(articleColumns) == len(articlePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	article := Article{}
	if err = randomize.Struct(seed, &article, articleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = article.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Article: %s", err)
	}

	count, err := Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &article, articleDBTypes, false, articlePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Article struct: %s", err)
	}

	if err = article.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Article: %s", err)
	}

	count, err = Articles(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
