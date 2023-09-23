package model

import (
	"github.com/emzola/venato/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProjectToProto converts a Project struct into a generated proto counterpart.
func ProjectToProto(p *Project) *gen.Project {
	return &gen.Project{
		Id:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		StartDate:     timestamppb.New(p.StartDate),
		TargetEndDate: timestamppb.New(p.TargetEndDate),
		ActualEndDate: timestamppb.New(p.ActualEndDate),
		CreatedOn:     timestamppb.New(p.CreatedOn),
		CreatedBy:     p.CreatedBy,
		ModifiedOn:    timestamppb.New(p.ModifiedOn),
		ModifiedBy:    p.ModifiedBy,
		Version:       p.Version,
	}
}

// ProjectFromProto converts a generated proto counterpart into a Project struct.
func ProjectFromProto(p *gen.Project) *Project {
	return &Project{
		ID:            p.Id,
		Name:          p.Name,
		Description:   p.Description,
		StartDate:     p.StartDate.AsTime(),
		TargetEndDate: p.TargetEndDate.AsTime(),
		ActualEndDate: p.ActualEndDate.AsTime(),
		CreatedOn:     p.CreatedOn.AsTime(),
		CreatedBy:     p.CreatedBy,
		ModifiedOn:    p.ModifiedOn.AsTime(),
		ModifiedBy:    p.ModifiedBy,
		Version:       p.Version,
	}
}
