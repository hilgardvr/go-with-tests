package maps

const (
	ErrNotFound = DictionaryErr("could not find word")
	ErrDuplicateEntry = DictionaryErr("entry already exists")
	ErrWordDoesNotExist = DictionaryErr("word doesn't already exist")
)

type Dictionary map[string]string

type DictionaryErr string

func (e DictionaryErr)Error() string {
	return string(e)
}

func (d Dictionary)Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", ErrNotFound
	}
	return def, nil
}

func (d Dictionary)Update(word, def string) error {
	_, ok := d[word]
	if !ok {
		return ErrWordDoesNotExist
	} else { 
		d[word] = def
		return nil
	}
}

func (d Dictionary)Delete(word string) {
	delete(d, word)
}

func (d Dictionary)Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		d[word] = def
	case nil:
		return ErrDuplicateEntry
	default:
		return err
	}
	return nil
}