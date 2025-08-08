package freshdesk

import (
	"context"
	"strings"
)

// ---------------------------------------------------
// Ticket Fields

const (
	TicketFieldIncludeSection = "section"

	TicketFieldTypeDefaultSubject     = "default_subject"
	TicketFieldTypeDefaultAgent       = "default_agent"
	TicketFieldTypeDefaultRequester   = "default_requester"
	TicketFieldTypeDefaultGroup       = "default_group"
	TicketFieldTypeDefaultCompany     = "default_company"
	TicketFieldTypeDefaultProduct     = "default_product"
	TicketFieldTypeDefaultPriority    = "default_priority"
	TicketFieldTypeDefaultSource      = "default_source"
	TicketFieldTypeDefaultStatus      = "default_status"
	TicketFieldTypeDefaultDescription = "default_description"
	TicketFieldTypeDefaultTicketType  = "default_ticket_type"
	TicketFieldTypeCustomDate         = "custom_date"
	TicketFieldTypeCustomDateTime     = "custom_date_time"
	TicketFieldTypeCustomDropdown     = "custom_dropdown"
	TicketFieldTypeCustomParagraph    = "custom_paragraph"
	TicketFieldTypeCustomText         = "custom_text"
	TicketFieldTypeCustomCheckbox     = "custom_checkbox"
	TicketFieldTypeCustomNumber       = "custom_number"
	TicketFieldTypeCustomDecimal      = "custom_decimal"
	TicketFieldTypeCustomFile         = "custom_file"
	TicketFieldTypeNestedField        = "nested_field"
)

// List All Ticket Fields
func (fd *Freshdesk) ListTicketFields(ctx context.Context, types ...string) ([]*TicketField, error) {
	url := fd.Endpoint("/admin/ticket_fields")
	if len(types) > 0 {
		s := strings.Join(types, ",")
		url += "?type=" + s
	}

	fields := []*TicketField{}
	err := fd.DoGet(ctx, url, &fields)
	return fields, err
}

func (fd *Freshdesk) CreateTicketField(ctx context.Context, tf *TicketFieldCreate) (*TicketField, error) {
	url := fd.Endpoint("/admin/ticket_fields")
	result := &TicketField{}
	if err := fd.DoPost(ctx, url, tf, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTicketField View a Ticket Field
// include: conversations, requester, company, stats
func (fd *Freshdesk) GetTicketField(ctx context.Context, fid int64, include ...string) (*TicketField, error) {
	url := fd.Endpoint("/admin/ticket_fields/%d", fid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}

	result := &TicketField{}
	err := fd.DoGet(ctx, url, result)
	return result, err
}

func (fd *Freshdesk) UpdateTicketField(ctx context.Context, fid int64, field *TicketFieldUpdate) (*TicketField, error) {
	url := fd.Endpoint("/admin/ticket_fields/%d", fid)
	result := &TicketField{}
	if err := fd.DoPut(ctx, url, field, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) DeleteTicketField(ctx context.Context, fid int64) error {
	url := fd.Endpoint("/admin/ticket_fields/%d", fid)
	return fd.DoDelete(ctx, url)
}
