package bitcoin_load_spike

type Logger interface {
	FilePrefix() string
	FileExtension() string
	Log(float64, float64, int)
	Outputs() []string
	Reset()
}
