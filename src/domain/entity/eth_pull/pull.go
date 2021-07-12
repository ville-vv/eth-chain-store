package eth_pull

import "context"

type PullChainServer struct {
}

func (p *PullChainServer) Scheme() string {
	return "PullChainServer"
}

func (p *PullChainServer) Init() error {
	return nil
}

func (p *PullChainServer) Start() error {

	return nil
}

func (p *PullChainServer) pull() {

}

func (p *PullChainServer) Exit(ctx context.Context) error {
	return nil
}
