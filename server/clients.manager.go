package server

type ClientsManager struct {
	AllClients []*Client
}

func (*ClientsManager) BroadcastMessage() {

}

func (*ClientsManager) AddClient(*Client) {

}

func (*ClientsManager) RemoveClient(name string) {

}
