pragma solidity ^0.8.4;

contract GravitonAuditor {
    struct swap {
        bytes uuid;
        bytes sender;
        string source_chain;
        bytes receiver;
        string destination_chain;
        uint256 amount;
        
        bytes source_transaction;
        bytes destination_transaction;
    }

    mapping(bytes32 => swap) public proof;
    function addSwap(bytes memory uuid, bytes memory sender, string memory source_chain, bytes memory receiver, string memory destination_chain, uint256 amount, bytes memory source_transaction,  bytes memory destination_transaction) public {
        swap memory a = swap({
            uuid: uuid,
            sender: sender,
            source_chain: source_chain,
            receiver: receiver,
            destination_chain: destination_chain,
            amount: amount,
            source_transaction: source_transaction,
            destination_transaction: destination_transaction
        });
        bytes32 k = keccak256(abi.encodePacked(uuid,sender,source_chain,receiver,destination_chain,amount));
        proof[k] = a;
    }
    function checkSwap(bytes calldata uuid,  bytes calldata sender, bytes calldata source_chain, bytes calldata reciever, bytes32 destination_chain, uint256 amount) public view returns (swap memory){
        bytes32 k = keccak256(abi.encodePacked(uuid,sender,source_chain,reciever,destination_chain,amount));
        swap memory res = proof[k];
        return res;
    }
}