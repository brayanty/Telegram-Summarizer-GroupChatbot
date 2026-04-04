package main

import (
	"fmt"
	"strings"
	"sync"
)

// MessageMember representa un mensaje de un usuario
type MessageMember struct {
	Name    string
	Message string
}

// String formatea el mensaje
func (m MessageMember) String() string {
	return fmt.Sprintf("*%s:* %s", m.Name, m.Message)
}

// ChatBuffer almacena los últimos N mensajes de un chat
type ChatBuffer struct {
	messages []MessageMember // Slice circular de mensajes
	maxSize  int             // Tamaño máximo
	index    int             // Índice actual para insertar
	count    int             // Cantidad real de mensajes guardados
	mu       sync.Mutex      // Para operaciones seguras en concurrencia
}

// NewChatBuffer crea un nuevo buffer con el tamaño especificado
func NewChatBuffer(size int) *ChatBuffer {
	return &ChatBuffer{
		messages: make([]MessageMember, size),
		maxSize:  size,
		index:    0,
		count:    0,
	}
}

// Add agrega un nuevo mensaje al buffer (sobrescribe el más antiguo si está lleno)
func (cb *ChatBuffer) Add(name, message string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Crear el mensaje
	msg := MessageMember{
		Name:    name,
		Message: message,
	}

	// Guardar en la posición actual
	cb.messages[cb.index] = msg

	// Avanzar el índice circularmente
	cb.index = (cb.index + 1) % cb.maxSize

	// Incrementar el contador hasta llegar al máximo
	if cb.count < cb.maxSize {
		cb.count++
	}
}

// GetAll devuelve todos los mensajes guardados en orden cronológico
func (cb *ChatBuffer) GetAll() []MessageMember {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	result := make([]MessageMember, cb.count)

	if cb.count == cb.maxSize {
		/* Buffer lleno: los mensajes están desde index hasta el
		final, luego desde 0 hasta index-1 */
		for i := 0; i < cb.count; i++ {
			pos := (cb.index + i) % cb.maxSize
			result[i] = cb.messages[pos]
		}
	} else {
		// Buffer no lleno: solo los primeros count elementos
		for i := 0; i < cb.count; i++ {
			result[i] = cb.messages[i]
		}
	}

	return result
}

// GetFormattedMessages devuelve los mensajes como string formateado
func (cb *ChatBuffer) GetFormattedMessages() string {
	messages := cb.GetAll()
	if len(messages) == 0 {
		return ""
	}

	var builder strings.Builder
	for _, msg := range messages {
		builder.WriteString(msg.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

// GetStats devuelve estadísticas del buffer (mensajes)
func (cb *ChatBuffer) GetStats() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return fmt.Sprintf("📊 *Estadísticas del mensajes:*\n- Mensajes guardados: %d/%d\n- Último índice: %d",
		cb.count, cb.maxSize, cb.index)
}

// Clear vacía completamente el buffer
func (cb *ChatBuffer) Clear() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.messages = make([]MessageMember, cb.maxSize)
	cb.index = 0
	cb.count = 0
}
