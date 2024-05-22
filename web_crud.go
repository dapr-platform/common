package common

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cast"
	"net/url"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func CommonPageQuery[T any](w http.ResponseWriter, r *http.Request, client dapr.Client, tableName string, idFieldName string) {
	var dataList []map[string]any
	var vars = r.URL.Query()
	var page = "1"
	var pageSize = "10"
	var orderField = ""
	var qstr = ""
	if vars != nil {
		if p, exists := vars["_page"]; exists {
			page = p[0]
		}
		if p, exists := vars["_page_size"]; exists {
			pageSize = p[0]
		}
		if p, exists := vars["_order"]; exists {
			orderField = p[0]
		}

		for k, v := range vars {
			if k == "_page" || k == "_page_size" || k == "_order" {
				continue
			}
			//val := strings.Replace(v[0], "%", "%25", -1)
			for _, vv := range v {
				val := url.QueryEscape(vv)
				qstr += k + "=" + val + "&"
			}
		}
		qstr = strings.TrimSuffix(qstr, "&")
		Logger.Debug("qstr:", qstr)

	}
	v, err := strconv.Atoi(page)
	if err != nil || v <= 0 {
		HttpError(w, ErrParam.AppendMsg("page param error"), http.StatusBadRequest)
		return
	}
	v, err = strconv.Atoi(pageSize)
	if err != nil || v <= 0 {
		HttpError(w, ErrParam.AppendMsg("pageSize param error"), http.StatusBadRequest)
		return
	}

	ret, err := client.InvokeMethod(r.Context(), DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?_count="+idFieldName+"&"+qstr, "get")
	if err != nil {
		HttpError(w, ErrServiceInvokeDB, http.StatusOK)

		return
	}
	var total []Count
	err = json.Unmarshal(ret, &total)
	if err != nil {
		HttpError(w, ErrReqBodyParse.AppendMsg(tableName+" total"), http.StatusOK)

		log.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal total", err.Error())

		return
	}
	var methods = "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + "?_page=" + page + "&_page_size=" + pageSize + "&" + qstr
	if orderField != "" {
		methods = methods + "&_order=" + orderField
	}
	ret, err = client.InvokeMethod(r.Context(), DB_SERVICE_NAME, methods, "get")
	if err != nil {
		HttpError(w, ErrServiceInvokeDB, http.StatusOK)

		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, methods, err.Error())

		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret))
	dec.UseNumber()
	var tmpList []T
	err = dec.Decode(&tmpList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		return
	}
	buf, err := json.Marshal(tmpList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		log.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
		return
	}
	var tmpMapList []map[string]any
	err = json.Unmarshal(buf, &tmpMapList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		log.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
		return
	}
	selectFields := vars.Get("_select")
	if selectFields != "" {
		arr := strings.Split(selectFields, ",")
		for _, mv := range tmpMapList {
			data := make(map[string]any, 0)
			for _, vv := range arr {
				data[vv] = mv[vv]
			}
			dataList = append(dataList, data)
		}
	} else {
		dataList = tmpMapList
	}
	if dataList == nil {
		dataList = make([]map[string]any, 0)
	}
	pagei, _ := strconv.Atoi(page)
	pageSizei, _ := strconv.Atoi(pageSize)
	p := Page{
		Page:     pagei,
		PageSize: pageSizei,
		Total:    total[0].Count,
		Items:    dataList,
	}
	HttpSuccess(w, OK.WithData(p))

}
func CommonQuery[T any](w http.ResponseWriter, r *http.Request, client dapr.Client, tableName string, idFieldName string) {
	var dataList []map[string]any
	var vars = r.URL.Query()
	var orderField = ""
	var qstr = ""
	if vars != nil {

		if p, exists := vars["_order"]; exists {
			orderField = p[0]
		}

		for k, v := range vars {
			if k == "_page" || k == "_page_size" || k == "_order" {
				continue
			}
			val := url.QueryEscape(v[0])
			qstr += k + "=" + val + "&"
		}
		qstr = strings.TrimSuffix(qstr, "&")
		Logger.Debug("qstr:", qstr)

	}

	var methods = "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + "?" + qstr
	if orderField != "" {
		methods = methods + "&_order=" + orderField
	}
	ret, err := client.InvokeMethod(r.Context(), DB_SERVICE_NAME, methods, "get")
	if err != nil {
		HttpError(w, ErrServiceInvokeDB, http.StatusOK)

		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, methods, err.Error())

		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret))
	dec.UseNumber()
	var tmpList []T
	err = dec.Decode(&tmpList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		return
	}
	buf, err := json.Marshal(tmpList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		log.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
		return
	}
	var tmpMapList []map[string]any
	err = json.Unmarshal(buf, &tmpMapList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		log.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
		return
	}
	selectFields := vars.Get("_select")
	if selectFields != "" {
		arr := strings.Split(selectFields, ",")
		for _, mv := range tmpMapList {
			data := make(map[string]any, 0)
			for _, vv := range arr {
				data[vv] = mv[vv]
			}
			dataList = append(dataList, data)
		}
	} else {
		dataList = tmpMapList
	}
	if dataList == nil {
		dataList = make([]map[string]any, 0)
	}
	HttpSuccess(w, OK.WithData(dataList))
}

func CommonUpsert(w http.ResponseWriter, r *http.Request, client dapr.Client, tableName string, keys string) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpError(w, ErrReqBodyRead.AppendMsg(err.Error()), http.StatusOK)
		return
	}
	defer r.Body.Close()
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		HttpError(w, ErrReqBodyParse.AppendMsg(err.Error()), http.StatusOK)
		return
	}
	buf, err := json.Marshal(jsonData)
	if err != nil {
		HttpError(w, ErrReqBodyParse.AppendMsg(err.Error()), http.StatusOK)
		return
	}
	dataContent := &dapr.DataContent{
		ContentType: "text/json",
		Data:        buf,
	}
	ret, err := client.InvokeMethodWithContent(r.Context(), DB_SERVICE_NAME, "/upsert/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?keys="+keys+"&batch=false", "post", dataContent)
	if err != nil {
		HttpError(w, ErrServiceInvokeDB.AppendMsg(err.Error()), http.StatusOK)
		return
	}
	HttpSuccess(w, OK.WithData(string(ret)))
	return
}

func CommonDelete(w http.ResponseWriter, r *http.Request, client dapr.Client, tableName string, idField string, queryField string) {

	id := chi.URLParam(r, queryField)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrParam.AppendMsg("id param error").ToBytes())
	}

	ret, err := client.InvokeMethod(r.Context(), DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+idField+"="+id, "delete")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(ErrServiceInvokeDB.AppendMsg(err.Error()).ToBytes())
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(OK.WithData(ret).ToBytes())
	}
}
func CommonGroupby(w http.ResponseWriter, r *http.Request, client dapr.Client, tableName string) (err error) {
	field := r.URL.Query().Get("_select")
	where := r.URL.Query().Get("_where")
	if field == "" {
		HttpError(w, ErrParam, http.StatusOK)
		return
	}
	selectStr := field + ",count(*) as sum"
	fromStr := tableName
	whereStr := ""
	if where != "" {
		whereStr = where + " group by " + field + " order by " + field
	} else {
		whereStr = "1=1 group by " + field + " order by " + field
	}

	data, err := CustomSql[map[string]any](r.Context(), client, selectStr, fromStr, whereStr)
	if err != nil {
		HttpError(w, ErrServiceInvokeDB.AppendMsg(err.Error()), http.StatusOK)
		return
	}
	result := make(map[string]any)
	for _, d := range data {
		result[field+"_"+cast.ToString(d[field])] = d["sum"]
	}

	HttpSuccess(w, OK.WithData(result))
	return
}
