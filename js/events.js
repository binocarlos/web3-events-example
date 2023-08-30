const fs = require('fs')
const bluebird = require('bluebird')
const path = require('path')
const ethers = require('ethers')

const PRIVATE_KEY = process.env.PRIVATE_KEY
const RPC_URL = process.env.RPC_URL
const CHAIN_ID = process.env.CHAIN_ID

if(!PRIVATE_KEY) {
  console.error(`PRIVATE_KEY not set`)
  process.exit(1)
}

if(!RPC_URL) {
  console.error(`RPC_URL not set`)
  process.exit(1)
}

if(!CHAIN_ID) {
  console.error(`CHAIN_ID not set`)
  process.exit(1)
}

const ABI_PATH = path.join(__dirname, '..', 'artifacts', 'SimpleEvent.abi')
const BYTECODE_PATH = path.join(__dirname, '..', 'artifacts', 'SimpleEvent.bin')

async function main() {
  if(!fs.existsSync(ABI_PATH)) throw new Error(`ABI_PATH ${ABI_PATH} is missing`)
  if(!fs.existsSync(BYTECODE_PATH)) throw new Error(`BYTECODE_PATH ${BYTECODE_PATH} is missing`)

  const abi = JSON.parse(fs.readFileSync(ABI_PATH))
  const bytecode = fs.readFileSync(BYTECODE_PATH).toString()

  const provider = RPC_URL.indexOf('ws') == 0 ?
    new ethers.providers.WebSocketProvider(RPC_URL, {
      chainId: parseInt(CHAIN_ID)
    }) :
    new ethers.providers.JsonRpcProvider(RPC_URL, {
      chainId: parseInt(CHAIN_ID)
    })
  
  const wallet = new ethers.Wallet(PRIVATE_KEY)
  const signer = wallet.connect(provider)
  const contractFactory = new ethers.ContractFactory(abi, bytecode, signer)
  const deployedContract = await contractFactory.deploy()

  console.log(`Contract deployed: ${deployedContract.address}`)

  deployedContract.on('NumberSet', (val, ev) => {
    console.log(`saw event value: ${val} from tx ${ev.transactionHash}`)
  })

  let num = 0

  while (true) {
    num++
    const tx = await deployedContract.setNumber(num)
    console.log(`tx: ${tx.hash}`)
    await tx.wait()
    console.log('')
    await bluebird.delay(1000)
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
