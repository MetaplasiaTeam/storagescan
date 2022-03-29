# StorageScan

StorageScan is a contract variable query tool on EVM chain (ETH BSC HECO...) 
Through the getStorageAt() function, it allows us to get the value of the variable according to the slot address, including the private variable,
enjoy it!

Generate variable objects from the Solidity code of the contract, under development...

## Quick Start

examples
- Network: Ropsten Testnet Network

- RPCNode: https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161

- Contract Address: 0xd9fc1c8ab7e6e06c5f67128b8000dce15f6baafa

contract solidity code
```solidity
contract StorageScan {
    // int type
    int8  private int1 = - 8; // 0x0
    int128  private int2 = 128;  // 0x0
    int256 private int3 = 256; // 0x1

    // uint type
    uint8  private uint1 = 8; // 0x2
    uint128 private  uint2 = 128; // 0x2
    uint256 private uint3 = 0x123456789abcef1; // 0x3

    // bool type
    bool private bool1 = true; // 0x4
    bool private bool2 = false; // 0x4

    // string type
    string private string1 = "abc"; // 0x5
    string private string2 = "solidity is an object-oriented, high-level language for implementing smart contracts.";

    // bytes typeva
    bytes1 private b1 = "a"; // 0x7
    bytes8  private b2 = "byte2"; //0x7
    bytes32 private b3 = "string bytes cannot exceed 32"; //0x8

    // address type
    address private addr1 = 0x2729E5DFDeeCB92C884470EF6CaD9e844e34502D; // 0x9



    // struct type
    struct Entity {
        uint id;  // //0xa
        string value; // //0xb
    }

    Entity i; // //0xa


    // slice value
    uint8[] private slice1 = [1, 2, 3, 4, 5]; // 0xc
    uint256[] private slice2 = [256, 257, 258, 259, 260]; // // 0xd
    bool[] private slice3 = [true, false, false, true, false]; // 0xe
    string[] private slice4 = ["abc", "solidity is an object-oriented, high-level language for implementing smart contracts."]; //0xf
    Entity[] private slice5; // 0x10


    // array value
    uint8[5] private array1 = [1, 2, 3, 4, 5]; // 0x11
    uint256[5] private array2 = [256, 257, 258, 259, 260]; // // 0x12-0x16
    bool[5] private array3 = [true, false, false, true, false]; // 0x17
    string[2] private array4 = ["abc", "solidity is an object-oriented, high-level language for implementing smart contracts."];//0x18-0x19
    Entity[2] private array5; // 0x1a-0x1d



    // mapping value
    mapping(uint256 => string)  private   mapping1; // 0x1e
    mapping(string => uint256)  private   mapping2; // 0x1f
    mapping(address => uint256) private   mapping3; // 0x20
    mapping(int256 => uint256)  private   mapping4; // 0x21
    mapping(bytes1 => uint256)  private   mapping5; // 0x22
    mapping(uint256 => Entity)  private   mapping6; // 0x23



    constructor() {
        i.id = 1;
        i.value = "entity";

        slice5.push(Entity(1, "slice50"));
        slice5.push(Entity(2, "slice51"));

        array5[0] = Entity(1, "arry50");
        array5[1] = Entity(2, "array51");

        mapping1[1] = "mapping1";
        mapping2["mapping2"] = 1;
        mapping3[0x2729E5DFDeeCB92C884470EF6CaD9e844e34502D] = 1;
        mapping4[- 256] = 1;
        mapping5["a"] = 1;
        mapping6[123] = Entity(1, "mapping6");

    }

    
}



```

get contract variable value
```go
import "github.com/MetaplasiaTeam/storagescan"

var RopstenTestNet = storagescan.GenGetStorageValueFunc(context.Background(), "https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161", common.HexToAddress("0x24302f327764f94c15d930f5ac70d362b4a156f9"))
intValue := storagescan.SolidityInt{
        SlotIndex: common.HexToHash("0x0"),
        Length:    256,
        Offset:    0,
    }
    
log.Printf("value:%v\n", intValue.Value(RopstenTestNet))
// output: value:-8

sructValue := storagescan.SolidityStruct{
    SlotIndex: common.HexToHash("0xa"),
    FiledValueMap: map[string]struct {
        V storagescan.Variable
        I uint64
    }{
        "id": {
            V: &storagescan.SolidityUint{
                Length: 256,
            },
            I: 0,
        },
        "value": {
            V: &storagescan.SolidityString{},
            I: 1,
        },
    },
}

structIValue := sructValue.Value(RopstenTestNet())
log.Printf("structValue:%v\n", structIValue)
// output: structValue: struct{id:1 value:entity}

valueFieldValue := structValue.(storagescan.StructValueI).Field("value")
log.Printf("'valueFieldValue:%v\n", valueFieldValue)
// output: valueFieldValue: entity

    







```

For more use cases, refer to `storagescan.go` in the `examples` directory.
