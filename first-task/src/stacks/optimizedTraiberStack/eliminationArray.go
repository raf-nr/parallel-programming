package optimizedTraiberStack

import "math/rand"

type exchangersArray[T any] struct {
	exchangers []exchanger[T]
	replays    int
	power      int
}

func freshExchangersArray[T any](power, replays int) exchangersArray[T] {
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
	index := rand.Intn(eArray.power)
	return eArray.exchangers[index].exchange(value, eArray.replays)
}
