package query

type DataReader interface {
	GetData(root *Node) (string, error)
}
