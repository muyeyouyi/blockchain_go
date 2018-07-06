package blockchain
/**
	创建创世区块
 */
func NewGenesisBlock() *Block {
	return NewBlock("我是创世区块",[]byte{})
}

/**
	创建区块链
 */
func CreateBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

/**
	区块链-增加区块
 */
func (bc *BlockChain)AddBlock(data string)  {
	preHash := bc.Blocks[len(bc.Blocks)-1].Hash//上个区块哈希值
	block := NewBlock(data, preHash)//创建新区快
	bc.Blocks = append(bc.Blocks, block)
}

/**
	区块链
 */
type BlockChain struct {
	Blocks []*Block
}