package errors

import "fmt"

var errorMessages = map[ErrorCode]string{
	UserNotFound:       "Usuario no encontrado",
	UserAlreadyExists:  "El usuario ya existe",
	UserInactive:       "El usuario esta inactivo",
	ValidationFailed:   "Fallo la validacion de datos",
	InvalidCredentials: "Credenciales incorrectas",
}

// GetMessage obtiene el mensaje para un código de error
func GetMessage(code ErrorCode) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "Error desconocido"
}

// GetMessageWithDetails retorna mensaje con detalles adicionales
func GetMessageWithDetails(code ErrorCode, details string) string {
	baseMessage := GetMessage(code)
	if details != "" {
		return fmt.Sprintf("%s: %s", baseMessage, details)
	}
	return baseMessage
}
