package monke

import (
	"net/url"
	"slices"
	"sort"
)

type ArticleLookup struct {
	created    []*Article            // list of articles ordered by their created date, descending
	index      map[string]*Article   // map of articles by their identifiers
	categories map[string][]*Article // map of articles by their category
	tags       map[string][]*Article // map of articles by their tags
}

type ArticleFilter struct {
	Category string `query:"category"`
	Tag      string `query:"tag"`
	Page     int    `query:"page"`
}

func (f ArticleFilter) ToUrlValues() url.Values {
	var params url.Values = make(url.Values)
	if len(f.Category) > 0 {
		params.Set("category", f.Category)
	}
	if len(f.Tag) > 0 {
		params.Set("tag", f.Tag)
	}
	return params
}

func (db *ArticleLookup) Init(a []*Article) {
	db.created = a
	db.index = make(map[string]*Article)
	db.categories = make(map[string][]*Article)
	db.tags = make(map[string][]*Article)

	sort.Slice(db.created, func(i, j int) bool {
		return db.created[i].Created.After(db.created[j].Created)
	})

	for _, article := range db.created {
		db.index[article.Id] = article
		db.categories[article.Category] = append(db.categories[article.Category], article)
		for _, tag := range article.Tags {
			db.tags[tag] = append(db.tags[tag], article)
		}
	}
}

func (db *ArticleLookup) GetArticles(filter ArticleFilter, limit int, offset int) []*Article {
	if limit < 1 {
		limit = max(0, len(db.created)-offset)
	}
	result := make([]*Article, 0, limit)
	var input []*Article = nil
	if filter.Category != "" {
		input = db.categories[filter.Category]
	} else {
		input = db.created
	}
	if filter.Tag != "" {
		for _, a := range input {
			if slices.Contains(a.Tags, filter.Tag) {
				if offset > 0 {
					offset -= 1
					continue
				}
				result = append(result, a)
				if len(result) == limit {
					break
				}
			}
		}
	} else if offset < len(input) {
		end := min(len(input), offset+limit)
		result = input[offset:end]
	}
	return result
}

func (db *ArticleLookup) GetArticle(id string) *Article {
	return db.index[id]
}

func (db *ArticleLookup) GetTags() []string {
	result := make([]string, 0, len(db.tags))
	for k := range db.tags {
		result = append(result, k)
	}
	slices.Sort(result)
	return result
}

func (db *ArticleLookup) GetTagsSizes() map[string]int {
	result := make(map[string]int)
	for k, v := range db.tags {
		result[k] = len(v)
	}
	return result
}
