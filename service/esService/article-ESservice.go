package esService

import (
	"context"
	"github.com/olivere/elastic/v7"
	"reflect"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
)

type ArticleESservice struct {
}

func (aes *ArticleESservice) FuzzyArticlesQuery(pageNum int, fuzVal, typeVal string) ([]mysqlModel.Article, error) {
	offSetter := sqlUtil.PageNumToSqlPage(pageNum, config.PageSize)
	//多字段精确匹配
	multiMatchPhraseQuery := elastic.NewBoolQuery().Should(
		elastic.NewMatchPhraseQuery("title", fuzVal),
		elastic.NewMatchPhraseQuery("content", fuzVal),
	)
	//工作类型匹配
	//设置必须匹配且并列匹配
	mustBoolQuery := elastic.NewBoolQuery().Must()
	matchATQuery := elastic.NewTermQuery("art_type", typeVal)
	matchStatusQuery := elastic.NewTermQuery("status", 1)
	matchShowQuery := elastic.NewTermQuery("show", 1)
	sortCreateTimeQuery := elastic.NewFieldSort("create_time").Desc()
	mustBoolQuery.Must(matchATQuery, matchStatusQuery, matchShowQuery, multiMatchPhraseQuery)
	//search
	searchByPhrase, err := dao.ESClient.Search().Index(config.ArticleESIndex).
		SortBy(sortCreateTimeQuery).
		From(offSetter).Size(config.PageSize).
		Query(mustBoolQuery).
		Do(context.Background())
	articles := []mysqlModel.Article{}
	if err != nil {
		return articles, err
	}

	for _, item := range searchByPhrase.Each(reflect.TypeOf(mysqlModel.Article{})) {
		articles = append(articles, item.(mysqlModel.Article))
	}
	return articles, nil
}