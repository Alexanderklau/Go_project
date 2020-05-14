package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
)

//Elastic es的连接
type Elastic struct {
	Client *elastic.Client
	host   string
}

//Connect 连接
func Connect(ip string) (*Elastic, error) {
	client, err := elastic.NewClient(elastic.SetURL(ip))
	if err != nil {
		return nil, err
	}
	_, _, err = client.Ping(ip).Do(context.Background())
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	_, err = client.ElasticsearchVersion(ip)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Elasticsearch version %s\n", esVersion)
	es := &Elastic{
		Client: client,
		host:   ip,
	}
	return es, nil
}

//InitES 初始化Es
func InitES() (*Elastic, error) {
	host := []string{"http://10.0.6.245:9200","http://10.0.6.246:9200","http://10.0.6.247:9200"}
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

//CreateIndex 创建一个index, mapping这个你可以不需要
func (Es *Elastic) CreateIndex(index, mapping string) bool {
	// 判断索引是否存在
	exists, err := Es.Client.IndexExists(index).Do(context.Background())
	if err != nil {
		fmt.Sprintf("<CreateIndex> some error occurred when check exists, index: %s, err:%s", index, err.Error())
		return false
	}
	if exists {
		fmt.Sprintf("<CreateIndex> index:{%s} is already exists", index)
		return true
	}
	createIndex, err := Es.Client.CreateIndex(index).Body(mapping).Do(context.Background())
	if err != nil {
		fmt.Sprintf("<CreateIndex> some error occurred when create. index: %s, err:%s", index, err.Error())
		return false
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		fmt.Sprintf("<CreateIndex> Not acknowledged, index: %s", index)
		return false
	}
	return true
}

//DelIndex 删除Index
func (Es *Elastic) DelIndex(index string) bool {
	// Delete an index.
	deleteIndex, err := Es.Client.DeleteIndex(index).Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Sprintf("<DelIndex> some error occurred when delete. index: %s, err:%s", index, err.Error())
		return false
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
		fmt.Sprintf("<DelIndex> acknowledged. index: %s", index)
		return false
	}
	return true
}

//Put 上传数据
func (Es *Elastic) Put(index string, typ string, bodyJSON interface{}) bool {
	_, err := Es.Client.Index().
		Index(index).
		Type(typ).
		BodyJson(bodyJSON).
		Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Sprintf("<Put> some error occurred when put.  err:%s", err.Error())
		return false
	}
	return true
}

//Del 删除指定id数据
func (Es *Elastic) Del(index, typ, id string) bool {
	del, err := Es.Client.Delete().
		Index(index).
		Type(typ).
		Id(id).
		Do(context.Background())
	if err != nil {
		// Handle error
		fmt.Sprintf("<Del> some error occurred when put.  err:%s", err.Error())
		return false
	}
	fmt.Sprintf("<Del> success, id: %s to index: %s, type %s\n", del.Id, del.Index, del.Type)
	return true
}

//Update 更新数据
func (Es *Elastic) Update(index, typ, id string, updateMap map[string]interface{}) bool {
	res, err := Es.Client.Update().
		Index(index).Type(typ).Id(id).
		Doc(updateMap).
		FetchSource(true).
		Do(context.Background())
	if err != nil {
		_ = fmt.Sprintf("<Update> some error occurred when update. index:%s, typ:%s, id:%s err:%s", index, typ, id, err.Error())
		return false
	}
	if res == nil {
		fmt.Sprintf("<Update> expected response != nil. index:%s, typ:%s, id:%s", index, typ, id)
		return false
	}
	if res.GetResult == nil {
		fmt.Sprintf("<Update> expected GetResult != nil. index:%s, typ:%s, id:%s", index, typ, id)
		return false
	}
	//data, _ := json.Marshal(res.GetResult.Source)
	//fmt.Println("<Update> update success. data:%s", data)
	return true
}


//GetTaskLogCount
func (Es *Elastic) GetTaskLogCount(index, starttime, endtime string) (int, error) {
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewRangeQuery("time").Gte(starttime), elastic.NewRangeQuery("time").Lte(endtime))
	//统计count
	count, err := Es.Client.Count(index).Type("doc").Query(boolQ).Do(context.Background())
	if err != nil {
		return 0, nil
	}
	return int(count), nil
}

//GetSourceByID
func (Es *Elastic) GetSourceByID(index, typ, esid string) (*elastic.GetResult, error) {
	source, err := Es.Client.Get().Index(index).Type(typ).Id(esid).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return source, nil
}

//////GetaskMsg
func (Es *Elastic) GetaskMsg(index, typ, starttime, endtime, keyword string, size, page int) (*elastic.SearchResult, error) {
	boolQ := elastic.NewBoolQuery()
	if keyword == "" {
		boolQ.Filter(elastic.NewRangeQuery("starttime").Gte(starttime), elastic.NewRangeQuery("starttime").Lte(endtime))
		res, err := Es.Client.Search(index).Type(typ).Query(boolQ).Size(size).From((page - 1) * size).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	keys := fmt.Sprintf("name:*%s*", keyword)
	boolQ.Filter(elastic.NewRangeQuery("starttime").Gte(starttime), elastic.NewRangeQuery("starttime").Lte(endtime), elastic.NewQueryStringQuery(keys))
	res, err := Es.Client.Search(index).Type(typ).Query(boolQ).Size(size).From((page - 1) * size).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}
