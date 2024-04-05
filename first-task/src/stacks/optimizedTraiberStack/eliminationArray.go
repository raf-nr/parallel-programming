package optimizedTraiberStack

import "math/rand"

type exchangersArray[T any] struct {
	exchangers []exchanger[T] // Array of elements through which exchange is carried out.
	replays    int            // The number of times try to make an exchange.
	power      int            // Exchanger array power.
}

func freshExchangersArray[T any](power, replays int) exchangersArray[T] {
	// Create a new instance of the array, filling all its elements with empty values.
	result := exchangersArray[T]{power: power, replays: replays}
	result.exchangers = make([]exchanger[T], result.power)
	for i := 0; i < result.power; i++ {
		fresh := exchanger[T]{}
		fresh.elem.Store(exchangerElem[T]{value: nil, state: empty})
		result.exchangers[i] = fresh
	}
	return result
}

func (eArray *exchangersArray[T]) visit(value *T) (*T, error) {
	index := rand.Intn(eArray.power) // Randomly select a cell in which the exchange will be attempted.
	return eArray.exchangers[index].exchange(value, eArray.replays)
}
