package query

type DataManager interface {
	GetData(root *Node) (string, error)
}
