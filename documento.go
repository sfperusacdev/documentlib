package documentlib

import (
	"context"
	"database/sql"

	"github.com/sfperusacdev/identitysdk"
	"github.com/sfperusacdev/identitysdk/services"
	"github.com/user0608/goones/errs"
	connection "github.com/user0608/pg-connection"
)

type DocumentStateManager struct {
	manager connection.StorageManager
	service *services.ExternalBridgeService
}

func New(manager connection.StorageManager, service *services.ExternalBridgeService) DocumentStateManager {
	return DocumentStateManager{manager: manager, service: service}
}
func (s *DocumentStateManager) getEstadoActual(ctx context.Context, document, codigoDocumento string) (string, error) {
	var tx = s.manager.Conn(ctx)
	var estado sql.NullString
	var rs = tx.Table(document).Select("estado").Where("codigo=?", codigoDocumento).Find(&estado)
	if rs.Error != nil {
		return "", errs.Pgf(rs.Error)
	}
	if !estado.Valid || estado.String == "" {
		const format = "El documento '%s' con el identificador '%s' no tiene un estado definido."
		var err = errs.BadRequestf(format, document, identitysdk.RemovePrefix(codigoDocumento))
		return "", err
	}
	return estado.String, nil
}
func (s *DocumentStateManager) cambiarEstado(ctx context.Context, documento, codigoDocumento, nuevo_estado string) error {
	var tx = s.manager.Conn(ctx)
	var rs = tx.Table(documento).Where("codigo=?", codigoDocumento).Update("estado", nuevo_estado)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	return nil
}

type documentStateComand struct {
	MontoMaximo           float64
	MontoMaximoAcumulado  float64
	MontoAcumuladoPeriodo string
	callback              func() error
}

func (h *documentStateComand) ConfirmarCambio() error { return h.callback() }

// documento es el nombre de tabla y codigoDocumento el primary key
func (s *DocumentStateManager) CambiarEstadoDocumento(ctx context.Context, documento, codigoDocumento, nuevo_estado string) (*documentStateComand, error) {
	estados, err := s.service.GetEstadosDocumentoSegunUser(ctx, documento)
	if err != nil {
		return nil, err
	}
	estadoActual, err := s.getEstadoActual(ctx, documento, codigoDocumento)
	if err != nil {
		return nil, err
	}
	if !estados.Contains(estadoActual) {
		const format = "Este usuario no está autorizado a cambiar el estado."
		return nil, errs.BadRequestDirect(format)
	}
	for _, itm := range estados {
		if itm.Estado == estadoActual && itm.IsFinal {
			const format = "El estado de este documento no puede ser modificado"
			return nil, errs.BadRequestDirect(format)
		} else if itm.Estado == estadoActual {
			break
		}
	}
	maxAcumulado, periodo := estados.MontoMaximoAcumulado(nuevo_estado)
	return &documentStateComand{
		callback:              func() error { return s.cambiarEstado(ctx, documento, codigoDocumento, nuevo_estado) },
		MontoMaximo:           estados.MontoMaximo(nuevo_estado),
		MontoMaximoAcumulado:  maxAcumulado,
		MontoAcumuladoPeriodo: periodo,
	}, nil
}
