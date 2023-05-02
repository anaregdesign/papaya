package nlp

type Dictionary struct {
	ID2Word map[int]string
	Word2ID map[string]int
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		ID2Word: make(map[int]string),
		Word2ID: make(map[string]int),
	}
}

func (d *Dictionary) AddWord(word string) {
	if _, ok := d.Word2ID[word]; !ok {
		id := len(d.ID2Word)
		d.ID2Word[id] = word
		d.Word2ID[word] = id
	}
}

func (d *Dictionary) AddWords(words []string) {
	for _, word := range words {
		d.AddWord(word)
	}
}

func (d *Dictionary) Words2BOW(words []string) BOW {
	d.AddWords(words)

	bow := NewBOW()
	for _, word := range words {
		id, ok := d.Word2ID[word]
		if ok {
			bow[id]++
		}
	}
	return bow
}

func (d *Dictionary) Words2CBOW(words []string, window int) []BOW {
	d.AddWords(words)

	bows := make([]BOW, len(words))
	for i, _ := range words {
		bows[i] = NewBOW()
		for j := i - window; j <= i+window; j++ {
			if j < 0 || j >= len(words) || j == i {
				continue
			}
			id, ok := d.Word2ID[words[j]]
			if ok {
				bows[i][id]++
			}
		}
	}
	return bows
}

type BOW map[int]int

func NewBOW() BOW {
	return make(map[int]int)
}
