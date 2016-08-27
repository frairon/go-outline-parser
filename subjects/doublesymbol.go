package subjects

// global symbol without receiver
func symbol() {}

type s float64

// global symbol with one receiver
func (s *s) symbol() {}

type t float64

// global symbol with receiver
func (t *t) symbol() {}
