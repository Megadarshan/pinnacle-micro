package handler

import (
	"context"

	token "github.com/Megadarshan/pinnacle-micro/managetoken/proto"
	"github.com/Megadarshan/pinnacle-micro/super-admin/tenant/database"
	tenant "github.com/Megadarshan/pinnacle-micro/super-admin/tenant/proto"

	"github.com/micro/micro/v3/service/client"
	log "github.com/micro/micro/v3/service/logger"
	// tenant "tenant/proto"
)

type Tenant struct{}

func TokenService() token.ManagetokenService {
	service := token.NewManagetokenService("managetoken", client.DefaultClient)
	return service
}

func (e *Tenant) GetTenantStatus(ctx context.Context, req *tenant.TenantStatusRequest, rsp *tenant.TenantStatusResponse) error {
	log.Info("Received Tenant.GetTenantStatus request**************************************************")
	rsp.StatusList = database.GetStatusList()
	// rsp.TenantList = resp
	return nil
}

func (e *Tenant) GetTenantType(ctx context.Context, req *tenant.TenantTypeRequest, rsp *tenant.TenantTypeResponse) error {
	log.Info("Received Tenant.GetTenantType request**************************************************")
	rsp.TypeList = database.GetTypeList()
	// rsp.TenantList = resp
	return nil
}

func (e *Tenant) CreateTenant(ctx context.Context, req *tenant.CreateTenantRequest, rsp *tenant.CreateTenantResponse) error {
	log.Info("Received Tenant.CreateTenant request**************************************************")
	msg, err := database.AddTenant(req)
	if err != nil {
		return err
	}
	rsp.Response = msg
	return err
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Tenant) ListTenants(ctx context.Context, req *tenant.ListTenantsRequest, rsp *tenant.ListTenantsResponse) error {
	log.Info("Received Tenant.ListTenants request**************************************************")
	rsp.TenantList = database.GetAllTenant()
	return nil
}

func (e *Tenant) UpdateTenant(ctx context.Context, req *tenant.UpdateTenantRequest, rsp *tenant.UpdateTenantResponse) error {
	log.Info("Received Tenant.UpdateTenant request**************************************************")
	// rsp = database.GetStatusList()
	// rsp.TenantList = resp
	return nil
}

func (e *Tenant) DeleteTenant(ctx context.Context, req *tenant.DeleteTenantRequest, rsp *tenant.DeleteTenantResponse) error {
	log.Info("Received Tenant.DeleteTenant request**************************************************")
	msg, err := database.DeleteTenant(req.TenantId)
	if err != nil {
		return err
	}
	rsp.Respnse = msg
	return err
}
