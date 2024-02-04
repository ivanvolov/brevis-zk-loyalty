import * as dotenv from 'dotenv';
import { DeployFunction } from 'hardhat-deploy/types';
import { HardhatRuntimeEnvironment } from 'hardhat/types';

dotenv.config();

const deployFunc: DeployFunction = async (hre: HardhatRuntimeEnvironment) => {
  const { deployments, getNamedAccounts } = hre;
  const { deploy } = deployments;
  const { deployer } = await getNamedAccounts();

  const args: string[] = ['0x7d4ed4077826Bf9BEB6232240A39db251e513e16']; // BrevisProof contract address on sepolia
  const deployment = await deploy('AccountAge', {
    from: deployer,
    log: true,
    args: args
  });

  await hre.run('verify:verify', {
    address: deployment.address,
    constructorArguments: args ?? deployment.args
  });
};

deployFunc.tags = ['AccountAge'];
deployFunc.dependencies = [];
export default deployFunc;
