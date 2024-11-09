package common

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"

	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

func DbPageQuery[T any](ctx context.Context, client dapr.Client, page, pageSize int, orderField string, tableName string, idFieldName string, queryString string) (pageResult *PageGeneric[T], err error) {

	ret, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString+"&_count="+idFieldName, "get")
	if err != nil {
		Logger.Error("DbPageQuery invoke method get count error:", err)
		err = errors.WithMessage(err, "invoke method get count error")

		return
	}
	var total []Count
	err = json.Unmarshal(ret, &total)
	if err != nil {
		Logger.Error("DbPageQuery json unmarshal error:", err)
		err = errors.WithMessage(err, "DbPageQuery json unmarshal error:")
		return
	}
	var methods = "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + "?" + queryString + "&_page=" + strconv.Itoa(page) + "&_page_size=" + strconv.Itoa(pageSize)
	if orderField != "" {
		methods = methods + "&_order=" + orderField
	}
	ret, err = client.InvokeMethod(ctx, DB_SERVICE_NAME, methods, "get")
	if err != nil {
		Logger.Error("DbPageQuery invoke method error:", err)
		err = errors.WithMessage(err, "DbPageQueryinvoke method error:"+err.Error())
		return
	}
	var dataList []T
	dec := json.NewDecoder(bytes.NewReader(ret))
	dec.UseNumber()
	err = dec.Decode(&dataList)

	if err != nil {
		Logger.Error("DbPageQuery json unmarshal error:", err)
		err = errors.WithMessage(err, "DbPageQuery json unmarshal error:"+err.Error())
		return
	} else {

		p := PageGeneric[T]{
			Page:     page,
			PageSize: pageSize,
			Total:    total[0].Count,
			Items:    dataList,
		}
		pageResult = &p
	}
	return
}

func DbQuery[T any](ctx context.Context, client dapr.Client, tableName string, queryString string) (result []T, err error) {
	ret, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, "get")
	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, err.Error())
		err = errors.WithMessage(err, "dbQuery error")
		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret)) //避免int64精度丢失
	dec.UseNumber()
	err = dec.Decode(&result)

	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal dataList", err.Error())
		err = errors.WithMessage(err, "unmashal error")
		return
	}
	return
}

func DbGetOne[T any](ctx context.Context, client dapr.Client, tableName string, queryString string) (result *T, err error) {
	ret, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, "get")
	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, err.Error())
		err = errors.WithMessage(err, "dbQuery error")
		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret)) //避免int64精度丢失
	dec.UseNumber()
	var list = []T{}
	err = dec.Decode(&list)

	if len(list) > 0 {
		result = &list[0]
	}

	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal dataList", err.Error())
		err = errors.WithMessage(err, "unmashal error")
		return
	}
	return
}

func DbGetCount[T any](ctx context.Context, client dapr.Client, tableName string, queryString string) (result *T, err error) {
	ret, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, "get")
	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+queryString, err.Error())
		err = errors.WithMessage(err, "dbQuery error")
		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret)) //避免int64精度丢失
	dec.UseNumber()

	var datas = []T{}
	err = dec.Decode(&datas)
	if len(datas) > 0 {
		result = &datas[0]
	}

	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal dataList", err.Error())
		err = errors.WithMessage(err, "unmashal error")
		return
	}
	return
}

func DbUpsert[T any](ctx context.Context, client dapr.Client, data T, tableName string, primaryKeys string) (err error) {
	buf, err := json.Marshal(data)
	if err != nil {
		err = errors.WithMessage(err, "dbupsert marshal error")
		return
	}
	str := string(buf)
	str = escapeSqlSingleQuote(str)
	dataContent := &dapr.DataContent{
		ContentType: "text/json",
		Data:        []byte(str),
	}
	_, err = client.InvokeMethodWithContent(ctx, DB_SERVICE_NAME, "/upsert/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?keys="+primaryKeys+"&batch=false", "post", dataContent)

	return
}

func DbBatchUpsert[T any](ctx context.Context, client dapr.Client, datas []T, tableName string, primaryKeys string) (err error) {
	if len(datas) == 0 {
		return
	}
	buf, err := json.Marshal(datas)
	if err != nil {
		err = errors.WithMessage(err, "DbBatchUpsert marshal error")
		return
	}
	str := string(buf)
	str = escapeSqlSingleQuote(str)
	dataContent := &dapr.DataContent{
		ContentType: "text/json",
		Data:        []byte(str),
	}

	_, err = client.InvokeMethodWithContent(ctx, DB_SERVICE_NAME, "/upsert/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?keys="+primaryKeys+"&batch=true", "post", dataContent)

	return
}

func DbInsert[T any](ctx context.Context, client dapr.Client, data T, tableName string) (resp *Response, err error) {
	buf, err := json.Marshal(data)
	if err != nil {
		err = errors.WithMessage(err, "DbInsert marshal error")
		return
	}
	dataContent := &dapr.DataContent{
		ContentType: "text/json",
		Data:        buf,
	}
	ret, err := client.InvokeMethodWithContent(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName, "post", dataContent)
	if err != nil {
		return
	} else {
		resp = OK.WithData(ret)
	}
	return
}

func DbDelete(ctx context.Context, client dapr.Client, tableName string, idField string, id string) (err error) {

	_, err = client.InvokeMethod(ctx, DB_SERVICE_NAME, "/"+DBNAME+"/"+DB_SCHEMA+"/"+tableName+"?"+idField+"="+id, "delete")
	return err
}

func DbDeleteExpired(ctx context.Context, client dapr.Client, tableName string, timeField string, expireTime time.Time) (err error) {

	method := "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + "?" + timeField + "='$lt." + expireTime.Format("2006-01-02T15:04:05") + "'"

	_, err = client.InvokeMethod(ctx, DB_SERVICE_NAME, method, "delete")
	return
}

func DbDeleteByOps(ctx context.Context, client dapr.Client, tableName string, field []string, ops []string, val []any) error {
	query := "?"
	for i := 0; i < len(field); i++ {
		s, useQuot, err := getValue(val[i])
		if err != nil {
			return err
		}

		if useQuot {
			op := getOp(ops[i])
			if strings.HasPrefix(op, "$") {
				query = query + field[i] + "=" + op + "" + s + "&"
			} else {
				query = query + field[i] + "='" + op + "" + s + "'&"
			}
		} else {
			query = query + field[i] + "=" + getOp(ops[i]) + "" + s + "&"
		}

	}
	query = strings.TrimSuffix(query, "&")

	method := "/" + DBNAME + "/" + DB_SCHEMA + "/" + tableName + query

	_, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, method, "delete")
	return err
}
func getOp(op string) string {
	switch op {
	case ">":
		return "$gt."
	case "<":
		return "$lt."
	case ">=":
		return "$gte."
	case "<=":
		return "$lte."
	case "==":
		return "$eq."
	case "!=":
		return "$ne."
	case "in":
		return "$in."
	}
	return ""
}
func getValue(val any) (ret string, useQuot bool, err error) {
	switch val.(type) {
	case int:
		ret = strconv.Itoa(val.(int))
	case string:
		ret = val.(string)
	case float64:
		ret = strconv.FormatFloat(val.(float64), 'f', -1, 64)

	case int64:
		ret = strconv.FormatInt(val.(int64), 10)
	case int32:
		ret = strconv.FormatInt(int64(val.(int32)), 10)

	case bool:
		if val.(bool) {
			ret = "1"
		} else {
			ret = "0"
		}

	case time.Time:
		ret = val.(time.Time).Format("2006-01-02T15:04:05")
		useQuot = true
	default:
		err = errors.New("unknown type: " + reflect.TypeOf(val).String())
	}
	return
}

func DbBatchInsert[T any](ctx context.Context, client dapr.Client, val []T, tablename string) (err error) {
	if len(val) == 0 {
		return
	}

	data, err := json.Marshal(val)
	if err != nil {
		return
	}

	content := &dapr.DataContent{
		ContentType: "application/json",
		Data:        data,
	}
	_, err = client.InvokeMethodWithContent(ctx, DB_SERVICE_NAME, "/batch/"+DBNAME+"/"+DB_SCHEMA+"/"+tablename, "POST", content)
	if err != nil {
		return
	}
	return nil
}
func DbRefreshContinuousAggregate(ctx context.Context, client dapr.Client, name, start, end string) (err error) {
	sqlScript := "/_QUERIES/mv/refresh_continuous_aggregate?name=" + name + "&start=" + start + "&end=" + end

	_, err = client.InvokeMethod(ctx, DB_SERVICE_NAME, sqlScript, "post")
	return
}
func CustomSql[T any](ctx context.Context, client dapr.Client, selectField, fromField, whereField string) (result []T, err error) {
	selectField = url.QueryEscape(selectField)
	fromField = url.QueryEscape(fromField)
	whereField = url.QueryEscape(whereField)
	sqlScript := "/_QUERIES/table/custom_sql?select_field=" + selectField + "&from_field=" + fromField + "&where_field=" + whereField
	ret, err := client.InvokeMethod(ctx, DB_SERVICE_NAME, sqlScript, "get")
	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrServiceInvokeDB.Status, sqlScript, err.Error())
		err = errors.WithMessage(err, "dbQuery error")
		return
	}
	dec := json.NewDecoder(bytes.NewReader(ret)) //避免int64精度丢失
	dec.UseNumber()
	err = dec.Decode(&result)

	if err != nil {
		log.Printf("errno=%d, method=%s,error=%s\n", ErrListUnMashal.Status, "unMashal dataList", err.Error())
		err = errors.WithMessage(err, "unmashal error")
		return
	}

	return
}

func escapeSqlSingleQuote(src string) string {
	return strings.Replace(src, "'", "''", -1)
}
