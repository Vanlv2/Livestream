import { contractABI, contractAddress } from './contractConfig.js';

let web3;
export let contract = null; // Ban đầu contract là null

export async function init() {
  if (contract) {
    // Nếu contract đã được khởi tạo, thoát ra khỏi hàm
    console.log("Contract đã được khởi tạo trước đó:", contract);
    return;
  }

  if (window.ethereum) {
    web3 = new Web3(window.ethereum);
    try {
      await window.ethereum.request({ method: "eth_requestAccounts" });
      console.log("Ethereum accounts connected");
    } catch (error) {
      console.error("Người dùng từ chối kết nối ví");
      return;
    }
  } else if (window.web3) {
    web3 = new Web3(window.web3.currentProvider);
  } else {
    console.log("Không tìm thấy trình duyệt hỗ trợ Ethereum. Cài MetaMask!");
    return;
  }

  // Khởi tạo contract nếu chưa khởi tạo
  contract = new web3.eth.Contract(contractABI, contractAddress);
  console.log("Contract đã được khởi tạo:", contract);
}
