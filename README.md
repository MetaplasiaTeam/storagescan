# StorageScan

StorageScan is a contract variable query tool on EVM chain (ETH BSC HECO...) 
Through the getStorageAt() function, it allows us to get the value of the variable according to the slot address, including the private variable,
enjoy it!

Generate variable objects from the Solidity code of the contract, under development...

## Quick Start


- Network: Ropsten Testnet Network

- RPCNode: https://ropsten.infura.io/v3/

- Contract Address: 0xd9fc1c8ab7e6e06c5f67128b8000dce15f6baafa

contract solidity code
```solidity
contract example {
    // int type
    int8  private int1 = -8; // 0x0
    int128  private int2 = 128;  // 0x0
    int256 private int3  = 256; // 0x1
    
}



```

get contract variable value
```go
import "storagescan"

var RopstenTestNet = storagescan.GenGetStorageValueFunc(context.Background(), "https://ropsten.infura.io/v3", common.HexToAddress("0x24302f327764f94c15d930f5ac70d362b4a156f9"))
intValue := storagescan.SolidityInt{
        SlotIndex: common.HexToHash("0x0"),
        Length:    256,
        Offset:    0,
    }
log.Printf("value:%v\n", intValue.Value(RopstenTestNet))

```
