package common

import (
	"bytes"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CommonPageQuery[T any](w http.ResponseWriter, r *http.Request, client client.Client, tableName string, idFieldName string) {
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

		Logger.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal total", err.Error())

		return
	}
	var methods = "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + "?_page=" + page + "&_page_size=" + pageSize + "&" + qstr
	if orderField != "" {
		methods = methods + "&_order=" + orderField
	}
	ret, err = client.InvokeMethod(r.Context(), DB_SERVICE_NAME, methods, "get")
	if err != nil {
		HttpError(w, ErrServiceInvokeDB, http.StatusOK)

		Logger.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, methods, err.Error())

		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret)) //避免int64精度丢失
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
		Logger.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
		return
	}
	var tmpMapList []map[string]any
	err = json.Unmarshal(buf, &tmpMapList)
	if err != nil {
		HttpError(w, ErrService, http.StatusOK)
		Logger.Printf("errno=%d, method=%s,error=%s\n", ErrService.Status, "Marshal dataList", err.Error())
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
