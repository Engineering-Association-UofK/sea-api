package cert

import (
	"context"
	"database/sql"
	"fmt"
	"sea-api/internal/models"
)

func (c *CertificateService) VerifyCertificate(hash string) (*models.CertificateVerify, error) {
	cert, err := c.certificateRepository.GetByHash(hash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &models.CertificateVerify{
				Valid: false,
			}, nil
		}
		return nil, err
	}
	event, err := c.eventService.GetEventByID(cert.EventID)
	if err != nil {
		return nil, err
	}
	user, err := c.userRepo.GetByUserID(cert.UserID)
	if err != nil {
		return nil, err
	}
	valid := true
	status := cert.Status
	if status == models.CertRevoked {
		valid = false
	}

	return &models.CertificateVerify{
		Valid:     valid,
		ID:        fmt.Sprint(cert.ID),
		NameAr:    user.NameAr,
		NameEn:    user.NameEn,
		EventName: event.Name,
		Status:    status,
		Grade:     fmt.Sprintf("%.2f", cert.Grade),
		Outcomes:  event.Outcomes,
		EndDate:   event.EndDate,
		IssueDate: cert.IssueDate,
	}, nil
}

func (c *CertificateService) VerifyDocument(hash string) (*models.DocumentVerifyResponse, error) {
	doc, err := c.documentRepository.GetByHash(hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.DocumentVerifyResponse{
				Valid: false,
			}, nil
		}
		return nil, err
	}

	relations, err := c.documentRepository.GetRelationsByDocumentID(doc.ID)
	if err != nil {
		return nil, err
	}

	metadata, err := c.documentRepository.GetMetadataByDocumentID(doc.ID)
	if err != nil {
		return nil, err
	}

	var Map []models.DocumentMetadata
	for _, relation := range relations {
		switch relation.ObjectType {
		case models.ObjEvent:
			event, err := c.eventService.GetEventByID(relation.ObjectID)
			if err != nil {
				return nil, err
			}
			Map = append(Map, models.DocumentMetadata{
				Label: "Event",
				Value: event.Name,
			})
		case models.ObjCollaborator:
			collab, err := c.CollaboratorService.GetByID(context.Background(), relation.ObjectID)
			if err != nil {
				return nil, err
			}
			Map = append(Map, models.DocumentMetadata{
				Label: "Collaborator",
				Value: collab.NameEn,
			})
		}
	}
	for _, meta := range metadata {
		Map = append(Map, models.DocumentMetadata{
			Label: meta.Key,
			Value: meta.Value,
		})
	}

	return &models.DocumentVerifyResponse{
		Valid:     true,
		Type:      doc.Type,
		CreatedAt: doc.CreateAt,
		Details:   Map,
	}, nil
}
