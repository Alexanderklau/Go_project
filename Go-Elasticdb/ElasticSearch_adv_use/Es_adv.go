package main

import (
"context"
"fmt"
"reflect"

"github.com/olivere/elastic"
)

//Contentmapping 测试
var Contentmapping = `{
	"settings": {
	  "number_of_shards": 5,
	  "number_of_replicas": 1,
	  "codec": "best_compression",
	  "max_result_window": "100000000"
	},
	"mappings": {
	  "doc": {
		"properties": {
		  "name": {
			"type": "keyword"
		  },
		  "id": {
			"type": "keyword"
		  },
		  "data": {
			"type": "nested",
			"properties": {
			  "value": {
				"type": "text",
				"fields": {
				  "keyword": {
					"ignore_above": 256,
					"type": "keyword"
				  }
				}
			  },
			  "key": {
				"type": "text",
				"fields": {
				  "keyword": {
					"ignore_above": 256,
					"type": "keyword"
				  }
				}
			  }
			}
		  },
		  "size": {
			"type": "long"
		  },
		  "last_mod_time": {
			"format": "yyyy-MM-dd HH:mm:ss",
			"type": "date"
		  },
		  "user": {
			"type": "keyword"
		  },
		  "content": {
			"analyzer": "ik_smart",
			"term_vector": "with_positions_offsets",
			"type": "text"
		  }
		}
	  }
	}
  }`

type ContentEsInfo struct {
	Name        string     `json:"name"`
	ID          string     `json:"id"`
	Size        uint64     `json:"size"`
	LastModTime string     `json:"last_mod_time"`
	User        string     `json:"user"`
	Data        []DataType `json:"data"`
	Content     string     `json:"content"`
}

type DataType struct {
	Key   string `json:"key"`
	Value string `json:"value" `
}

//Elastic es的连接
type Elastic struct {
	Client *elastic.Client
	host   string
}

//Connect 连接Es
func Connect(ip string) (*Elastic, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(ip))
	if err != nil {
		return nil, err
	}
	_, _, err = client.Ping(ip).Do(context.Background())
	if err != nil {
		return nil, err
	}
	_, err = client.ElasticsearchVersion(ip)
	if err != nil {
		return nil, err
	}
	es := &Elastic{
		Client: client,
		host:   ip,
	}
	return es, nil
}

//InitES 初始化Es
func InitES() (*Elastic, error) {
	host := []string{"http://10.0.9.28:9200"}
	Eslistsnum := len(host)
	if Eslistsnum == 0 {
		return nil, fmt.Errorf("Cluster Not Es Node")
	}
	//创建新的连接
	for i, ip := range host {
		//判断是不是最后一个节点ip
		if (Eslistsnum - 1) != i {
			es, err := Connect(ip)
			//如果连接出错，则跳过
			if err != nil {
				fmt.Println(err)
				continue
			}
			return es, nil
		} else {
			es, err := Connect(ip)
			if err != nil {
				return nil, err
			}
			return es, nil
		}
	}
	return nil, nil
}

//CreateIndex 创建一个index
func (Es *Elastic) CreateIndex(index, mapping string) bool {
	// 判断索引是否存在
	exists, err := Es.Client.IndexExists(index).Do(context.Background())
	if err != nil {
		fmt.Printf("<CreateIndex> some error occurred when check exists, index: %s, err:%s", index, err.Error())
		return false
	}
	if exists {
		fmt.Printf("<CreateIndex> index:{%s} is already exists", index)
		return true
	}
	createIndex, err := Es.Client.CreateIndex(index).Body(mapping).Do(context.Background())
	if err != nil {
		fmt.Printf("<CreateIndex> some error occurred when create. index: %s, err:%s", index, err.Error())
		return false
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		fmt.Printf("<CreateIndex> Not acknowledged, index: %s", index)
		return false
	}
	return true
}

//Put 上传数据
func (Es *Elastic) Put(index string, typ string, bodyJSON interface{}) (bool, error) {
	_, err := Es.Client.Index().
		Index(index).
		Type(typ).
		BodyJson(bodyJSON).
		Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Printf("<Put> some error occurred when put.  err:%s", err.Error())
		return false, err
	}
	return true, nil
}

//GetMsg 获取Msg
func (Es *Elastic) GetMsg(indexname, typ string) {
	var contentinfo ContentEsInfo
	res, _ := Es.Client.Search(indexname).Type(typ).Do(context.Background())
	//从搜索结果中取数据的方法
	for _, item := range res.Each(reflect.TypeOf(contentinfo)) {
		if t, ok := item.(ContentEsInfo); ok {
			fmt.Println(t)
		}
	}
}

//ShieldAnotherfield 屏蔽指定字段
func (Es *Elastic) ShieldAnotherfield(indexname, typ string) {
	var contentinfo ContentEsInfo
	fsc := elastic.NewFetchSourceContext(true).Include("name", "type", "user", "size", "last_mod_time", "data")
	res, _ := Es.Client.Search(indexname).Type(typ).FetchSourceContext(fsc).Do(context.Background())
	//从搜索结果中取数据的方法
	for _, item := range res.Each(reflect.TypeOf(contentinfo)) {
		if t, ok := item.(ContentEsInfo); ok {
			fmt.Println(t)
		}
	}
}

//FromSize 翻页方法
func (Es *Elastic) FromSize(indexname, typ string, size, from int) {
	var contentinfo ContentEsInfo
	res, _ := Es.Client.Search(indexname).Type(typ).Size(size).From(from).Do(context.Background())
	//从搜索结果中取数据的方法
	for _, item := range res.Each(reflect.TypeOf(contentinfo)) {
		if t, ok := item.(ContentEsInfo); ok {
			fmt.Println(t)
		}
	}
}

//HighlightMsg 高亮方法
func (Es *Elastic) HighlightMsg(indexname, typ string, size, from int, keyword string) {
	// var contentinfo ContentEsInfo

	boolQ := elastic.NewBoolQuery()
	boolZ := elastic.NewBoolQuery()

	// 定义highlight
	highlight := elastic.NewHighlight()
	// 指定需要高亮的字段
	highlight = highlight.Fields(elastic.NewHighlighterField("content"))
	// 指定高亮的返回逻辑 <span style='color: red;'>...msg...</span>
	highlight = highlight.PreTags("<span style='color: red;'>").PostTags("</span>")

	escontent := elastic.NewMatchQuery("content", keyword)
	boolZ.Filter(boolQ.Should(escontent))

	res, _ := Es.Client.Search(indexname).Type(typ).Highlight(highlight).Query(boolZ).Do(context.Background())
	for _, highliter := range res.Hits.Hits {
		fmt.Println(highliter.Highlight["content"][0])
	}
}

//Precisesearch 精准检索
func (Es *Elastic) Precisesearch(indexname, typ string, size, from int, keyword string) {
	// var contentinfo ContentEsInfo

	boolQ := elastic.NewBoolQuery()
	boolZ := elastic.NewBoolQuery()

	// 定义highlight
	highlight := elastic.NewHighlight()
	// 指定需要高亮的字段
	highlight = highlight.Fields(elastic.NewHighlighterField("content"))
	// 指定高亮的返回逻辑 <span style='color: red;'>...msg...</span>
	highlight = highlight.PreTags("<span style='color: red;'>").PostTags("</span>")

	// 短句匹配
	escontent := elastic.NewMatchPhrasePrefixQuery("content", keyword).MaxExpansions(10)

	boolZ.Filter(boolQ.Should(escontent))

	res, _ := Es.Client.Search(indexname).Type(typ).Highlight(highlight).Query(boolZ).Do(context.Background())
	for _, highliter := range res.Hits.Hits {
		fmt.Println(highliter.Highlight["content"][0])
	}
}

func main() {
	es, err := InitES()
	if err != nil {
		return
	}
	es.HighlightMsg("content_test", "doc", 10, 0, "三檛")
	// 上传doc
	// z := ContentEsInfo{Name: "yk", ID: "ak47", Size: 10423, Data: []DataType{{Key: "1", Value: "2"}}, User: "yk123", LastModTime: "2020-01-01 12:00:00", Content: "dssads"}
	// es.Put("content_test", "doc", z)
	// 创建index
	// es.CreateIndex("content_test", Contentmapping)
}