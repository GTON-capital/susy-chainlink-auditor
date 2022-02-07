package sources

type DataSource interface {
	GetFloat32(string) float32
	GetUint32(string) uint32
}
