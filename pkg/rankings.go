package pkg

// TODO:
// Propably make the get the proximity score from the average proximity of the elements in that content
type Ranking struct {
	NoteLocation   string
	Frequency      int
	Proximityscore float32
	Keywords       []string
	Searchscore    float32
	Tags           []string
}
