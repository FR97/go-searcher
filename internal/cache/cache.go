package cache

// How many team does each term appear in a document
type FileTermFrequency struct {
	TF             map[string]uint
	TotalTermCount uint
	IndexTime      int64
}

// In how many files does a term appear
type TermFileFrequency map[string]uint

type Cache struct {
	FileToTermFreq map[string]FileTermFrequency `json:"file_to_term"`
	TermToFileFreq TermFileFrequency            `json:"term_to_file"`
}

func NewCache() Cache {
	return Cache{
		FileToTermFreq: map[string]FileTermFrequency{},
		TermToFileFreq: TermFileFrequency{},
	}
}
