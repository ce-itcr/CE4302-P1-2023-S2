package main

import (
	"fmt"
	"sync"
)

// MemoriaPrincipal representa la memoria principal del sistema.
type MemoriaPrincipal struct {
	datos map[int]int
	mutex sync.Mutex
}

func NewMemoriaPrincipal() *MemoriaPrincipal {
	return &MemoriaPrincipal{
		datos: make(map[int]int),
	}
}

func (mp *MemoriaPrincipal) Leer(direccion int) int {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()
	valor, ok := mp.datos[direccion]
	if !ok {
		return 0
	}
	return valor
}

func (mp *MemoriaPrincipal) Escribir(direccion int, valor int) {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()
	mp.datos[direccion] = valor
}

func main() {
	memoriaPrincipal := NewMemoriaPrincipal()
	canal := make(chan int)

	// Hilo para la memoria principal
	go func() {
		direccion := 0x100
		valor := 42
		memoriaPrincipal.Escribir(direccion, valor)
		canal <- 1 // Notifica que la escritura se ha completado.
	}()

	// Acceso a la memoria principal desde el hilo del interconector
	direccion := 0x100
	resultado := 0

	// Espera a que el hilo de la memoria principal complete la escritura.
	<-canal

	// Lee desde la memoria principal.
	resultado = memoriaPrincipal.Leer(direccion)

	fmt.Printf("Valor en dirección %X: %d\n", direccion, resultado)
}
