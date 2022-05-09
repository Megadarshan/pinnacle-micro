package database

import (

	// "test/go_server/database"

	// tenant "tenant/proto"

	// tenant "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"
	// tenant "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"

	"fmt"
	"strconv"
	"strings"

	tenant "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"

	errors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
)

func GetStatusList() []*tenant.TenantStatus {
	cols := []string{
		"code",
		"CASE WHEN (CURRENT_DATE BETWEEN start_date and end_date) THEN true ELSE false END",
		"descr"}
	table := "status_master"
	where := []string{"type = 'TENANT'"}
	var list = []*tenant.TenantStatus{}

	query := "SELECT " + strings.Join(cols, ",") + " from " + table + " where " + strings.Join(where, " and ")
	log.Info("SQL : " + query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Info("DB Error : ", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		row := tenant.TenantStatus{}
		err = rows.Scan(&row.Code, &row.Status, &row.Descr)
		if err != nil {
			log.Info("ERROR white iterating rows from Query result. " + err.Error())
		}
		list = append(list, &row)
	}
	if err := rows.Err(); err != nil {
		log.Info(err.Error())
	}
	return list
}

func GetTypeList() []*tenant.TenantType {
	cols := []string{
		"type_code",
		"CASE WHEN (CURRENT_DATE BETWEEN start_date and end_date) THEN true ELSE false END",
		"type_descr"}
	table := "type_master"
	where := []string{"type = 'TENANT'"}
	var list = []*tenant.TenantType{}

	query := "SELECT " + strings.Join(cols, ",") + " from " + table + " where " + strings.Join(where, " and ")
	log.Info("SQL : " + query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Info("DB Error : ", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		row := tenant.TenantType{}
		err = rows.Scan(&row.Code, &row.Status, &row.Descr)
		if err != nil {
			log.Info("ERROR white iterating rows from Query result. " + err.Error())
		}
		list = append(list, &row)
	}
	if err := rows.Err(); err != nil {
		log.Info(err.Error())
	}
	return list
}

func AddTenant(newTenant *tenant.CreateTenantRequest) (string, error) {

	var exists bool
	checkQuery := fmt.Sprintf("SELECT exists (%s)", "SELECT tenant_name from tenant where lower(tenant_name) = lower('"+newTenant.Name+"')")
	exstErr := DB.QueryRow(checkQuery).Scan(&exists)
	if exstErr != nil {
		return "", errors.InternalServerError("Tenant.CreateTenant", "Unable select DB table 'tenant'. ERROR: %s", exstErr.Error())
	}

	if exists {
		return "", errors.Conflict("Tenant.CreateTenant", "Tenant name '%s' already Exists. Please use a new name", newTenant.Name)
	}

	cols := []string{
		"tenant_name",
		"logo",
		"status",
		"start_date",
		"end_date",
		"tenant_type",
	}
	table := "tenant"
	values := []string{
		"'" + newTenant.Name + "'",
		"'" + newTenant.Logo + "'",
		"'" + newTenant.Status + "'",
		"TO_DATE('" + newTenant.StartDate + "','YYYYMMDD')",
		"TO_DATE('" + newTenant.EndDate + "','YYYYMMDD')",
		"'" + newTenant.Type + "'",
	}
	var newId int

	query := "INSERT INTO " + table + " (" + strings.Join(cols, ",") + ") values (" + strings.Join(values, ",") + ") RETURNING tenant_id"
	err := DB.QueryRow(query).Scan(&newId)
	if err != nil {
		return "", errors.InternalServerError("Tenant.CreateTenant", "Unable to insert into DB. ERROR: %s", err.Error())
	}
	msg := "New tenant created with ID : " + strconv.Itoa(newId)
	return msg, nil
}

func GetAllTenant() []*tenant.TenantDetail {

	cols := []string{
		"tbl.tenant_id",
		"COALESCE(tbl.logo,'empty')",
		"tbl.tenant_name",
		"COALESCE((select descr from status_master where type = 'TENANT' and code = tbl.status),'Not Configured')",
		"to_char(tbl.start_date,'yyyymmdd')",
		"to_char(tbl.end_date,'yyyymmdd')",
		"COALESCE((select type_descr from type_master where type = 'TENANT' and type_code = tbl.tenant_type),'Not Configured')",
		"to_char(tbl.created_datetime,'YYYYMMDDHH24MISS')",
		"to_char(tbl.updated_datetime,'YYYYMMDDHH24MISS')"}
	table := "tenant tbl"
	var list = []*tenant.TenantDetail{}

	query := "SELECT " + strings.Join(cols, ",") + " from " + table
	rows, err := DB.Query(query)
	if err != nil {
		log.Info("DB Error : ", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		row := tenant.TenantDetail{}
		err = rows.Scan(&row.Id, &row.Logo, &row.Name, &row.Status, &row.StartDate, &row.EndDate, &row.TenantType, &row.CreatedOn, &row.LastModDate)

		if err != nil {
			log.Info(err.Error())
		}
		list = append(list, &row)
	}
	if err := rows.Err(); err != nil {
		log.Info(err.Error())
	}
	return list
}

func DeleteTenant(tenantId int32) (string, error) {

	sqlStatement := "DELETE FROM tenant WHERE tenant_id = $1"
	_, err := DB.Exec(sqlStatement, tenantId)
	if err != nil {
		return "", errors.InternalServerError("Tenant.CreateTenant", "Unable to insert into DB. ERROR: %s", err.Error())
	}
	return "Tenant with Id " + strconv.Itoa(int(tenantId)) + " deleted", nil

}
