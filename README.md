# DocumentStateManager 

## Descripción

`DocumentStateManager` gestiona el cambio de estados de documentos, validando si un usuario tiene permisos para cambiar el estado y actualizando la base de datos.

## Componentes Clave

### `DocumentStateManager`

Estructura que maneja el estado de los documentos. Se comunica con el gestor de almacenamiento y el servicio externo para validar y actualizar el estado.

#### Métodos:

* `New(manager connection.StorageManager, service *services.ExternalBridgeService) DocumentStateManager`:

  * Crea una nueva instancia de `DocumentStateManager`.

* `getEstadoActual(ctx context.Context, document, codigoDocumento string) (string, error)`:

  * Obtiene el estado actual del documento.

* `cambiarEstado(ctx context.Context, documento, codigoDocumento, nuevo_estado string) error`:

  * Cambia el estado del documento.

* `CambiarEstadoDocumento(ctx context.Context, documento, codigoDocumento, nuevo_estado string) (*documentStateComand, error)`:

  * Valida si el usuario puede cambiar el estado del documento y genera un comando para confirmar el cambio.

### `documentStateComand`

Comando para confirmar el cambio de estado de un documento, incluye información relevante como montos máximos.

#### Métodos:

* `ConfirmarCambio() error`:

  * Ejecuta el cambio de estado.

## Manejo de Errores

Se usan errores personalizados con `errs.Pgf` y `errs.BadRequestf` para generar mensajes de error detallados.

## Ejemplo de Uso

```go
package main

import (
	"context"
	"fmt"
	"log"
	"documentlib"
	"github.com/sfperusacdev/identitysdk/services"
	connection "github.com/user0608/pg-connection"
)

func main() {
	// Crear instancias
	var storageManager connection.StorageManager
	var externalService *services.ExternalBridgeService

	// Inicializar DocumentStateManager
	manager := documentlib.New(storageManager, externalService)

	// Cambiar estado de documento
	documento := "document_table"
	codigoDocumento := "123456"
	nuevoEstado := "approved"

	cmd, err := manager.CambiarEstadoDocumento(context.Background(), documento, codigoDocumento, nuevoEstado)
	if err != nil {
		log.Fatal(err)
	}

	// Confirmar cambio
	if err := cmd.ConfirmarCambio(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Estado del documento actualizado con éxito")
}
```

## Dependencias

* `github.com/sfperusacdev/identitysdk`
* `github.com/sfperusacdev/identitysdk/services`
* `github.com/user0608/pg-connection`
* `github.com/user0608/goones/errs`

## Licencia

Detalles de la licencia (si aplica).
